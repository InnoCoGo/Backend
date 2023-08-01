package core

type Trip struct {
	TripId          int64  `json:"id" db:"id"`
	AdminId         int64  `json:"admin_id" db:"admin_id"`
	AdminUsername   string `json:"admin_username" db:"admin_username"`
	AdminTgId       int64  `json:"admin_tg_id" db:"admin_tg_id"`
	IsDriver        bool   `json:"is_driver" db:"is_driver" biding:"required"`
	PlacesMax       int    `json:"places_max" db:"places_max" biding:"required"`
	PlacesTaken     int    `json:"places_taken" db:"places_taken"`
	ChosenTimestamp string `json:"chosen_timestamp" db:"chosen_timestamp" biding:"required"`
	FromPoint       int    `json:"from_point" db:"from_point" biding:"required"`
	ToPoint         int    `json:"to_point" db:"to_point" biding:"required"`
	Description     string `json:"description" db:"description"`
	TranslatedDesc  string `json:"translated_desc" db:"translated_desc"`
}

type InputAdjTrips struct {
	CompanionType  string `json:"companion_type"  biding:"required"`
	LeftTimestamp  string `json:"left_timestamp" db:"chosen_timestamp" biding:"required"`
	RightTimestamp string `json:"right_timestamp" db:"chosen_timestamp" biding:"required"`
	FromPoint      int    `json:"from_point" db:"from_point" biding:"required"`
	ToPoint        int    `json:"to_point" db:"to_point" biding:"required"`
}
