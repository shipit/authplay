# Auth N Play

Simple project to exercise `OAuth` web flows.

## Changelog

Starts server and listens on port `8888`, returns simple string for `/` path. It's the only path atm.

## To Run

```bash
cd test
go run test_server.go
```

`CTRL-C` to kill the server, it captures `SIGINT` and cleans up.
