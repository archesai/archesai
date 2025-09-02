package database

import (
	"testing"
	// "github.com/testcontainers/testcontainers-go"
	// "github.com/testcontainers/testcontainers-go/modules/postgres"
	// "github.com/testcontainers/testcontainers-go/wait"
)

// func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
// 	var (
// 		dbName = "database"
// 		dbPwd  = "password"
// 		dbUser = "user"
// 	)

// 	dbContainer, err := postgres.Run(
// 		context.Background(),
// 		"postgres:latest",
// 		postgres.WithDatabase(dbName),
// 		postgres.WithUsername(dbUser),
// 		postgres.WithPassword(dbPwd),
// 		testcontainers.WithWaitStrategy(
// 			wait.ForLog("database system is ready to accept connections").
// 				WithOccurrence(2).
// 				WithStartupTimeout(5*time.Second)),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	database = dbName
// 	password = dbPwd
// 	username = dbUser

// 	dbHost, err := dbContainer.Host(context.Background())
// 	if err != nil {
// 		return dbContainer.Terminate, err
// 	}

// 	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
// 	if err != nil {
// 		return dbContainer.Terminate, err
// 	}

// 	host = dbHost
// 	port = dbPort.Port()

// 	return dbContainer.Terminate, err
// }

func TestMain(_ *testing.M) {
	// teardown, err := mustStartPostgresContainer()
	// if err != nil {
	// 	log.Fatalf("could not start postgres container: %v", err)
	// }

	// m.Run()

	// if teardown != nil && teardown(context.Background()) != nil {
	// 	log.Fatalf("could not teardown postgres container: %v", err)
	// }
}

func TestNew(_ *testing.T) {

}

func TestHealth(_ *testing.T) {

}

func TestClose(_ *testing.T) {

}
