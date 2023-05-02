# Go API Snake Demo _(go-api-snake-demo)_

This project demonstrates a basic REST API using Go & Fiber. See `api-test.rest` for example tests using [`REST Client`](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)_(humao.rest-client)_ VS Code extension

## Install

```bash
go mod init github.com/nubrid/go-api-snake-demo

go get -u github.com/gofiber/fiber/v2 github.com/go-playground/validator/v10

npx kill-port 3000 && go run cmd/main.go
```
