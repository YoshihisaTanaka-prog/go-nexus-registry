package nexus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
	"web_app/dbClient"
	"web_app/ent"
)

func GetLibraries(c *gin.Context, kind string, name string, v1 int, v2 int, v3 int, limit int) {
	fmt.Fprintln(os.Stdout, kind, limit, name, v1, v2, v3)
	libraries, nextCursor, err := dbClient.SavedLibrary.FindLibraries(kind, name, v1, v2, v3, limit)
	if err == nil {
		c.JSON(200, gin.H{
			"data":   libraries,
			"cursor": nextCursor,
			"limit":  len(libraries),
		})
		return
	}
	fmt.Fprintln(os.Stderr, err)
	c.JSON(500, gin.H{})
}

var publishRWMutex sync.RWMutex

var nexusConfig = struct {
	URL                 string
	Username            string
	Password            string
	StagingRepository   string
	PublishedRepository string
} {
	"http://nexus:8081",
	"_api_npm",
	"apiNpm",
	"npm-staging",
	"npm",
}

func publishUnit(library *ent.SavedLibrary) (ok bool) {
	safeName := url.PathEscape(library.Name)
	safeVersion := url.PathEscape(library.Version)
	downloadUrl := fmt.Sprintf("%s/repository/%s/%s/-/%s-%s.tgz", nexusConfig.URL, nexusConfig.StagingRepository, safeName, safeName, safeVersion)
	fmt.Fprintln(os.Stdout, "publish", "   -> Downloading from", downloadUrl)
	
	// --- 1. Stagingディレクトリからファイルをダウンロードしてメモリにキャッシュ ---

	downloadReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", err)
		return false
	}

	downloadReq.SetBasicAuth(nexusConfig.Username, nexusConfig.Password)
	downloadClient := &http.Client{}

	downloadResp, err := downloadClient.Do(downloadReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", err)
		return false
	}
	defer downloadResp.Body.Close()
	
	if downloadResp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "publish", "download failed with status:", downloadResp.Status)
		body, err := ioutil.ReadAll(downloadResp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "publish", "download failed:", err)
			return false
		}
		fmt.Fprintln(os.Stderr, "publish", "download failed:", string(body))
		return false
	}
	
	fileBytes, err := io.ReadAll(downloadResp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "reading file data error:", err)
		return false
	}
	fmt.Fprintln(os.Stdout, "   -> Successfully downloaded", len(fileBytes), "bytes into memory.")

	// --- 2. Nexusへアップロードするリクエストを準備 ---

	uploadReqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(uploadReqBody)

	// アップロードするファイル名
	uploadFilename := fmt.Sprintf("%s-%s.tgz", library.Name, library.Version)
	// ファイルパートを作成
	part, err := writer.CreateFormFile("npm.asset", uploadFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "creating file data for upload request body error:", err)
		return false
	}
	// メモリ上のバイトデータをパートにコピー
	_, err = part.Write(fileBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "creating file data for upload request body error:", err)
		return false
	}
	// multipart writerを閉じて、終端境界を追加
	err = writer.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "creating file data for upload request body error:", err)
		return false
	}

	// --- 3. Nexusへアップロードを実行 ---
	uploadURL := fmt.Sprintf("%s/service/rest/v1/components?repository=%s", nexusConfig.URL, nexusConfig.PublishedRepository)
	uploadReq, err := http.NewRequest("POST", uploadURL, uploadReqBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "failed to create upload request:", err)
		return false
	}

	uploadReq.SetBasicAuth(nexusConfig.Username, nexusConfig.Password)

	// Content-Kindヘッダーをmultipart writerが生成したものに設定 (境界情報が含まれる)
	uploadReq.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Fprintln(os.Stdout, "   -> Uploading to", uploadURL)
	uploadClient := &http.Client{}
	uploadResp, err := uploadClient.Do(uploadReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "publish", "failed to execute upload request:", err)
		return false
	}
	
	defer uploadResp.Body.Close()
	
	if downloadResp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "publish", "upload failed with status:", downloadResp.Status)
		body, err := ioutil.ReadAll(uploadResp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "publish", "upload failed:", err)
			return false
		}
		fmt.Fprintln(os.Stderr, "publish", "upload failed:", string(body))
		return false
	}

	fmt.Fprintln(os.Stdout, "publish", "   -> Successfully uploaded", library.Name, "version", library.Version, ".\n")

	return true
}


func publish(library *ent.SavedLibrary) (ok bool) {
	defer time.Sleep(time.Millisecond * 100)
	if !publishUnit(library) {
		return false
	}
	libraries, err := dbClient.SavedLibrary.FindNewerPublishedLibraries(library)
	if err != nil {
		fmt.Fprintln(os.Stdout, "publish", "Internal Error in DB:", err)
		return false
	}
	for _, l := range libraries {
		if !publishUnit(l) {
			return false
		}
	}
	return true
}

type searchAssetsItem struct {
	Id string `json:"id" binding:"required"`
}

type searchedAssetsResponse struct {
	Items []searchAssetsItem `json:"items" binding:"required"`
}

func unpublish(library *ent.SavedLibrary) (ok bool) {
	searchUrl := fmt.Sprintf("%s/service/rest/v1/search?repository=%s&name=%s&version=%s", nexusConfig.URL, nexusConfig.PublishedRepository, library.Name, library.Version)
	fmt.Fprintln(os.Stdout, "unpublish", "   -> Searching from", searchUrl)

	searchReq, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", err)
		return false
	}

	searchReq.SetBasicAuth(nexusConfig.Username, nexusConfig.Password)
	searchClient := &http.Client{}

	searchResp, err := searchClient.Do(searchReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", err)
		return false
	}
	defer searchResp.Body.Close()
	
	if searchResp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "unpublish", "searching asset failed with status:", searchResp.Status)
		body, err := ioutil.ReadAll(searchResp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unpublish", "searching asset failed:", err)
			return false
		}
		fmt.Fprintln(os.Stderr, "unpublish", "searching asset failed:", string(body))
		return false
	}
	
	assetsInfoBytes, err := io.ReadAll(searchResp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", "searching asset data error:", err)
		return false
	}

	var assetsInfo searchedAssetsResponse
	if err := json.Unmarshal(assetsInfoBytes, &assetsInfo); err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", "parsing asset data error:", err)
		return false
	}
	items := assetsInfo.Items
	if len(items) == 0 {
		fmt.Fprintln(os.Stderr, "unpublish", "no asset data was found")
		return false
	}

	deleteUrl := fmt.Sprintf("%s/service/rest/v1/components/%s", nexusConfig.URL, items[0].Id)
	fmt.Println(deleteUrl)

	deleteReq, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", err)
		return false
	}

	deleteReq.SetBasicAuth(nexusConfig.Username, nexusConfig.Password)
	deleteClient := &http.Client{}

	deleteResp, err := deleteClient.Do(deleteReq)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unpublish", err)
		return false
	}
	defer deleteResp.Body.Close()
	
	if deleteResp.StatusCode != http.StatusNoContent {
		fmt.Fprintln(os.Stderr, "unpublish", "deleting asset failed with status:", deleteResp.Status)
		body, err := ioutil.ReadAll(deleteResp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unpublish", "deleting asset failed:", err)
			return false
		}
		fmt.Fprintln(os.Stderr, "unpublish", "deleting asset failed:", string(body))
		return false
	}

	fmt.Fprintln(os.Stdout, "unpublish", "   -> Successfully deleted", library.Name, "version", library.Version, ".\n")
	return true
}

func updateIsPublishedUnit(c *gin.Context, library *ent.SavedLibrary, isPublished bool) {
	publishRWMutex.Lock()
	defer publishRWMutex.Unlock()

	id := library.ID

	err := dbClient.SavedLibrary.SetStatusUpdating(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "update-is-published", err)
		c.JSON(500, gin.H{
			"message": "Internal Error in DB",
		})
		return
	}
	defer dbClient.SavedLibrary.ResetStatusUpdating(id)
	
	var ok bool
	if isPublished {
		ok = publish(library)
	} else {
		ok = unpublish(library)
	}

	if !ok {
		fmt.Fprintln(os.Stderr, "update-is-published", err)
		c.JSON(500, gin.H{
			"message": "Internal Error in Nexus",
		})
		return
	}
	
	library, err = dbClient.SavedLibrary.UpdateIsPublished(id, isPublished)
	if err != nil {
		fmt.Fprintln(os.Stderr, "update-is-published", err)
		c.JSON(500, gin.H{
			"message": "Internal Error in DB",
		})
		return
	}
	c.IndentedJSON(200, library)
}

func UpdateIsPublished(c *gin.Context, id string, isPublished bool) {
	fmt.Fprintln(os.Stdout, "update-is-published", id, isPublished)
	publishRWMutex.RLock()

	library, err := dbClient.SavedLibrary.FindById(id)
	if err != nil {
		publishRWMutex.RUnlock()
		fmt.Fprintln(os.Stderr, err)
		c.JSON(500, gin.H{
			"message": "Internal Error in DB",
		})
		return
	}
	if library == nil {
		publishRWMutex.RUnlock()
		c.JSON(422, gin.H{
			"message": fmt.Sprintln("Could not find record id:", id),
		})
		return
	}

	libStatus := library.Status
	publishRWMutex.RUnlock() 

	switch libStatus {
	case "uploaded":
	  updateIsPublishedUnit(c, library, isPublished)
	default:
		c.JSON(409, gin.H{
			"message": fmt.Sprintln("target record id:", id, "is executed on other process."),
		})
	}
}