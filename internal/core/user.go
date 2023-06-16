package core

type User struct {
	Id             int    `json:"-" db:"id"`
	FirstName      string `json:"firstName" db:"first_name" binding:"required"`
	LastName       string `json:"lastName" db:"last_name" binding:"required"`
	Username       string `json:"username" db:"username" binding:"required"`
	PasswordOrHash string `json:"password" binding:"required"`
	Rating         int    `json:"rating" db:"rating"`
	NumPeopleRated int    `json:"num_people_rated" db:"num_people_rated"`
	TgId           int    `json:"tgId" db:"tg_id"`
}
