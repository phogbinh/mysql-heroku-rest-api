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
	userNameColumnName     = "name"
	userPasswordColumnName = "password"
	// Errors
	errorText                                             = "Error "
	errorDatabaseTableText                                = " the database table " + DatabaseUsersTableName
	errorSelectGetAllUsersFromDatabaseUsersTable          = errorText + "selecting all users from" + errorDatabaseTableText + symbolutil.Colon
	errorScanGetAllUsersFromDatabaseUsersTableRowsPointer = errorText + "scanning all users from" + errorDatabaseTableText + "'s rows pointer" + symbolutil.Colon
	errorGetUserFromContext                               = errorText + "getting user from context" + symbolutil.Colon
	errorPrepareInsertUserToDatabaseUsersTable            = errorText + "preparing to insert user to" + errorDatabaseTableText + symbolutil.Colon
	errorInsertUserToDatabaseUsersTable                   = errorText + "inserting user to" + errorDatabaseTableText + symbolutil.Colon
	errorSelectGetUserFromDatabaseUsersTable              = errorText + "selecting an user from" + errorDatabaseTableText + symbolutil.Colon
	errorGetManyUsersGetUserFromDatabaseUsersTable        = errorText + "want to get one but got many users from" + errorDatabaseTableText + symbolutil.Colon
	errorPrepareUpdateUserPasswordToDatabaseUsersTable    = errorText + "preparing to update user password to" + errorDatabaseTableText + symbolutil.Colon
	errorUpdateUserPasswordToDatabaseUsersTable           = errorText + "updating user password to" + errorDatabaseTableText + symbolutil.Colon
	errorPrepareDeleteUserFromDatabaseUsersTable          = errorText + "preparing to delete user to" + errorDatabaseTableText + symbolutil.Colon
	errorDeleteUserFromDatabaseUsersTable                 = errorText + "deleting user to" + errorDatabaseTableText + symbolutil.Colon
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
	_, createTableError := databasePtr.Exec("CREATE TABLE IF NOT EXISTS users (" + userNameColumnName + " VARCHAR(255) PRIMARY KEY, " + userPasswordColumnName + " VARCHAR(255) NOT NULL)")
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
	selectRowsPtr, selectError := databasePtr.Query("SELECT * FROM " + DatabaseUsersTableName)
	if selectError != nil {
		return nil, Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorSelectGetAllUsersFromDatabaseUsersTable + selectError.Error()}
	}
	defer selectRowsPtr.Close()
	return getAllUsers(selectRowsPtr)
}

func getAllUsers(databaseUsersTableRowsPtr *sql.Rows) ([]User, Status) {
	var users []User
	for databaseUsersTableRowsPtr.Next() {
		var user User
		scanError := databaseUsersTableRowsPtr.Scan(&user.Name, &user.Password)
		if scanError != nil {
			return nil, Status{
				httpStatusCode: http.StatusInternalServerError,
				errorMessage:   errorScanGetAllUsersFromDatabaseUsersTableRowsPointer + scanError.Error()}
		}
		users = append(users, user)
	}
	return users, Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

// CreateUserToDatabaseUsersTableAndResponseJsonOfUser creates the user given in the context to the database table 'users'.
// Also, it responses to the client the json of the given user.
func CreateUserToDatabaseUsersTableAndResponseJsonOfUser(databasePtr *sql.DB) gin.HandlerFunc {
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

// ResponseJsonOfUserFromDatabaseUsersTable responses to the client the json of the user given in the context parameter from the database table 'users'.
func ResponseJsonOfUserFromDatabaseUsersTable(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameColumnName)
		user, status := getUserFromDatabaseUsersTable(userName, databasePtr)
		if status.httpStatusCode != http.StatusOK {
			context.String(status.httpStatusCode, status.errorMessage)
			return
		}
		context.JSON(http.StatusOK, user)
	}
}

func getUserFromDatabaseUsersTable(userName string, databasePtr *sql.DB) (User, Status) {
	var dumpUser User
	selectRows, selectError := databasePtr.Query("SELECT * FROM "+DatabaseUsersTableName+" WHERE "+userNameColumnName+" = ?", userName)
	if selectError != nil {
		return dumpUser, Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorSelectGetUserFromDatabaseUsersTable + selectError.Error()}
	}
	defer selectRows.Close()
	users, getStatus := getAllUsers(selectRows)
	if getStatus.httpStatusCode != http.StatusOK {
		return dumpUser, getStatus
	}
	if len(users) != 1 {
		return dumpUser, Status{
			httpStatusCode: http.StatusInternalServerError,
			errorMessage:   errorGetManyUsersGetUserFromDatabaseUsersTable}
	}
	return users[0], Status{
		httpStatusCode: http.StatusOK,
		errorMessage:   ""}
}

// UpdateUserPasswordInDatabaseUsersTableAndResponseJsonOfUser updates the password of the user in the database table 'users' whose name is given in the context parameter and the requested JSON object.
// Also, it responses to the client the json of the given user.
func UpdateUserPasswordInDatabaseUsersTableAndResponseJsonOfUser(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameColumnName)
		newPasswordUser, getStatus := getUserFromContext(context)
		if getStatus.httpStatusCode != http.StatusOK {
			context.String(getStatus.httpStatusCode, getStatus.errorMessage)
			return
		}
		if userName != newPasswordUser.Name {
			context.String(http.StatusBadRequest, "The user name given in the context parameter - "+userName+" - does not match the user name provided by the requested JSON object - "+newPasswordUser.Name+".")
			return
		}
		updateStatus := updateUserPasswordToDatabaseUsersTable(newPasswordUser, databasePtr)
		if updateStatus.httpStatusCode != http.StatusOK {
			context.String(updateStatus.httpStatusCode, updateStatus.errorMessage)
			return
		}
		context.JSON(http.StatusOK, newPasswordUser)
	}
}

func updateUserPasswordToDatabaseUsersTable(userOfNewPassword User, databasePtr *sql.DB) Status {
	preparedStatement, prepareError := databasePtr.Prepare("UPDATE " + DatabaseUsersTableName + " SET " + userPasswordColumnName + " = ? WHERE " + userNameColumnName + " = ?")
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

// DeleteUserFromDatabaseUsersTableAndResponseJsonOfUserName deletes the user whose name is given in the context parameter from the database table 'users'.
// Also, it responses to the client the json of the given user name.
func DeleteUserFromDatabaseUsersTableAndResponseJsonOfUserName(databasePtr *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		userName := context.Param(userNameColumnName)
		deleteUserFromDatabaseUsersTable(userName, databasePtr)
		context.JSON(http.StatusOK, gin.H{"name": userName})
	}
}

func deleteUserFromDatabaseUsersTable(userName string, databasePtr *sql.DB) Status {
	preparedStatement, prepareError := databasePtr.Prepare("DELETE FROM " + DatabaseUsersTableName + " WHERE " + userNameColumnName + " = ?")
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
