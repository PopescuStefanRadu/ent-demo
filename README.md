Project using ent.

### Build and run via Docker

```shell
docker build -t ent:latest .
```


```shell
docker run -p 8080:8080 ent
```

### Code entrypoints

Server: `cmd/http/server/main.go`

Application: `pkg/app/app.go`

### TODO
 
 - add more documentation
   - design decisions
   - public functionality when required
 - test unhappy paths
 - add config loading
 - add correlated logging and tracing
 - add short commit sha in build (see: https://docs.docker.com/build/guide/build-args/)
 - use postgres instead of sqlite3, use ory/dockertest and dind(docker in docker) to use postgres in tests
 - cache tool installation in GitHub Actions
 - ci check for branch name regex to satisfy `#issueId - message` format
