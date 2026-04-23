package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos"`
	DNI       string    `json:"dni"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Rol       string    `json:"rol"`
	Activo    bool      `json:"activo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserDTO struct {
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellidos"`
	DNI       string `json:"dni"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Rol       string `json:"rol"`
}

type UpdateUserDTO struct {
	Nombre    *string `json:"nombre,omitempty"`
	Apellidos *string `json:"apellidos,omitempty"`
	DNI       *string `json:"dni,omitempty"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
	Rol       *string `json:"rol,omitempty"`
}

type ListUsersQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type UserListResponse struct {
	Users      []User `json:"users"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}
