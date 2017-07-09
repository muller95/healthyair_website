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
go get github.com/asaskevich/govalidator

mysql -u root -e 'SOURCE init_db.sql;'

curl http://download.tarantool.org/tarantool/1.7/gpgkey | sudo apt-key add -
release=`lsb_release -c -s`

# install https download transport for APT
sudo apt-get -y install apt-transport-https

# append two lines to a list of source repositories
sudo rm -f /etc/apt/sources.list.d/*tarantool*.list
sudo tee /etc/apt/sources.list.d/tarantool_1_7.list «- EOF
deb http://download.tarantool.org/tarantool/1.7/debian/ $release main
deb-src http://download.tarantool.org/tarantool/1.7/debian/ $release main
EOF

# install
sudo apt-get update
sudo apt-get -y install tarantool

tarantool init_tarantool.lua