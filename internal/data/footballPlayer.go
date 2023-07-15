package data

import (
	"database/sql"
	"errors"
)

type FootballPlayer struct {
	Id       int64  `json:"id"`
	Name     string `json:"title"`
	LastName string `json:"last-name"`
	Value    int32  `json:"value",omitempty`
	Team     string `json:"team",omitempty`
}
type FootballPlayerModel struct {
	DB *sql.DB
}

func (f FootballPlayerModel) Insert(player *FootballPlayer) error {
	query := `INSERT INTO football_players (name, last_name, value, team)
			VALUES ($1, $2, $3, $4)
			RETURNING id`

	args := []interface{}{player.Name, player.LastName, player.Value, player.Team}

	return f.DB.QueryRow(query, args...).Scan(&player.Id)
}

func (f FootballPlayerModel) Get(id int64) (*FootballPlayer, error) {
	if id < 1 {
		return nil, errors.New("record Not found")
	}
	query := `SELECT id,name,last_name,value,team FROM football_players where id = $1`

	var footballPlayer FootballPlayer

	err := f.DB.QueryRow(query, id).Scan(
		&footballPlayer.Id,
		&footballPlayer.Name,
		&footballPlayer.LastName,
		&footballPlayer.Value,
		&footballPlayer.Team)
	if err != nil {
		return nil, errors.New("unavailable record")
	}
	return &footballPlayer, nil
}

func (f FootballPlayerModel) Update(footballPlayer *FootballPlayer) error {
	query := `Update football_players Set name=$1, last_name=$2, value=$3, team=$4
				where id = $5
				RETURNING id`
	args := []interface{}{
		footballPlayer.Name,
		footballPlayer.LastName,
		footballPlayer.Value,
		footballPlayer.Team,
		footballPlayer.Id}

	return f.DB.QueryRow(query, args...).Scan(&footballPlayer.Id)
}

func (f FootballPlayerModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not find")
	}
	query := `DELETE FROM football_players where id = $1`

	results, err := f.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("record not find")
	}
	return nil
}

func (f FootballPlayerModel) GetAll() ([]*FootballPlayer, error) {

	query := `SELECT * FROM football_players Order by id`

	rows, err := f.DB.Query(query)

	if err != nil {

		return nil, err
	}

	var footballplayers []*FootballPlayer

	for rows.Next() {
		var footballplayer FootballPlayer
		err := rows.Scan(&footballplayer.Id, &footballplayer.Name,
			&footballplayer.LastName,
			&footballplayer.Value, &footballplayer.Team)
		if err != nil {
			return nil, err
		}
		footballplayers = append(footballplayers, &footballplayer)
	}
	return footballplayers, nil

}
