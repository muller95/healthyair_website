#!/bin/bash

export HEALTHYAIR_SQL_USER=root
export HEALTHYAIR_SQL_PASSWORD=1
export HEALTHYAIR_SQL_PORT=3308
export HEALTHYAIR_CERTIFICATE_PATH=~/healthyair_website/ssl_cert/server.crt
export HEALTHYAIR_KEY_PATH=~/healthyair_website/ssl_cert/server.key
export HEALTHYAIR_SERVER_PORT=10000

go get "github.com/valyala/fasthttp"
go get "github.com/go-sql-driver/mysql"
go get "github.com/google/uuid"
go get "github.com/tarantool/go-tarantool"