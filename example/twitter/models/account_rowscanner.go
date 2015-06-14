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
			&m.ID,
			&m.ProfileID,
			&m.Password,
			&m.URL,
			&m.PasswordStatusConstant,
			&m.Salt,
			&m.EmailAddress,
			&m.EmailStatusConstant,
			&m.StatusConstant,
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

	*(dest.(*[]*Account)) = records

	return nil
}
