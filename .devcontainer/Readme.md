### DB setup
The db is created by default within the devcontainer using the env configuration 
file in the [.devcontainer](/) directory. You can access the db with 
`psql -U postgres -h localhost -p 5432` or `psql -h db -U postgres broker`. 
Password is `postgres`. <br>

The container comes with `go version go1.25.5 linux/arm64` installed by default, and 
all external packages in the `go.mod` file are preinstalled by default. 
