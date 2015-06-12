package main

import (
	"time"
)

type UserCreate struct {
	ID        int64     `json:"id" sql:"SERIAL"`
	Email     string    `json:"email"`
	Password  string
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	ZipCode   int       `json:"zip_code"`
	IsActive  bool      `json:"is_active"`
	CreateDt  time.Time `json:"create_dt"`
	UpdateDt  time.Time `json:"update_dt"`
}

type UserDisplay struct {
	ID        int64     `json:"id" sql:"SERIAL"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	ZipCode   int       `json:"zip_code"`
	IsActive  bool      `json:"is_active"`
	CreateDt  time.Time `json:"create_dt"`
	UpdateDt  time.Time `json:"update_dt"`
}

type UserAuth struct {
	AuthID		int64	`sql:"SERIAL"`
	UserId		int64
	Email		string
	HashPw		string
	CreatedDt	time.Time
	UpdatedDt	time.Time
	LastLoginDt	time.Time
}

type Event struct {
	ID          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Description string    `json:"description"`
	CreatedDt   time.Time `json:"created_dt"`
	StartDt     time.Time `json:"start_dt"`
	EndDt       time.Time `json:"end_dt"`
}
