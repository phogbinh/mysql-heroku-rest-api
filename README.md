# mysql-heroku-rest-api
The project uses Go, Heroku with ClearDB (MySQL) to implement a back-end API managing a simple database table named `users`. It supports some basic RESTful APIs as described in the cURL API section below. The application URL can be found [here](https://powerful-reaches-73385.herokuapp.com).

## cURL API
### Get all users
#### Description
Get all users' names and passwords from the database table `users`.
#### Response
An JSON object containing all users' names and passwords fetched from the database.
#### Example
`curl -X GET https://powerful-reaches-73385.herokuapp.com/users`

### Create an user
#### Description
Create an user to the database table `users`.
#### Response
An JSON object containing the requested user's name and password.
#### Example
`curl -X POST -d "{ \"name\": \"bill\", \"password\": \"1\" }" https://powerful-reaches-73385.herokuapp.com/users`

### Get an user
#### Description
Get an user from the database table `users`.
#### Response
An JSON object containing the user's name and password fetched from the database.
#### Example
`curl -X GET https://powerful-reaches-73385.herokuapp.com/users/bill`

### Update an user password
#### Description
Update an user password in the database table `users`.
#### Response
An JSON object containing the requested user's name and password.
#### Example
`curl -X PUT -d "{ \"name\": \"bill\", \"password\": \"666\" }" https://powerful-reaches-73385.herokuapp.com/users/bill`

### Delete an user
#### Description
Delete an user from the database table `users`.
#### Response
An JSON object containing the user's name given in the requested URL.
#### Example
`curl -X DELETE https://powerful-reaches-73385.herokuapp.com/users/bill`

## [Solutions for Known Errors](/docs/solutions-for-known-errors.md)
A documentation describing solutions for known errors.
