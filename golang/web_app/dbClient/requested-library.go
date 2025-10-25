package dbClient

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"web_app/ent"
	"web_app/ent/requestedlibrary"
)

type requestedLibraryNameSpace struct {}

var requestedLibraryNS = requestedLibraryNameSpace{} 

func createRequestedLibraryUnit(
	client *ent.Client,
	id string,
	kind string,
	name string,
	version string,
	requestedBy string,
) (ok bool) {
	_, err := client.RequestedLibrary.Create().SetID(id).
		SetName(name).
		SetKind(kind).
		SetVersion(version).
		SetRequestedBy(requestedBy).
		Save(*ctx)
		
	if err == nil {
		return true
	}
	fmt.Fprintln(os.Stderr, err)
	return false
}

func (requestedLibraryNameSpace)Create(id uuid.UUID, kind string, name string, version string, requestedBy string) (ok bool) {
	localId := fmt.Sprintf("%s", id)

	if ok := createRequestedLibraryUnit(
		psqlClient,
		localId, kind, name, version,
		requestedBy,
	); !ok {
		return false
	}

	return createRequestedLibraryUnit(
		ramClient,
		localId, kind, name, version,
		requestedBy,
	)
}

func (requestedLibraryNameSpace)FindById(id uuid.UUID) (ok bool) {
	_, err := ramClient.RequestedLibrary.Query().
		Where(requestedlibrary.ID(fmt.Sprintf("%s", id))).
		All(*ctx)
	if err == nil {
		return true
	}
	return false
}

func (requestedLibraryNameSpace)FindByKindAndNameAndVersion(kind string, name string, version string) (uuId uuid.UUID, status string, err error) {
	library, err := psqlClient.RequestedLibrary.Query().
		Where(
			requestedlibrary.Name(name),
			requestedlibrary.Kind(kind),
			requestedlibrary.Version(version),
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

func (requestedLibraryNameSpace)OnUploaded(id uuid.UUID) (ok bool) {
	_, err := psqlClient.RequestedLibrary.
		UpdateOneID(fmt.Sprintf("%s", id)).
		SetStatus("uploaded").
		Save(*ctx)
	if err == nil {
		return true
	}
	return false
}

var RequestedLibrary = requestedLibraryNS