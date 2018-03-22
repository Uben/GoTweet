-- drop table if exists users cascade;
-- drop table if exists user_meta;
-- drop table if exists sessions;
-- drop table if exists user_follows;
-- drop table if exists tweets cascade;
-- drop table if exists favorites;

create table if not exists users (
	id			serial primary key,
	name	 	varchar(255),
	email	 	varchar(255),
	username	varchar(255),
	password	varchar(255),
	created_at	timestamp,
	updated_at	timestamp
);

create table if not exists user_meta (
	id 			serial primary key,
	user_id 	integer references users(id),
	description varchar(255),
	url 		varchar(255),
	created_at	timestamp,
	updated_at	timestamp
);

create table if not exists user_follows (
	id 			 serial primary key,
	follower_id  integer references users(id),
	following_id integer references users(id),
	created_at 	 timestamp
);

create table if not exists sessions (
	id			serial primary key,
	user_id		integer references users(id),
	token		varchar(255),
	created_at 	timestamp
);

create table if not exists tweets (
	id				serial primary key,
	user_id			integer references users(id),
	msg				varchar(255),
	name 			varchar(255),
	username 		varchar(255),
	favorite_count	integer default 0,
	retweet_count	integer default 0,
	is_retweet  	boolean default FALSE,
	is_origin_live  boolean default TRUE,
	origin_tweet_id int default 0,
	origin_user_id 	int default 0,
	origin_name 	varchar(255),
	origin_username varchar(255),
	created_at		timestamp
);

create table if not exists favorites (
	id 			serial primary key,
	user_id 	integer references users(id),
	tweet_id 	integer references tweets(id),
	created_at	timestamp
);

create table if not exists retweets (
	id 				serial primary key,
	user_id 		integer references users(id),
	tweet_id 		integer references tweets(id),
	created_at		timestamp
);

-- alter table tweets drop column is_origin_live;
-- alter table tweets add column is_origin_live boolean default TRUE;




