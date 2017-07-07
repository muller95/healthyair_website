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

mysql -u root -e 'DROP DATABASE healthyair;'

mysql -u root -e 'CREATE DATABASE healthyair CHARACTER SET utf8 COLLATE utf8_general_ci;'

mysql -u root -e 'USE healthyair;'

mysql -u root -e 'CREATE TABLE users (
	id INT PRIMARY KEY AUTO_INCREMENT,
	email text,
	passwd text,
	name text
);'

mysql -u root -e 'CREATE TABLE stations (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name TEXT,
	user_id INT,
	category_id INT
);'

mysql -u root -e 'CREATE TABLE measures (
	id INT PRIMARY KEY AUTO_INCREMENT,
	t DOUBLE,
	rh DOUBLE,
	co2 DOUBLE,
	time TIMESTAMP,
	station_id INT
);'

mysql -u root -e 'CREATE TABLE categories (
	id int PRIMARY KEY AUTO_INCREMENT,
	name TEXT,
	t_low_bad DOUBLE,
	t_low_norm DOUBLE,
	t_good DOUBLE,
	t_high_norm DOUBLE,
	t_high_bad DOUBLE,
	rh_low_bad DOUBLE,
	rh_low_norm DOUBLE,
	rh_good DOUBLE,
	rh_high_norm DOUBLE,
	rh_high_bad DOUBLE,
	co2_bad DOUBLE,
	co2_norm DOUBLE,
	co2_good DOUBLE
);'

mysql -u root -e 'ALTER TABLE stations ADD INDEX (user_id);'	
mysql -u root -e 'ALTER TABLE users ADD UNIQUE (id);'
mysql -e 'ALTER TABLE stations ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;'

mysql -u root -e 'ALTER TABLE measures ADD INDEX (station_id);'	
mysql -u root -e 'ALTER TABLE stations ADD UNIQUE (user_id);'	
mysql -u root -e 'ALTER TABLE measures ADD FOREIGN KEY (station_id) REFERENCES stations (user_id) ON DELETE CASCADE ON UPDATE CASCADE;'

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