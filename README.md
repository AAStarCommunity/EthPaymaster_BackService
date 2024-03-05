# EthPaymaster-Back
EthPaymaster relay Back-end Service

Basic flow :
![](https://raw.githubusercontent.com/jhfnetboy/MarkDownImg/main/img/202403052039293.png)


# Quick Start

## 1. Swagger

### 1.1 install

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

### 1.2 init swag

```shell
swag init -g ./cmd/server/main.go
```

> FAQ: [Unknown LeftDelim and RightDelim in swag.Spec](https://github.com/swaggo/swag/issues/1568)

## 2. Run

```shell
go mod tidy
go run ./cmd/server/main.go
```
