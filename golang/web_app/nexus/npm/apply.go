package npm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
	"web_app/dbClient"
	"web_app/pubsub"
)

type ApplyProps struct {
	Kind   string `json:"kind"  binding:"required"`
	Name   string `json:"name"  binding:"required"`
	Index *int    `json:"index" binding:"required"`
	V1    *int    `json:"v1"`
	V2    *int    `json:"v2"`
	V3    *int    `json:"v3"`
}

var dockerImageMap = map[string]string{
	"npm": "node:22",
}

func InstallLibraries(body ApplyProps, userId string, uuid uuid.UUID, dockerImageName string) (ok bool) {
	libraryArg := body.Name

	if body.V1 != nil {
		libraryArg += fmt.Sprintf("@%d", *body.V1)
		if body.V2 != nil {
			libraryArg += fmt.Sprintf(".%d", *body.V2)
			if body.V3 != nil {
				libraryArg += fmt.Sprintf(".%d", *body.V3)
			}
		}
	}

	cmd := exec.Command(
		"docker", "run",
		"--rm",
		"-v", fmt.Sprintf("%s/tmp/npm/%s:/app", os.Getenv("ROOT_DIR_PATH"), uuid),
		"-v", fmt.Sprintf("%s/templates/npmrc.txt:/app/.npmrc", os.Getenv("ROOT_DIR_PATH")),
		"-w", "/app",
		dockerImageName,
		"npm", "i", libraryArg,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ライブラリ", body.Kind, "のインストールに失敗しました。\n", err)
		pubsub.Publish(userId, uuid, -1, 1)
		return false
	}

	go func(){
		// node_modulesを削除
		if err := os.RemoveAll(fmt.Sprintf("/app/tmp/npm/%s/node_modules", uuid)); err != nil {
			fmt.Fprintln(os.Stderr, body.Kind, "の node_modules フォルダの削除に失敗しました。\n", err)
			pubsub.Publish(userId, uuid, -2, 1)
		}
	}()

	pubsub.Publish(userId, uuid, -1, 0)
	return true
}

func AuditLibraries(libName string, userId string, uuid uuid.UUID, dockerImageName string) (ok bool) {
	cmd := exec.Command(
		"docker", "run",
		"--rm",
		"-v", fmt.Sprintf("%s/tmp/npm/%s:/app", os.Getenv("ROOT_DIR_PATH"), uuid),
		"-v", fmt.Sprintf("%s/templates/npmrc.txt:/app/.npmrc", os.Getenv("ROOT_DIR_PATH")),
		"-w", "/app",
		dockerImageName,
		"npm", "audit", "fix",
	)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err == nil {
		pubsub.Publish(userId, uuid, -3, 0)
		return true
	} else {
		fmt.Fprintln(os.Stderr, "ライブラリ", libName, "の修正に失敗しました。\n", err)
		pubsub.Publish(userId, uuid, -3, 1)
		return false
	}
}

type ParsedMainLibrary struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	NumOfLibraries int    `json:"numOfLibraries"`
}

type ParsedSubLibrary struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Resolved string `json:"resolved"`
}

type parsedLibraries struct {
	MainLibrary    ParsedMainLibrary   `json:"mainLibrary"`
	SubLibraries []ParsedSubLibrary `json:"subLibraries"`
}

func ParseLibraries(libName string, userId string, uuid uuid.UUID, dockerImageName string) (mainLibrary *ParsedMainLibrary, subLibraries []ParsedSubLibrary, err error) {
	cmd := exec.Command(
		"docker", "run",
		"--rm",
		"-v", fmt.Sprintf("%s/tmp/npm/%s:/app", os.Getenv("ROOT_DIR_PATH"), uuid),
		"-v", fmt.Sprintf("%s/templates/npmrc.txt:/app/.npmrc", os.Getenv("ROOT_DIR_PATH")),
		"-v", fmt.Sprintf("%s/templates/npm-parser.js:/app/index.js", os.Getenv("ROOT_DIR_PATH")),
		"-w", "/app",
		dockerImageName,
		"node", "index",
	)
	
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ライブラリ", libName, "の情報取得に失敗しました。\n", err)
		pubsub.Publish(userId, uuid, -4, 1)
		return nil, []ParsedSubLibrary{}, err
	}
	
	var libraryInfo *parsedLibraries
	if err := json.Unmarshal(out, &libraryInfo); err != nil {
		fmt.Fprintln(os.Stderr, "ライブラリ", libName, "の情報取得に失敗しました。\n", err)
		pubsub.Publish(userId, uuid, -4, 1)
		return nil, []ParsedSubLibrary{}, err
	}
	return &libraryInfo.MainLibrary, libraryInfo.SubLibraries, nil
}

func UploadLibraries(libKind string, subLibraries []ParsedSubLibrary) {
	var wg sync.WaitGroup
	var nexusConfig = struct {
		URL        string
		Username   string
		Password   string
		Repository string
	} {
		"http://nexus:8081",
		"_api_npm",
		"apiNpm",
		"npm-staging",
	}
	for i, subLibrary := range subLibraries {
		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Millisecond * time.Duration(i * 50))

			name, version, resolvedUrl := subLibrary.Name, subLibrary.Version, subLibrary.Resolved
			fmt.Fprintln(os.Stdout, "Processing:", name, "version", version)

			library, doSkip, err := dbClient.SavedLibrary.FindOrCreate(libKind, name, version)

			if err != nil {
				fmt.Fprintln(os.Stderr, "  ", err)
				return
			}

			if doSkip {
				status := library.Status
				if (status == "uploading") {
					fmt.Fprintln(os.Stderr, "   The library:", name, "for", libKind, "is been pushing on other process.")
				}
				if (status == "uploaded") {
					fmt.Fprintln(os.Stderr, "   The library:", name, "for", libKind, "has already been pushed.")
				}
				return
			}

			// --- 1. ファイルをダウンロードしてメモリにキャッシュ ---
			
			fmt.Fprintln(os.Stdout, "   -> Downloading from", resolvedUrl)
			downloadResp, _ := http.Get(resolvedUrl)
			defer downloadResp.Body.Close()

			if downloadResp.StatusCode != http.StatusOK {
				fmt.Fprintln(os.Stderr, "bad status on download:", downloadResp.Status)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}

			fileBytes, err := io.ReadAll(downloadResp.Body)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to read response body:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}
			fmt.Fprintln(os.Stdout, "   -> Successfully downloaded", len(fileBytes), "bytes into memory.")

			// --- 2. Nexusへアップロードするリクエストを準備 ---

			uploadReqBody := &bytes.Buffer{}
			writer := multipart.NewWriter(uploadReqBody)

			// アップロードするファイル名
			uploadFilename := fmt.Sprintf("%s-%s.tgz", name, version)
			// ファイルパートを作成
			part, err := writer.CreateFormFile("npm.asset", uploadFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to create form file:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}
			// メモリ上のバイトデータをパートにコピー
			_, err = part.Write(fileBytes)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to write file bytes to form:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}
			// multipart writerを閉じて、終端境界を追加
			err = writer.Close()
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to close multipart writer:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}
			
			// --- 3. Nexusへアップロードを実行 ---
			uploadURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", nexusConfig.URL, nexusConfig.Repository)
			uploadReq, err := http.NewRequest("POST", uploadURL, uploadReqBody)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to create upload request:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}

			uploadReq.SetBasicAuth(nexusConfig.Username, nexusConfig.Password)

			// Content-Kindヘッダーをmultipart writerが生成したものに設定 (境界情報が含まれる)
			uploadReq.Header.Set("Content-Type", writer.FormDataContentType())

			fmt.Fprintln(os.Stdout, "   -> Uploading to", uploadURL)
			uploadClient := &http.Client{}
			uploadResp, err := uploadClient.Do(uploadReq)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to execute upload request:", err)
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
			}
			defer uploadResp.Body.Close()

			if uploadResp.StatusCode >= 400 {
				// エラーレスポンスのボディを読んで詳細を確認
				errorBody, _ := io.ReadAll(uploadResp.Body)
				fmt.Fprintln(os.Stderr, "upload failed with status", uploadResp.Status, ":", string(errorBody))
				dbClient.SavedLibrary.OnUploadFailed(library.ID)
				return
			}

			fmt.Fprintln(os.Stdout, "   -> Successfully uploaded", name, "version", version, ".\n")
			time.Sleep(time.Second)
			dbClient.SavedLibrary.OnUploaded(library.ID)
		} ()
	}
	wg.Wait()
}