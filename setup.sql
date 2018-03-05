drop table users cascade;
drop table user_meta;
drop table sessions;
drop table user_follows;
drop table tweets cascade;
drop table favorites;

create table users (
	id			serial primary key,
	name	 	varchar(255),
	email	 	varchar(255),
	username	varchar(255),
	password	varchar(255),
	created_at	timestamp,
	updated_at	timestamp
);

create table user_meta (
	id 			serial primary key,
	user_id 	integer references users(id),
	description varchar(255),
	url 		varchar(255),
	created_at	timestamp,
	updated_at	timestamp
);

create table user_follows (
	id 			 serial primary key,
	follower_id  integer references users(id),
	following_id integer references users(id),
	created_at 	 timestamp,
	updated_at 	 timestamp
);

create table sessions (
	id			serial primary key,
	user_id		integer references users(id),
	token		varchar(255),
	created_at 	timestamp,
	updated_at 	timestamp
);

create table tweets (
	id			serial primary key,
	user_id		integer references users(id),
	msg			varchar(255),
	created_at	timestamp,
	updated_at	timestamp
);

create table favorites (
	id 			serial primary key,
	user_id 	integer references users(id),
	tweet_id 	integer references tweets(id),
	created_at	timestamp,
	updated_at	timestamp
);