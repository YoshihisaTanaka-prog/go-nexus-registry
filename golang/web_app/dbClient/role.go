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
	if mode == "neplus" {
		createLocalBase(mode, name, "__system__")
	}
	roles, err := psqlClient.Role.Query().Where(role.Mode(mode)).All(*ctx)
	if err != nil {
		customError.Exit1("DB Error: Creating", mode, "mode Role. Record:", err)
	}
	if len(roles) == 0 {
		_, err = createLocalBase(mode, name, "__system__")
		if err != nil {
			customError.Exit1("DB Error: Creating", mode, "mode Role. Record:", err)
		}
	}
}

func initRole() {
	var wg sync.WaitGroup
	wg.Add(5)
	go initRoleUnit(&wg, "admins",  "Nexus管理者")
	go initRoleUnit(&wg, "viewers", "閲覧者")
	go initRoleUnit(&wg, "neplus",  "NePlus管理者")
	go initRoleUnit(&wg, "neplus",  "NePlus編集者")
	go initRoleUnit(&wg, "apis",    "admins")
	wg.Wait()
	fmt.Fprintln(os.Stdout, "初期データを投入しました。")
}

func createLocalBase(mode string, name string, requestedBy string) (role *ent.Role, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	localId := fmt.Sprintf("%s", id)
	isForNexus := slices.Contains([]string{"admins", "viewers", "custom"}, mode)
	
	return psqlClient.Role.Create().SetID(localId).
		SetName(name).
		SetMode(mode).
		SetIsForNexus(isForNexus).
		SetRequestedBy(requestedBy).
		Save(*ctx)
}

func (roleNameSpace)FindById(id string) (foundRole *ent.Role, err error) {
	roles, err := psqlClient.Role.Query().Where(role.ID(id)).All(*ctx)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	return roles[0], nil
}

func (roleNameSpace)GetNexusRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().Where(role.IsForNexus(true)).Order(role.ByName()).All(*ctx)
}

func (roleNameSpace)GetViewersRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().Where(role.Mode("viewers")).All(*ctx)
}

func (roleNameSpace)GetAllRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().All(*ctx)
}

func (roleNameSpace)GetNePlusRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().Where(role.Mode("neplus")).All(*ctx)
}

func (roleNameSpace)GetApiRoles() (roles []*ent.Role, err error) {
	return psqlClient.Role.Query().Where(role.Mode("apis")).All(*ctx)
}

func (roleNameSpace)CreateNexusRole(name string, requestedBy string) (role *ent.Role, err error) {
	return createLocalBase("custom", name, requestedBy)
}

func updateLocalBase(id string, name string, requestedBy string, isForNexus bool) (newRole *ent.Role, err error) {
	foundRole, err := psqlClient.Role.Query().Where(role.ID(id)).All(*ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Searching Role Error:", err)
		return nil, err
	}
	if len(foundRole) > 0 {
		if isForNexus != foundRole[0].IsForNexus {
			customError.Exit1("DB Error: Updating Role. Record:", id, "The isForNexus column cannot be updated.")
		}
	}

	newRole, err = psqlClient.Role.UpdateOneID(id).
		SetName(name).
		SetRequestedBy(requestedBy).
		Save(*ctx)
	
	if err != nil {
		fmt.Fprintln(os.Stderr, "Updating Role Error", err)
		return nil, err
	}
	return newRole, err
}

func (roleNameSpace)UpdateNexusRole(id string, name string, requestedBy string) (role *ent.Role, err error) {
	return updateLocalBase(id, name, requestedBy, true)
}

func (roleNameSpace)DeleteRole(id string) (err error) {
	return psqlClient.Role.DeleteOneID(id).Exec(*ctx)
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
