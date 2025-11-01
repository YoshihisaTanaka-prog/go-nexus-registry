package dbClient

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"slices"
	"strings"
	"time"
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

func findByIdWithTx(id string, tx *ent.Tx) (foundLib *ent.SavedLibrary, ok bool) {
	libraries, err := tx.SavedLibrary.Query().Where(savedlibrary.ID(id)).All(*ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cache Data Error:", err)
		tx.Rollback()
		return nil, false
	}

	if len(libraries) == 0 {
		return nil, true
	}
	return libraries[0], true
}

func syncSavedLibraryCache(library *ent.SavedLibrary) (ok bool) {
	
	tx, err := ramClient.Tx(*ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Creating Transaction for Cache Data Error:", err)
		return false
	}

	foundLib, ok := findByIdWithTx(library.ID, tx)

	if !ok {
		return false
	}

	if foundLib == nil {
		_, err := tx.SavedLibrary.Create().SetID(library.ID).
			SetFullName(library.FullName).
			SetSimpleName(library.SimpleName).
			SetKind(library.Kind).
			SetVersion(library.Version).
			SetV1(0).
			SetV2(0).
			SetV3(0).
			Save(*ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Creating Cache Data Error:", err)
			tx.Rollback()
			return false
		}
		go func() {
			time.Sleep(time.Second * 70)
			tx, err := ramClient.Tx(*ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Creating Transaction to Delete Cache Data Error:", err)
				return
			}
			foundLib, _ := findByIdWithTx(library.ID, tx)
			if foundLib != nil {
				if foundLib.UpdatedAt.Add(time.Minute).Before(time.Now()) {
					err := tx.SavedLibrary.DeleteOneID(library.ID).Exec(*ctx)
					if err != nil {
						fmt.Fprintln(os.Stderr, "Deleting Cache Data Error:", err)
						tx.Rollback()
					}
				}
			}
			err = tx.Commit()
			if err == nil {
				fmt.Fprintln(os.Stdout, "Deleted Cache Data id:", library.ID)
			} else {
				fmt.Fprintln(os.Stderr, "Deleting Cache Data Error:", err)
			}
		}()
	}
	
	_, err = tx.SavedLibrary.UpdateOneID(library.ID).
		SetStatus(library.Status).
		SetV1(library.V1).
		SetV2(library.V2).
		SetV3(library.V3).
		SetIsPublished(library.IsPublished).
		SetUpdatedAt(library.UpdatedAt).
		Save(*ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Synchronizing Cache Data Error:", err)
		tx.Rollback()
		return false
	}

	err = tx.Commit()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Synchronizing Cache Data Error:", err)
		return false
	}
	return true
}

func findSavedLibraryByKindAndNameAndVersionUnit(client *ent.Client, kind string, name string, version string) (*ent.SavedLibrary, error) {
	libraries, err := client.SavedLibrary.Query().
		Where(
			savedlibrary.Kind(kind),
			savedlibrary.FullName(name),
			savedlibrary.Version(version),
		).
		All(*ctx)
	
	if err != nil {
		return nil, err
	}

	if len(libraries) == 0{
		return nil, nil
	}

	return libraries[0], nil
}

func findSavedLibraryByKindAndNameAndVersion(kind string, name string, version string) (*ent.SavedLibrary, error) {
	library, _ := findSavedLibraryByKindAndNameAndVersionUnit(ramClient, kind, name, version)
	if library != nil {
		return library, nil
	}

	library, err := findSavedLibraryByKindAndNameAndVersionUnit(psqlClient, kind, name, version)
	if library != nil {
		syncSavedLibraryCache(library)
	}
	return library, err
}

func (savedLibraryNameSpace)FindOrCreate(kind string, name string, version string) (library *ent.SavedLibrary, doSkip bool, err error) {
	library, _ = findSavedLibraryByKindAndNameAndVersion(kind, name, version)
	if library != nil {
		doNextStatuses := []string{
			"failed",
		}
		return library, !(slices.Contains(doNextStatuses, library.Status)), nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, true, err
	}

	localId := fmt.Sprintf("%s", id)

	var v1, v2, v3 int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &v1, &v2, &v3); err != nil {
		return nil, true, err
	}

	splitedName := strings.Split(name, "/")

	simpleName := splitedName[len(splitedName)-1]

	library, err = psqlClient.SavedLibrary.Create().SetID(localId).
		SetFullName(name).
		SetSimpleName(simpleName).
		SetKind(kind).
		SetVersion(version).
		SetV1(v1).
		SetV2(v2).
		SetV3(v3).
		Save(*ctx)
		
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, true, err
	}

	syncSavedLibraryCache(library)
	return library, false, nil
}

func (savedLibraryNameSpace)FindLibraries(kind string, name string, v1 int, v2 int, v3 int, limit int) ([]*ent.SavedLibrary, *Cursor, error) {
	q := psqlClient.SavedLibrary.Query().
		Where(
			savedlibrary.Kind(kind),
		).
		Order(
			savedlibrary.ByFullName(),
			savedlibrary.ByV1(),
			savedlibrary.ByV2(),
			savedlibrary.ByV3(),
		)
	
	if name != "" {
		q = q.Where(
			savedlibrary.Or(
				savedlibrary.FullNameGT(name),
				savedlibrary.And(
					savedlibrary.FullNameEQ(name),
					savedlibrary.V1GT(v1),
				),
				savedlibrary.And(
					savedlibrary.FullNameEQ(name),
					savedlibrary.V1EQ(v1),
					savedlibrary.V2GT(v2),
				),
				savedlibrary.And(
					savedlibrary.FullNameEQ(name),
					savedlibrary.V1EQ(v1),
					savedlibrary.V2EQ(v2),
					savedlibrary.V3GT(v3),
				),
			),
		)
	}

	libraries, err := q.Limit(limit).All(*ctx)

	if err != nil {
		return []*ent.SavedLibrary{}, nil, err
	}

	var nextCursor *Cursor
	if len(libraries) > 0 {
		last := libraries[len(libraries) - 1]
		nextCursor = &Cursor{
			Name: last.FullName,
			V1:   last.V1,
			V2:   last.V2,
			V3:   last.V3,
		}
	}

	return libraries, nextCursor, nil
}

func (savedLibraryNameSpace)FindNewerPublishedLibraries(targetLibrary *ent.SavedLibrary) ([]*ent.SavedLibrary, error) {
	kind, name, v1, v2, v3 := targetLibrary.Kind, targetLibrary.FullName, targetLibrary.V1, targetLibrary.V2, targetLibrary.V3
	q := psqlClient.SavedLibrary.Query().
		Where(
			savedlibrary.Kind(kind),
			savedlibrary.FullName(name),
			savedlibrary.IsPublished(true),
		).
		Order(
			savedlibrary.ByV1(),
			savedlibrary.ByV2(),
			savedlibrary.ByV3(),
		)
	
	q = q.Where(
		savedlibrary.Or(
			savedlibrary.V1GT(v1),
			savedlibrary.And(
				savedlibrary.V1EQ(v1),
				savedlibrary.V2GT(v2),
			),
			savedlibrary.And(
				savedlibrary.V1EQ(v1),
				savedlibrary.V2EQ(v2),
				savedlibrary.V3GT(v3),
			),
			savedlibrary.And(
				savedlibrary.V1EQ(v1),
				savedlibrary.V2EQ(v2),
				savedlibrary.V3EQ(v3),
			),
		),
	)

	foundLibraries, err := q.All(*ctx)

	if err != nil {
		return []*ent.SavedLibrary{}, err
	}

	libraries := []*ent.SavedLibrary{}

	for _, l := range foundLibraries {
		rawVersion := fmt.Sprintf("%s.%s.%s", l.V1, l.V2, l.V3)
		if rawVersion != l.Version {
			libraries = append(libraries, l)
		}
	}

	for _, l := range foundLibraries {
		rawVersion := fmt.Sprintf("%s.%s.%s", l.V1, l.V2, l.V3)
		if rawVersion == l.Version {
			libraries = append(libraries, l)
		}
	}

	return libraries, nil
}

func findSavedLibraryByIdUnit(client *ent.Client, id string) (*ent.SavedLibrary, error) {
	libraries, err := client.SavedLibrary.Query().
		Where(savedlibrary.ID(id)).
		All(*ctx)
	if err != nil {
		return nil, err
	}
	if len(libraries) == 0 {
		return nil, nil
	}
	return libraries[0], nil
}

func findSavedLibraryById(id string) (*ent.SavedLibrary, error) {
	library, _ := findSavedLibraryByIdUnit(ramClient, id)
	if library != nil {
		return library, nil
	}
	library, err := findSavedLibraryByIdUnit(psqlClient, id)
	if err != nil {
		return nil, err
	}
	syncSavedLibraryCache(library)
	return library, nil
}

func (savedLibraryNameSpace)FindById(id string) (*ent.SavedLibrary, error) {
	return findSavedLibraryById(id)
}

func (savedLibraryNameSpace)OnUploaded(id string) (ok bool) {
	library, err := psqlClient.SavedLibrary.
		UpdateOneID(id).
		SetStatus("uploaded").
		Save(*ctx)
	if err != nil {
		return false
	}
	syncSavedLibraryCache(library)
	return true
}

func (savedLibraryNameSpace)OnUploadFailed(id string) (ok bool) {
	library, err := psqlClient.SavedLibrary.
		UpdateOneID(id).
		SetStatus("failed").
		Save(*ctx)
	if err != nil {
		return false
	}
	syncSavedLibraryCache(library)
	return true
}

func (savedLibraryNameSpace)SetStatusUpdating(id string) error {
	library, err := psqlClient.SavedLibrary.
		UpdateOneID(id).
		SetStatus("updating").
		Save(*ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	syncSavedLibraryCache(library)
	return nil
}

func (savedLibraryNameSpace)ResetStatusUpdating(id string) error {
	library, err := psqlClient.SavedLibrary.
		UpdateOneID(id).
		SetStatus("uploaded").
		Save(*ctx)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	syncSavedLibraryCache(library)
	return nil
}

func (savedLibraryNameSpace)UpdateIsPublished(id string, isPublished bool) (*ent.SavedLibrary, error) {
	library, err := psqlClient.SavedLibrary.
		UpdateOneID(id).
		SetIsPublished(isPublished).
		Save(*ctx)
	if library != nil {
		syncSavedLibraryCache(library)
	}
	return library, err
}

var SavedLibrary = savedLibraryNS