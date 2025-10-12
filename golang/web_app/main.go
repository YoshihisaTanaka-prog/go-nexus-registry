package main

import (
	"time"
	"web_app/api"
	"web_app/dbClient"
)

func main() {
	time.Sleep(time.Second * 1)
	dbClient.InitDb()
	time.Sleep(time.Second * 1)
	api.Start()
}
