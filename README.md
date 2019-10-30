# mysql-heroku-rest-api
The project uses Go, Heroku with ClearDB (mySQL) to implement a back-end API managing a simple database table of users. The application URL is https://powerful-reaches-73385.herokuapp.com

## cURL API
### Get all users
#### Description
Get all users name and password in the database.
#### Response
Returns an JSON of all users name and password.
#### Example
`curl -X GET https://powerful-reaches-73385.herokuapp.com/users`

### Create user
#### Description
Create an user with password to the database.
#### Response
Returns an JSON of newly created user name and password.
#### Example
`curl -X POST -d "{ \"name\": \"bill\", \"password\": \"1\" }" https://powerful-reaches-73385.herokuapp.com/users`

### Get an user
#### Description
Get an user name and password by the user name.
#### Response
Returns an JSON of the user name and password.
#### Example
`curl -X GET https://powerful-reaches-73385.herokuapp.com/users/bill`

### Update password of an user
#### Description
Update the password of an user in the database.
#### Response
Returns an JSON of the user name and password.
#### Example
`curl -X PUT -d "{ \"name\": \"bill\", \"password\": \"666\" }" https://powerful-reaches-73385.herokuapp.com/users/bill`

### Delete an user
#### Description
Delete an user by the user name.
#### Response
Returns an JSON of the user name.
#### Example
`curl -X DELETE https://powerful-reaches-73385.herokuapp.com/users/bill`
