package dbClient

import (
	"context"
	"fmt"
	"os"
)

func CreateRequestedLibrary(id string, name string, version string, v1 int, v2 int, v3 int, requestedBy string) {
	library, err := client.RequestedLibrary.Create().SetID(id).
		SetName(name).
		SetVersion(version).
		SetV1(v1).
		SetV2(v2).
		SetV3(v3).
		SetRequestedBy(requestedBy).
		Save(context.Background())
		
	if err != nil {
		fmt.Fprintln(os.Stderr, "テストデータの作成失敗")
		return
	}
	fmt.Fprintln(os.Stdout, library)
}