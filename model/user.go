package model

import (
	"github.com/lib/pq"
)

type User struct {
	Id        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Password  string         `json:"password" db:"password"`
	Address   pq.StringArray `json:"address" db:"address"`
	RoleId    int            `json:"roleId" db:"role_id"`
	IsAdmin   bool           `json:"isAdmin" db:"is_admin"`
	CreatedBy int            `json:"createdBy db:created_by"`
}
type InputUser struct {
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Password  string         `json:"password" db:"password"`
	Address   pq.StringArray `json:"address" db:"address"`
	RoleId    int            `json:"roleId" db:"role_id"`
	CreatedBy int            `json:"createdBy" db:"created_by"`
}

type OutputUser struct {
	Id        int            `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Address   pq.StringArray `json:"address" db:"address"`
	RoleId    int            `json:"roleId" db:"role_id"`
	IsAdmin   bool           `json:"isAdmin" db:"is_admin"`
	CreatedBy int            `json:"createdBy db:created_by"`
}
