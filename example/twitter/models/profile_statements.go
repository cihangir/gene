// Generated struct for Profile.
package models

import (
	"strings"

	"github.com/lann/squirrel"
)

// GenerateCreateSQL generates plain sql for the given Profile
func (p *Profile) GenerateCreateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id")
		values = append(values, p.ID)
	}
	if p.ScreenName != "" {
		columns = append(columns, "screen_name")
		values = append(values, p.ScreenName)
	}
	if p.URL != "" {
		columns = append(columns, "url")
		values = append(values, p.URL)
	}
	if p.Location != "" {
		columns = append(columns, "location")
		values = append(values, p.Location)
	}
	if p.Description != "" {
		columns = append(columns, "description")
		values = append(values, p.Description)
	}
	if p.LinkColor != "" {
		columns = append(columns, "link_color")
		values = append(values, p.LinkColor)
	}
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at")
		values = append(values, p.CreatedAt)
	}
	return psql.Columns(columns...).Values(values...).ToSql()
}

// GenerateUpdateSQL generates plain update sql statement for the given Profile
func (p *Profile) GenerateUpdateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(p.TableName())
	if p.ScreenName != "" {
		psql = psql.Set("screen_name", p.ScreenName)
	}
	if p.URL != "" {
		psql = psql.Set("url", p.URL)
	}
	if p.Location != "" {
		psql = psql.Set("location", p.Location)
	}
	if p.Description != "" {
		psql = psql.Set("description", p.Description)
	}
	if p.LinkColor != "" {
		psql = psql.Set("link_color", p.LinkColor)
	}
	if p.AvatarURL != "" {
		psql = psql.Set("avatar_url", p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		psql = psql.Set("created_at", p.CreatedAt)
	}
	return psql.Where("id = ?", p.ID).ToSql()
}

// GenerateDeleteSQL generates plain delete sql statement for the given Profile
func (p *Profile) GenerateDeleteSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, p.ID)
	}
	if p.ScreenName != "" {
		columns = append(columns, "screen_name = ?")
		values = append(values, p.ScreenName)
	}
	if p.URL != "" {
		columns = append(columns, "url = ?")
		values = append(values, p.URL)
	}
	if p.Location != "" {
		columns = append(columns, "location = ?")
		values = append(values, p.Location)
	}
	if p.Description != "" {
		columns = append(columns, "description = ?")
		values = append(values, p.Description)
	}
	if p.LinkColor != "" {
		columns = append(columns, "link_color = ?")
		values = append(values, p.LinkColor)
	}
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url = ?")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, p.CreatedAt)
	}
	if len(columns) != 0 {
		psql = psql.Where(strings.Join(columns, " AND "), values...)
	}
	return psql.ToSql()
}

// GenerateSelectSQL generates plain select sql statement for the given Profile
func (p *Profile) GenerateSelectSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, p.ID)
	}
	if p.ScreenName != "" {
		columns = append(columns, "screen_name = ?")
		values = append(values, p.ScreenName)
	}
	if p.URL != "" {
		columns = append(columns, "url = ?")
		values = append(values, p.URL)
	}
	if p.Location != "" {
		columns = append(columns, "location = ?")
		values = append(values, p.Location)
	}
	if p.Description != "" {
		columns = append(columns, "description = ?")
		values = append(values, p.Description)
	}
	if p.LinkColor != "" {
		columns = append(columns, "link_color = ?")
		values = append(values, p.LinkColor)
	}
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url = ?")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, p.CreatedAt)
	}
	if len(columns) != 0 {
		psql = psql.Where(strings.Join(columns, " AND "), values...)
	}
	return psql.ToSql()
}

// TableName returns the table name for Profile
func (p *Profile) TableName() string {
	return "profile"
}
