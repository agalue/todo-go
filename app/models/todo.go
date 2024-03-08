package models

import "time"

type Base struct {
	Title       string `json:"title,omitempty" gorm:"not null"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty" gorm:"default:1"`
}

type Todo struct {
	Base
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Status struct {
	Completed bool `json:"completed,omitempty"`
}
