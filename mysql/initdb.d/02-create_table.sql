CREATE TABLE `rapid_go`.`m_users`
(
    id int not null primary key,
	name varchar(255) not null,
    avatar_path varchar(255) not null,
    sex varchar(255) not null,
    enabled boolean not null,
	created_at datetime default CURRENT_TIMESTAMP not null,
	updated_at datetime default CURRENT_TIMESTAMP not null
)