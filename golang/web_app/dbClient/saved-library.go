package dbClient

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"reflect"
	"web_app/ent"
	"web_app/ent/savedlibrary"
)

type Cursor struct {
	Name string `json:"name"`
	V1   int    `json:"v1"`
	V2   int    `json:"v2"`
	V3   int    `json:"v3"`
}

type savedLibraryNameSpace struct {
	Cursor
}

var savedLibraryNS = savedLibraryNameSpace{} 

func createSavedLibraryUnit(
	client *ent.Client,
	id string,
	kind string,
	name string,
	version string,
	v1 int, v2 int, v3 int,
) (ok bool) {
	_, err := client.SavedLibrary.Create().SetID(id).
		SetName(name).
		SetKind(kind).
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

func (savedLibraryNameSpace)Create(id uuid.UUID, kind string, name string, version string) (ok bool) {
	localId := fmt.Sprintf("%s", id)

	var v1, v2, v3 int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &v1, &v2, &v3); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	if ok := createSavedLibraryUnit(
		psqlClient,
		localId, kind, name, version,
		v1, v2, v3,
	); !ok {
		return false
	}

	return createSavedLibraryUnit(
		ramClient,
		localId, kind, name, version,
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

func (savedLibraryNameSpace)FindByKindAndNameAndVersion(kind string, name string, version string) (id uuid.UUID, status string, err error) {
	library, err :=psqlClient.SavedLibrary.Query().
		Where(
			savedlibrary.Kind(kind),
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

func (savedLibraryNameSpace)FindLibraries(kind string, name string, v1 int, v2 int, v3 int, limit int) ([]*ent.SavedLibrary, *Cursor, error) {
	q := psqlClient.SavedLibrary.Query().
		Where(
			savedlibrary.Kind(kind),
		).
		Order(
			savedlibrary.ByName(),
			savedlibrary.ByV1(),
			savedlibrary.ByV2(),
			savedlibrary.ByV3(),
		).
		Limit(limit)
	
	if name != "" {
		q = q.Where(
			savedlibrary.Or(
				savedlibrary.NameGT(name),
				savedlibrary.And(
					savedlibrary.NameEQ(name),
					savedlibrary.V1GT(v1),
				),
				savedlibrary.And(
					savedlibrary.NameEQ(name),
					savedlibrary.V1EQ(v1),
					savedlibrary.V2GT(v2),
				),
				savedlibrary.And(
					savedlibrary.NameEQ(name),
					savedlibrary.V1EQ(v1),
					savedlibrary.V2EQ(v2),
					savedlibrary.V3GT(v3),
				),
			),
		)
	}

	libraries, err := q.All(*ctx)

	if err != nil {
		return []*ent.SavedLibrary{}, nil, err
	}

	var nextCursor *Cursor
	if len(libraries) > 0 {
		last := libraries[len(libraries) - 1]
		nextCursor = &Cursor{
			Name: last.Name,
			V1:   last.V1,
			V2:   last.V2,
			V3:   last.V3,
		}
	}

	return libraries, nextCursor, nil
}

var SavedLibrary = savedLibraryNS