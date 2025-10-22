package dbClient

import (
	"context"
	"os"
	"web_app/ent"
)

var (
	psqlDsn = "host=postgres port=5432 user=" + os.Getenv("POSTGRES_USER") + " dbname=addon_db password=" + os.Getenv("POSTGRES_PASSWORD") + " sslmode=disable"
  psqlClient *ent.Client
	ramDsn = "file:ent?mode=memory&cache=shared&_fk=1"
	ramClient *ent.Client
	ctx *context.Context
)
