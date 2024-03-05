#!/bin/bash

{
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./cmd/server/main.go
git config user.name devops
git config suer.email devops@aastar.xyz
git add .
git commit -m "swagger updated"
git push
} || {
set -e
}
