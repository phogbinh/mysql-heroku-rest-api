package databaseutil

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/phogbinh/mysql-heroku-rest-api/symbolutil"
)

const (
	DatabaseUsersTableName = "users"
	userNameAttribute      = "name"
	userPasswordAttribute  = "password"
	// Errors
	errorText                                          = "Error "
	errorDatabaseTableText                             = " the database table " + DatabaseUsersTableName
	errorSelectGetAllUsersFromDatabaseUsersTable       = errorText + "selecting all users from" + errorDatabaseTableText + symbolutil.Colon
	errorScanGetAllUsersFromDatabaseUsersTable         = errorText + "scanning all users from" + errorDatabaseTableText + symbolutil.Colon
	errorGetUserFromContext                            = errorText + "getting user from context" + symbolutil.Colon
	errorPrepareInsertUserToDatabaseUsersTable         = errorText + "preparing to insert user to" + errorDatabaseTableText + symbolutil.Colon
	errorInsertUserToDatabaseUsersTable                = errorText + "inserting user to" + errorDatabaseTableText + symbolutil.Colon
	errorSelectGetUserFromDatabaseUsersTable           = errorText + "selecting an user from" + errorDatabaseTableText + symbolutil.Colon
	errorScanGetUserFromDatabaseUsersTable             = errorText + "scanning an user from" + errorDatabaseTableText + symbolutil.Colon
	errorPrepareUpdateUserPasswordToDatabaseUsersTable = errorText + "preparing to update user password to" + errorDatabaseTableText + symbolutil.Colon
	errorUpdateUserPasswordToDatabaseUsersTable        = errorText + "updating user password to" + errorDatabaseTableText + symbolutil.Colon
	errorPrepareDeleteUserFromDatabaseUsersTable       = errorText + "preparing to delete user to" + errorDatabaseTableText + symbolutil.Colon
	errorDeleteUserFromDatabaseUsersTable              = errorText + "deleting user to" + errorDatabaseTableText + symbolutil.Colon
)

// An User represents an user tuple in the database table 'users'.
type User struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// A Status contains an http status code and its associated error message.
type Status struct {
	httpStatusCode int
	errorMessage   string
}

// CreateDatabaseUsersTableIfNotExists creates a table named 'users' for the given database pointer if the table has not already existed.
func CreateDatabaseUsersTableIfNotExists(databasePtr *sql.DB) error {
	_, createTableError := databasePtr.Exec("CREATE TABLE IF NOT EXISTS users (" + userNameAttribute + " VARCHAR(255) PRIMARY KEY, " + userPasswordAttribute + " VARCHAR(255) NOT NULL)")
	return createTableError
}

// ResponseJsonOfAllUsersFromDatabaseUsersTable responses to the client the json of all users from the database table 'users'.
func ResponseJsonOfAllUsersFromDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		users, status := getAllUsersFromDatabaseUsersTable(databasePtr)
		if status.httpStatusCode != http.StatusOK {
			context.String(status.httpStatusCode, status.errorMessage)
			return
		}
		context.JSON(http.StatusOK, users)
	}
}

func getAllUsersFromDatabaseUsersTable(databasePtr *sql.DB) ([]User, Status) {
	var users []User
	selectRowsPtr, selectError := databasePtr.Query("SELECT * FROM " + DatabaseUsersTableName)
	if selectError != nil {
		return nil, Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorSelectGetAllUsersFromDatabaseUsersTable + selectError.Error()}
	}
	defer selectRowsPtr.Close()
	for selectRowsPtr.Next() {
		var user User
		scanError := selectRowsPtr.Scan(&user.Name, &user.Password)
		if scanError != nil {
			return nil, Status{
				httpStatusCode: http.StatusInternalServerError,
				errorMessage:   errorScanGetAllUsersFromDatabaseUsersTable + scanError.Error()}
		}
		users = append(users, user)
	}
	return users, Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

func CreateUserToDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		user, getStatus := getUserFromContext(context)
		if getStatus.httpStatusCode != http.StatusOK {
			context.String(getStatus.httpStatusCode, getStatus.errorMessage)
			return
		}
		insertStatus := insertUserToDatabaseUsersTable(user, databasePtr)
		if insertStatus.httpStatusCode != http.StatusOK {
			context.String(insertStatus.httpStatusCode, insertStatus.errorMessage)
			return
		}
		context.JSON(http.StatusOK, user)
	}
}

func getUserFromContext(context *gin.Context) (User, Status) {
	var user User
	bindError := context.ShouldBindJSON(&user)
	if bindError != nil {
		return user, Status{
			httpStatusCode: http.StatusBadRequest,
			errorMessage:   errorGetUserFromContext + bindError.Error()}
	}
	return user, Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

func insertUserToDatabaseUsersTable(user User, databasePtr *sql.DB) Status {
	preparedStatement, prepareError := databasePtr.Prepare("INSERT INTO " + DatabaseUsersTableName + " VALUES(?, ?)")
	if prepareError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorPrepareInsertUserToDatabaseUsersTable + prepareError.Error()}
	}
	_, insertError := preparedStatement.Exec(user.Name, user.Password)
	if insertError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorInsertUserToDatabaseUsersTable + insertError.Error()}
	}
	return Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

func ReturnJsonOfUserFromDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameAttribute)
		user, status := getUserFromDatabaseUsersTable(userName, databasePtr)
		if status.httpStatusCode != http.StatusOK {
			context.String(status.httpStatusCode, status.errorMessage)
			return
		}
		context.JSON(http.StatusOK, user)
	}
}

func getUserFromDatabaseUsersTable(userName string, databasePtr *sql.DB) (User, Status) {
	var user User
	selectRows, selectError := databasePtr.Query("SELECT * FROM "+DatabaseUsersTableName+" WHERE "+userNameAttribute+" = ?", userName)
	if selectError != nil {
		return user, Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorSelectGetUserFromDatabaseUsersTable + selectError.Error()}
	}
	defer selectRows.Close()
	for selectRows.Next() {
		scanError := selectRows.Scan(&user.Name, &user.Password)
		if scanError != nil {
			return user, Status{
				httpStatusCode: http.StatusInternalServerError,
				errorMessage:   errorScanGetUserFromDatabaseUsersTable + scanError.Error()}
		}
	}
	return user, Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

func UpdatePasswordAndReturnJsonOfUserFromDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameAttribute)
		userOfNewPassword, getStatus := getUserFromContext(context)
		if getStatus.httpStatusCode != http.StatusOK {
			context.String(getStatus.httpStatusCode, getStatus.errorMessage)
			return
		}
		if userName != userOfNewPassword.Name {
			context.String(http.StatusBadRequest, userName+" cannot change the password of "+userOfNewPassword.Name+".")
			return
		}
		updateStatus := updateUserPasswordToDatabaseUsersTable(userOfNewPassword, databasePtr)
		if updateStatus.httpStatusCode != http.StatusOK {
			context.String(updateStatus.httpStatusCode, updateStatus.errorMessage)
			return
		}
		context.JSON(http.StatusOK, userOfNewPassword)
	}
}

func updateUserPasswordToDatabaseUsersTable(userOfNewPassword User, databasePtr *sql.DB) Status {
	preparedStatement, prepareError := databasePtr.Prepare("UPDATE " + DatabaseUsersTableName + " SET " + userPasswordAttribute + " = ? WHERE " + userNameAttribute + " = ?")
	if prepareError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorPrepareUpdateUserPasswordToDatabaseUsersTable + prepareError.Error()}
	}
	_, updateError := preparedStatement.Exec(userOfNewPassword.Password, userOfNewPassword.Name)
	if updateError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorUpdateUserPasswordToDatabaseUsersTable + updateError.Error()}
	}
	return Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

func DeleteAndReturnJsonOfUserFromDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameAttribute)
		deleteUserFromDatabaseUsersTable(userName, databasePtr)
		context.JSON(http.StatusOK, gin.H{"name": userName})
	}
}

func deleteUserFromDatabaseUsersTable(userName string, databasePtr *sql.DB) Status {
	preparedStatement, prepareError := databasePtr.Prepare("DELETE FROM " + DatabaseUsersTableName + " WHERE " + userNameAttribute + " = ?")
	if prepareError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorPrepareDeleteUserFromDatabaseUsersTable + prepareError.Error()}
	}
	_, deleteError := preparedStatement.Exec(userName)
	if deleteError != nil {
		return Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorDeleteUserFromDatabaseUsersTable + prepareError.Error()}
	}
	return Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}
