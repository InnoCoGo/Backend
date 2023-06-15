package core

type Trip struct {
	TripId         int    `json:"-" db:"id"`
	AdminId        int    `json:"adminId" db:"admin_id"`
	IsPassanger    bool   `json:"isPassanger" db:"is_passanger" biding:"required"`
	PlacesMax      int    `json:"placesMax" db:"places_max" biding:"required"`
	PlacesTaken    int    `json:"placesTaken" db:"places_taken"`
	ChosenDateTime string `json:"chosenDateTime" db:"chosen_date_time" biding:"required"`
	FromPoint      int    `json:"fromPoint" db:"from_point" biding:"required"`
	ToPoint        int    `json:"toPoint" db:"to_point" biding:"required"`
	Description    string `json:"description" db:"description" biding:"required"`
}
