drop table users cascade;
drop table sessions;

create table users (
	id 		 serial primary key,
	name 	 varchar(255),
	email 	 varchar(255),
	password varchar(255),
	created_at timestamp,
	updated_at timestamp
);

create table sessions (
	id 		serial primary key,
	user_id integer references users(id),
	token	varchar(255),
	created_at timestamp,
	updated_at timestamp
);

create table tweets (
	id 			serial primary key,
	user_id integer references users(id),
	msg varchar(255),
	created_at timestamp,
	updated_at timestamp
)