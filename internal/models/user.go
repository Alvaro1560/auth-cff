package models

import "time"

type User struct {
	ID             string    `json:"id" db:"id" valid:"required,uuid"`
	Username       string    `json:"username" db:"username" valid:"required,stringlength(5|50),matches(^[a-zA-Z0-9_]+$)"`
	CodeStudent    string    `json:"code_student,omitempty" db:"code_student" valid:"required,stringlength(0|255)"`
	Dni            string    `json:"dni,omitempty" db:"dni" valid:"required,stringlength(0|255)"`
	Names          string    `json:"names,omitempty" db:"names" valid:"-"`
	LastnameFather string    `json:"lastname_father,omitempty" db:"lastname_father" valid:"-"`
	LastnameMother string    `json:"lastname_mother,omitempty" db:"lastname_mother" valid:"-"`
	Email          string    `json:"email,omitempty" db:"email" valid:"required"`
	Password       string    `json:"password,omitempty" db:"password" valid:"-"`
	IsDelete       int       `json:"is_delete,omitempty" db:"is_delete" valid:"-"`
	IsBlock        *bool     `json:"is_block,omitempty" db:"is_block" valid:"-"`
	SessionID      string    `json:"session_id" bson:"session_id"`
	ClientID       int       `json:"client_id,omitempty" bson:"client_id"`
	Colors         Color     `json:"colors,omitempty" bson:"colors"`
	RealIP         string    `json:"real_ip,omitempty" bson:"real_ip"`
	Status         int       `json:"status,omitempty" db:"status" valid:"-"`
	UserId         string    `json:"user_id,omitempty" db:"user_id"`
	Roles          []*string `json:"roles,omitempty" bson:"roles"`
	DocTypes       []*int    `json:"doc_types,omitempty" bson:"doc_types"`
	Projects       []*string `json:"projects,omitempty" bson:"projects"`
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type LoggedUsers struct {
	Event     string    `json:"event" bson:"event"`
	HostName  string    `json:"host_name" bson:"host_name"`
	IpRemote  string    `json:"ip_remote" bson:"ip_remote"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type PasswordHistory struct {
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Color struct {
	Primary   string `json:"primary" bson:"primary"`
	Secondary string `json:"secondary" bson:"secondary"`
	Tertiary  string `json:"tertiary" bson:"tertiary"`
}
