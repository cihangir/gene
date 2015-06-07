package models

import "database/sql"

func (a *Account) RowsScan(rows *sql.Rows, dest interface{}) error {
	if rows == nil {
		return nil
	}

	var records []*Account
	for rows.Next() {
		m := NewAccount()
		err := rows.Scan(
			&m.CreatedAt,
			&m.EmailAddress,
			&m.EmailStatusConstant,
			&m.ID,
			&m.Password,
			&m.PasswordStatusConstant,
			&m.ProfileID,
			&m.Salt,
			&m.StatusConstant,
			&m.URL,
			&m.URLName,
		)
		if err != nil {
			return err
		}
		records = append(records, m)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	*(dest.(*[]*Account)) = records

	return nil
}
