## Steps to setup the project

1. `schema.sql` file for the table description.
2. Specify the Oauth clientID and clientSecret in the main.go
3. Replace the callback url(if changed) in GCP too.
4. Specify the MySQL config parameters.
5. Specify the Redis config params in main.go
6. `go get` all the dependencies
7. `go run main.go` to start the server,