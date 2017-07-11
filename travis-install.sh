#!/bin/bash

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
sudo tee /etc/apt/sources.list.d/tarantool_1_7.list <<- EOF
deb http://download.tarantool.org/tarantool/1.7/ubuntu/ $release main
deb-src http://download.tarantool.org/tarantool/1.7/ubuntu/ $release main
EOF

# install
sudo apt-get update
sudo apt-get -y install tarantool

mkdir logs
mkdir sessions

tarantool init_tarantool.lua 2>./logs/tarantool-stderr.log  >./logs/tarantool-stdout.log &
