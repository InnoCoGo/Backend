package core

type User struct {
	Id             int    `json:"-" db:"id"`
	FirstName      string `json:"firstName" db:"f_name" binding:"required"`
	LastName       string `json:"lastName" db:"l_name" binding:"required"`
	Username       string `json:"username" db:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Rating         int    `json:"rating" db:"rating"`
	NumPeopleRated int    `json:"num_people_rated" db:"num_people_rated"`
	// TgId     string `json:"tgId" db:"tg_id"`
}
