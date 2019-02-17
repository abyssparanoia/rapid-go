CREATE TABLE m_users
{
    id int not null
		primary key,
	name varchar(255) not null comment 'ユーザー名前',
    avatar_path varchar(255) not null comment 'プロフィール画像',
    sex varchar(255) not null comment '性別',
    enabled boolean not null,
	created_at datetime default CURRENT_TIMESTAMP not null,
	updated_at datetime default CURRENT_TIMESTAMP not null
}