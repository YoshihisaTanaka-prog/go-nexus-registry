package dbClient

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"reflect"
	"web_app/ent"
	"web_app/ent/savedlibrary"
)

type savedLibraryNameSpace struct {}

var savedLibraryNS = savedLibraryNameSpace{} 

func createSavedLibraryUnit(
	client *ent.Client,
	id string,
	name string,
	version string,
	v1 int, v2 int, v3 int,
) (ok bool) {
	_, err := client.SavedLibrary.Create().SetID(id).
		SetName(name).
		SetVersion(version).
		SetV1(v1).
		SetV2(v2).
		SetV3(v3).
		Save(*ctx)
		
	if err == nil {
		return true
	}
	fmt.Fprintln(os.Stderr, err)
	return false
}

func (savedLibraryNameSpace)Create(id uuid.UUID, name string, version string) (ok bool) {
	localId := fmt.Sprintf("%s", id)

	var v1, v2, v3 int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &v1, &v2, &v3); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	if ok := createSavedLibraryUnit(
		psqlClient,
		localId, name, version,
		v1, v2, v3,
	); !ok {
		return false
	}

	return createSavedLibraryUnit(
		ramClient,
		localId, name, version,
		v1, v2, v3,
	)
}

func (savedLibraryNameSpace)FindById(id uuid.UUID) (ok bool) {
	library, err := ramClient.SavedLibrary.Query().
		Where(savedlibrary.ID(fmt.Sprintf("%s", id))).
		All(*ctx)
	if err == nil {
		fmt.Fprintln(os.Stdout, library, reflect.ValueOf(library))
		return true
	}
	return false
}

func (savedLibraryNameSpace)FindByNameAndVersion(name string, version string) (id uuid.UUID, status string, err error) {
	library, err :=psqlClient.SavedLibrary.Query().
		Where(
			savedlibrary.Name(name),
			savedlibrary.Version(version),
		).
		Only(*ctx)
	if err == nil {
		uuId, err := uuid.Parse(library.ID)
		if err != nil {
			return uuid.Nil, "", err
		}
		return uuId, library.Status, nil
	}
	return uuid.Nil, "", err
}

func (savedLibraryNameSpace)OnUploaded(id uuid.UUID) (ok bool) {
	_, err := psqlClient.SavedLibrary.
		UpdateOneID(fmt.Sprintf("%s", id)).
		SetStatus("uploaded").
		Save(*ctx)
	if err == nil {
		return true
	}
	return false
}

func (savedLibraryNameSpace)OnUploadFailed(id uuid.UUID) (ok bool) {
	_, err := psqlClient.SavedLibrary.
		UpdateOneID(fmt.Sprintf("%s", id)).
		SetStatus("failed").
		Save(*ctx)
	if err == nil {
		return true
	}
	return false
}

var SavedLibrary = savedLibraryNS