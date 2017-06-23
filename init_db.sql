CREATE DATABASE healthyair CHARACTER SET utf8 COLLATE utf8_general_ci;

USE healthyair;

CREATE TABLE users (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name text,
	passwd text
); 

CREATE TABLE stations (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name TEXT,
	user_id INT,
	category_id INT
);

CREATE TABLE measures (
	id INT PRIMARY KEY AUTO_INCREMENT,
	t DOUBLE,
	rh DOUBLE,
	co2 DOUBLE,
	time TIMESTAMP,
	station_id INT
);

CREATE TABLE categories (
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
);

ALTER TABLE stations ADD INDEX (user_id);	
ALTER TABLE users ADD UNIQUE (id);	
ALTER TABLE stations ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE measures ADD INDEX (station_id);	
ALTER TABLE stations ADD UNIQUE (user_id);	
ALTER TABLE measures ADD FOREIGN KEY (station_id) REFERENCES stations (user_id) ON DELETE CASCADE ON UPDATE CASCADE;

