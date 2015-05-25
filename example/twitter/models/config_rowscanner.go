package models

import "database/sql"

func (c *Config) RowsScan(rows *sql.Rows, dest interface{}) error {
	if rows == nil {
		return nil
	}

	var records []*Config
	for rows.Next() {
		m := NewConfig()
		err := rows.Scan(
			&m.Postgres,
		)
		if err != nil {
			return err
		}
		records = append(records, m)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*(dest.(*[]*Config)) = records

	return nil
}
