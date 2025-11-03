package dbClient

import (
	"context"
	"os"
	"sync"
	"web_app/ent"
)

var (
	psqlDsn = "host=postgres port=5432 user=" + os.Getenv("POSTGRES_USER") + " dbname=addon_db password=" + os.Getenv("POSTGRES_PASSWORD") + " sslmode=disable"
	psqlClient *ent.Client
	ramDsn = "file::memory:?cache=shared&_fk=1"
	ramClient *ent.Client
	ctx *context.Context
	ramMutex sync.Mutex
	psqlMutex sync.Mutex
)
