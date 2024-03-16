package envvar

import (
	"testing"

	"github.com/kinneko-de/restaurant-file-store-svc/internal/app/server"
	"github.com/kinneko-de/restaurant-file-store-svc/test/testing/mongodb"
)

const MongoDatabase = "restaurant-file-store-db"

func SetAllNeceassaryEnvironemntVariables(t *testing.T) {
	t.Setenv(server.OtelMetricEndpointEnv, "http://localhost")
	t.Setenv(server.ServiceNameEnv, "blub")
	t.Setenv(server.MongoDBUriEnv, mongodb.MongoDbServer)
	t.Setenv(server.MongoDbDatabaseNameEnv, MongoDatabase)
}
