package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/phogbinh/mysql-heroku-rest-api/databaseutil"
	"github.com/phogbinh/mysql-heroku-rest-api/symbolutil"
)

const (
	databaseDriverName                 = "mysql"
	userNamePath                       = ":name"
	portEnvironmentVariableName        = "PORT"
	databaseUrlEnvironmentVariableName = "DATABASE_URL"
)

func main() {
	port := os.Getenv(portEnvironmentVariableName)
	if port == "" {
		log.Fatal("$" + portEnvironmentVariableName + " must be set.")
	}
	databasePtr := getDatabaseHandler()
	defer databasePtr.Close()
	router := gin.Default()
	initializeRouterHandlers(router, databasePtr)
	router.Run(symbolutil.Colon + port)
}

func getDatabaseHandler() *sql.DB {
	databasePtr, openDatabaseError := sql.Open(databaseDriverName, os.Getenv(databaseUrlEnvironmentVariableName))
	if openDatabaseError != nil {
		log.Fatalf("Error opening database: %q.", openDatabaseError)
	}
	return databasePtr
}

func initializeRouterHandlers(router *gin.Engine, databasePtr *sql.DB) {
	router.GET(
		symbolutil.RightSlash+databaseutil.DatabaseUsersTableName,
		databaseutil.ReturnJsonOfAllUsersFromDatabaseUsersTable(databasePtr))

	router.POST(
		symbolutil.RightSlash+databaseutil.DatabaseUsersTableName,
		databaseutil.CreateUserToDatabaseUsersTable(databasePtr))

	router.GET(
		symbolutil.RightSlash+databaseutil.DatabaseUsersTableName+symbolutil.RightSlash+userNamePath,
		databaseutil.ReturnJsonOfUserFromDatabaseUsersTable(databasePtr))

	router.PUT(
		symbolutil.RightSlash+databaseutil.DatabaseUsersTableName+symbolutil.RightSlash+userNamePath,
		databaseutil.UpdatePasswordAndReturnJsonOfUserFromDatabaseUsersTable(databasePtr))

	router.DELETE(
		symbolutil.RightSlash+databaseutil.DatabaseUsersTableName+symbolutil.RightSlash+userNamePath,
		databaseutil.DeleteAndReturnJsonOfUserFromDatabaseUsersTable(databasePtr))
}