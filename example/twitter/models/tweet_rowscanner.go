package models

import "database/sql"

func (t *Tweet) RowsScan(rows *sql.Rows, dest interface{}) error {
	if rows == nil {
		return nil
	}

	var records []*Tweet
	for rows.Next() {
		m := NewTweet()
		err := rows.Scan(
			&m.ID,
			&m.ProfileID,
			&m.Body,
			&m.Location,
			&m.RetweetCount,
			&m.FavouritesCount,
			&m.PossiblySensitive,
			&m.CreatedAt,
		)
		if err != nil {
			return err
		}
		records = append(records, m)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*(dest.(*[]*Tweet)) = records

	return nil
}
