package core

type Trip struct {
	TripId         int    `json:"-" db:"id"`
	AdminId        int    `json:"admin_id" db:"admin_id"`
	IsPassanger    bool   `json:"is_passanger" db:"is_passanger" biding:"required"`
	PlacesMax      int    `json:"places_max" db:"places_max" biding:"required"`
	PlacesTaken    int    `json:"places_taken" db:"places_taken"`
	ChosenDateTime string `json:"chosen_timestamp" db:"chosen_timestamp" biding:"required"`
	FromPoint      int    `json:"from_point" db:"from_point" biding:"required"`
	ToPoint        int    `json:"to_point" db:"to_point" biding:"required"`
	Description    string `json:"description" db:"description"`
}
