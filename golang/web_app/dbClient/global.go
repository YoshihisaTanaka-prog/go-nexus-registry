package dbClient

import (
	"os"
	"web_app/ent"
)

var dsn = "host=postgres port=5432 user=" + os.Getenv("POSTGRES_USER") + " dbname=addon_db password=" + os.Getenv("POSTGRES_PASSWORD") + " sslmode=disable"
var client *ent.Client