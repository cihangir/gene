// Generated struct for Account.
package models

import (
	"strings"

	"github.com/lann/squirrel"
)

// GenerateCreateSQL generates plain sql for the given Account
func (a *Account) GenerateCreateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(a.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(a.ID) != float64(0) {
		columns = append(columns, "id")
		values = append(values, a.ID)
	}
	if float64(a.ProfileID) != float64(0) {
		columns = append(columns, "profile_id")
		values = append(values, a.ProfileID)
	}
	if a.Password != "" {
		columns = append(columns, "password")
		values = append(values, a.Password)
	}
	if a.URL != "" {
		columns = append(columns, "url")
		values = append(values, a.URL)
	}
	if a.PasswordStatusConstant != "" {
		columns = append(columns, "password_status_constant")
		values = append(values, a.PasswordStatusConstant)
	}
	if a.Salt != "" {
		columns = append(columns, "salt")
		values = append(values, a.Salt)
	}
	if a.EmailAddress != "" {
		columns = append(columns, "email_address")
		values = append(values, a.EmailAddress)
	}
	if a.EmailStatusConstant != "" {
		columns = append(columns, "email_status_constant")
		values = append(values, a.EmailStatusConstant)
	}
	if a.StatusConstant != "" {
		columns = append(columns, "status_constant")
		values = append(values, a.StatusConstant)
	}
	if !a.CreatedAt.IsZero() {
		columns = append(columns, "created_at")
		values = append(values, a.CreatedAt)
	}
	return psql.Columns(columns...).Values(values...).ToSql()
}

// GenerateUpdateSQL generates plain update sql statement for the given Account
func (a *Account) GenerateUpdateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(a.TableName())
	if float64(a.ProfileID) != float64(0) {
		psql = psql.Set("profile_id", a.ProfileID)
	}
	if a.Password != "" {
		psql = psql.Set("password", a.Password)
	}
	if a.URL != "" {
		psql = psql.Set("url", a.URL)
	}
	if a.PasswordStatusConstant != "" {
		psql = psql.Set("password_status_constant", a.PasswordStatusConstant)
	}
	if a.Salt != "" {
		psql = psql.Set("salt", a.Salt)
	}
	if a.EmailAddress != "" {
		psql = psql.Set("email_address", a.EmailAddress)
	}
	if a.EmailStatusConstant != "" {
		psql = psql.Set("email_status_constant", a.EmailStatusConstant)
	}
	if a.StatusConstant != "" {
		psql = psql.Set("status_constant", a.StatusConstant)
	}
	if !a.CreatedAt.IsZero() {
		psql = psql.Set("created_at", a.CreatedAt)
	}
	return psql.Where("id = ?", a.ID).ToSql()
}

// GenerateDeleteSQL generates plain delete sql statement for the given Account
func (a *Account) GenerateDeleteSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(a.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(a.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, a.ID)
	}
	if float64(a.ProfileID) != float64(0) {
		columns = append(columns, "profile_id = ?")
		values = append(values, a.ProfileID)
	}
	if a.Password != "" {
		columns = append(columns, "password = ?")
		values = append(values, a.Password)
	}
	if a.URL != "" {
		columns = append(columns, "url = ?")
		values = append(values, a.URL)
	}
	if a.PasswordStatusConstant != "" {
		columns = append(columns, "password_status_constant = ?")
		values = append(values, a.PasswordStatusConstant)
	}
	if a.Salt != "" {
		columns = append(columns, "salt = ?")
		values = append(values, a.Salt)
	}
	if a.EmailAddress != "" {
		columns = append(columns, "email_address = ?")
		values = append(values, a.EmailAddress)
	}
	if a.EmailStatusConstant != "" {
		columns = append(columns, "email_status_constant = ?")
		values = append(values, a.EmailStatusConstant)
	}
	if a.StatusConstant != "" {
		columns = append(columns, "status_constant = ?")
		values = append(values, a.StatusConstant)
	}
	if !a.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, a.CreatedAt)
	}
	if len(columns) != 0 {
		psql = psql.Where(strings.Join(columns, " AND "), values...)
	}
	return psql.ToSql()
}

// GenerateSelectSQL generates plain select sql statement for the given Account
func (a *Account) GenerateSelectSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From(a.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(a.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, a.ID)
	}
	if float64(a.ProfileID) != float64(0) {
		columns = append(columns, "profile_id = ?")
		values = append(values, a.ProfileID)
	}
	if a.Password != "" {
		columns = append(columns, "password = ?")
		values = append(values, a.Password)
	}
	if a.URL != "" {
		columns = append(columns, "url = ?")
		values = append(values, a.URL)
	}
	if a.PasswordStatusConstant != "" {
		columns = append(columns, "password_status_constant = ?")
		values = append(values, a.PasswordStatusConstant)
	}
	if a.Salt != "" {
		columns = append(columns, "salt = ?")
		values = append(values, a.Salt)
	}
	if a.EmailAddress != "" {
		columns = append(columns, "email_address = ?")
		values = append(values, a.EmailAddress)
	}
	if a.EmailStatusConstant != "" {
		columns = append(columns, "email_status_constant = ?")
		values = append(values, a.EmailStatusConstant)
	}
	if a.StatusConstant != "" {
		columns = append(columns, "status_constant = ?")
		values = append(values, a.StatusConstant)
	}
	if !a.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, a.CreatedAt)
	}
	if len(columns) != 0 {
		psql = psql.Where(strings.Join(columns, " AND "), values...)
	}
	return psql.ToSql()
}

// TableName returns the table name for Account
func (a *Account) TableName() string {
	return "account.account"
}
