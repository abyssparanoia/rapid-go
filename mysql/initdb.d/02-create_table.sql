CREATE TABLE `rapid_go`.`users`
(
    id int not null primary key,
	name varchar(255) not null,
    sex varchar(255) not null,
    enabled boolean not null,
	created_at int not null,
	updated_at int not null
)