package dbClient

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"os"
	"slices"
	"sync"
	"web_app/customError"
	"web_app/ent"
	"web_app/ent/role"
)

type roleNameSpace struct {}
var roleNS = roleNameSpace{} 

func initRoleUnit(wg *sync.WaitGroup, mode string, name string) {
	defer wg.Done()
	roles, err := psqlClient.Role.Query().Where(role.Mode(mode)).All(*ctx)
	if err != nil {
		customError.Exit1("DB Error: Creating", mode, "Record:", err)
	}
	if (len(roles) == 0) {
		_, err = createLocalBase(mode, name, "__system__")
	}
	if err != nil {
		customError.Exit1("DB Error: Creating", mode, "Record:", err)
	}
}

func initRole() {
	var wg sync.WaitGroup
	wg.Add(3)
	go initRoleUnit(&wg, "viewers", "閲覧者")
	go initRoleUnit(&wg, "editors", "編集者")
	go initRoleUnit(&wg, "admins",  "カスタムアプリ管理者")
	wg.Wait()
	fmt.Fprintln(os.Stdout, "初期データを投入しました。")
}

func createLocalBase(mode string, name string, requestedBy string) (role *ent.Role, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	localId := fmt.Sprintf("%s", id)
	isForNexus := slices.Contains([]string{"viewers", "editors", "custom"}, mode)
	
	return psqlClient.Role.Create().SetID(localId).
		SetName(name).
		SetMode(mode).
		SetIsForNexus(isForNexus).
		SetRequestedBy(requestedBy).
		Save(*ctx)
}

func (roleNameSpace)GetNexusRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().Where(role.IsForNexus(true)).Order(role.ByName()).All(*ctx)
}

func (roleNameSpace)CreateNexusRole(name string, requestedBy string) (role *ent.Role, err error) {
	return createLocalBase("custom", name, requestedBy)
}

func updateLocalBase(id string, name string, requestedBy string, isForNexus bool) (role *ent.Role, err error) {
	return psqlClient.Role.UpdateOneID(id).
		SetName(name).
		SetRequestedBy(requestedBy).
		SetIsForNexus(isForNexus).
		Save(*ctx)
}

func (roleNameSpace)UpdateNexusRole(id string, name string, requestedBy string) (role *ent.Role, err error) {
	return updateLocalBase(id, name, requestedBy, true)
}

func (roleNameSpace)UpdateAddonRole(id string, name string) (role *ent.Role, err error) {
	return updateLocalBase(id, name, "__system__", false)
}

func (roleNameSpace)JudgeError(err error) (statusCode int, message string) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" { // UniqueViolation
			return 409, "その名前は既に存在します。"
		}
	}

	return 400, fmt.Sprintf("%v", err)
}

var Role = roleNS
