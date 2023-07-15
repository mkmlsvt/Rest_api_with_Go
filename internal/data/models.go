package data

import "database/sql"

type Models struct {
	FootballPlayers FootballPlayerModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		FootballPlayers: FootballPlayerModel{
			DB: db,
		},
	}
}
