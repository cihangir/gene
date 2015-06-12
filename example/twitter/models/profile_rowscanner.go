package models

import "database/sql"

func (p *Profile) RowsScan(rows *sql.Rows, dest interface{}) error {
	if rows == nil {
		return nil
	}

	var records []*Profile
	for rows.Next() {
		m := NewProfile()
		err := rows.Scan(
			&m.ID,
			&m.ScreenName,
			&m.URL,
			&m.Location,
			&m.Description,
			&m.LinkColor,
			&m.AvatarURL,
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

	*(dest.(*[]*Profile)) = records

	return nil
}
