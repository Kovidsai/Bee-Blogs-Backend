package models

import "time"

type Comment struct {
	ID        int       `orm:"auto;pk"`
	Content   string    `orm:"type(text)"`
	Post      *Post     `orm:"rel(fk)"` // Foreign key to Post
	Author    *User     `orm:"rel(fk)"` // Foreign key to User
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
}
