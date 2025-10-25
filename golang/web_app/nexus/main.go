package nexus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"time"
	"web_app/dbClient"
	"web_app/nexus/npm"
	"web_app/pubsub"
)

type ApplyProps = npm.ApplyProps

type npmNameSpace struct {}

var npmNS = npmNameSpace{} 

func (npmNameSpace)Apply(c *gin.Context, body ApplyProps, userId string, uuId uuid.UUID) {
	fmt.Fprintln(os.Stdout, "apply:", userId, uuId, body)
	pubsub.PublishUuid(userId, uuId)

	time.Sleep(time.Millisecond * 100)
	dockerImageName, ok := getDockerImageName(body)
	if ok {
		c.JSON(200, gin.H{})
	} else {
		c.JSON(422, gin.H{"error": fmt.Sprintf("「%s」という種類のライブラリには未対応です。", body.Kind)})
	}
	go func(){
		if ok := npm.InstallLibraries(body, userId, uuId, dockerImageName); !ok {
			return
		}
		if ok := npm.AuditLibraries(body.Kind, userId, uuId, dockerImageName); !ok {
			return
		}
		mainLibrary, subLibraries, err := npm.ParseLibraries(body.Kind, userId, uuId, dockerImageName)
		name := (*mainLibrary).Name
		version := (*mainLibrary).Version
		if err != nil {
			return
		}
		savedUuid := uuId
		if ok := dbClient.RequestedLibrary.Create(uuId, body.Kind, name, version, userId); !ok {
			if id, status, err := dbClient.RequestedLibrary.FindByKindAndNameAndVersion(body.Kind, name, version); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			} else if (status == "uploading") {
				return
			} else {
				go func(){
					// node_modulesを削除
					targetDirName := fmt.Sprintf("/app/tmp/npm/%s", uuId)
					if err := os.RemoveAll(targetDirName); err != nil {
						fmt.Fprintln(os.Stderr, targetDirName, "フォルダの削除に失敗しました。\n", err)
						pubsub.Publish(userId, uuId, -2, 1)
					}
				}()
				savedUuid = id
			}
		}
		npm.UploadLibraries(body.Kind, subLibraries)
		dbClient.RequestedLibrary.OnUploaded(savedUuid)
		fmt.Fprintln(os.Stdout, "\nUploaded all libraries.\n\n")
	}()
}

func (npmNameSpace)GetLibraries(c *gin.Context, kind string, name string, v1 int, v2 int, v3 int, limit int) {
	npm.GetLibraries(c *gin.Context, kind, name, v1, v2, v3, limit)
}

var Npm = npmNS

var dockerImageMap = map[string]string{
	"npm": "node:22",
}

func getDockerImageName(body ApplyProps) (dockerImageName string, ok bool) {
	// PubSubの渋滞緩和用
	time.Sleep(time.Second * time.Duration(*body.Index))

	dockerImageName, exists := dockerImageMap[body.Kind]
	if !exists {
		return "", false
	}
	return dockerImageName, true
}
