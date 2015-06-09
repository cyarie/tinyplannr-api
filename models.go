package main

import (
	"time"
)

type User struct {
	ID				int			`json:"id" sql:"AUTO_INCREMENT"`
	Email			string		`json:"email"`
	CreatedDt		time.Time	`json:"created_dt"`
}

type Event struct {
	ID				int            `json:"id"`
	UserId			int            `json:"user_id"`
	Description		string         `json:"description"`
	CreatedDt		time.Time      `json:"created_dt"`
	StartDt			time.Time      `json:"start_dt"`
	EndDt			time.Time      `json:"end_dt"`
}
