package envvar

import (
	"testing"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server"
)

const MongoDBUri = "mongodb://rootuser:rootpassword@mongodb:27017"
const MongoDatabase = "restaurant-file-store-db"

func SetAllNeceassaryEnvironemntVariables(t *testing.T) {
	t.Setenv(server.OtelMetricEndpointEnv, "http://localhost")
	t.Setenv(server.ServiceNameEnv, "blub")
	t.Setenv(server.MongoDBUriEnv, MongoDBUri)
	t.Setenv(server.MongoDbDatabaseNameEnv, MongoDatabase)
}
