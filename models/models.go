package Models

import (
	"database/sql"
	"time"
)

type User struct {
	Id         int
	Name       string
	Email      string
	Username   string
	Hash       string
	Created_at time.Time
	Updated_at time.Time
}

type UserMeta struct {
	Id          int
	User_id     int
	Description sql.NullString
	Url         sql.NullString
	Created_at  time.Time
	Updated_at  time.Time
}

type Follow struct {
	Id           int
	Follower_id  int
	Following_id int
	Created_at   time.Time
}

type Session struct {
	Id         int
	User_id    int
	Token      string
	Created_at time.Time
}

type Tweet struct {
	Id             int
	User_id        int
	Message        string
	Name           sql.NullString
	Username       sql.NullString
	RCount         int
	FCount         int
	Is_retweet     bool
	Is_origin_live bool
	Otweet_id      int
	Ouser_id       int
	Oname          sql.NullString
	Ousername      sql.NullString
	Created_at     time.Time
}

type Favorite struct {
	Id         int
	User_id    int
	Tweet_id   int
	Created_at time.Time
}

type Retweets struct {
	Id         int
	User_id    int
	Tweet_id   int
	Created_at time.Time
}
