// Generated struct for models.
package models

import (
	"strings"

	"github.com/lann/squirrel"
)

// GenerateCreateSQL generates plain sql for the given Tweet
func (t *Tweet) GenerateCreateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(t.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if t.Body != "" {
		columns = append(columns, "body")
		values = append(values, t.Body)
	}
	if !t.CreatedAt.IsZero() {
		columns = append(columns, "created_at")
		values = append(values, t.CreatedAt)
	}
	if float64(t.FavouritesCount) != float64(0) {
		columns = append(columns, "favourites_count")
		values = append(values, t.FavouritesCount)
	}
	if float64(t.ID) != float64(0) {
		columns = append(columns, "id")
		values = append(values, t.ID)
	}
	if t.Location != "" {
		columns = append(columns, "location")
		values = append(values, t.Location)
	}
	if t.PossiblySensitive != false {
		columns = append(columns, "possibly_sensitive")
		values = append(values, t.PossiblySensitive)
	}
	if float64(t.ProfileID) != float64(0) {
		columns = append(columns, "profile_id")
		values = append(values, t.ProfileID)
	}
	if float64(t.RetweetCount) != float64(0) {
		columns = append(columns, "retweet_count")
		values = append(values, t.RetweetCount)
	}
	return psql.Columns(columns...).Values(values...).ToSql()
}

// GenerateUpdateSQL generates plain update sql statement for the given Tweet
func (t *Tweet) GenerateUpdateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(t.TableName())
	if t.Body != "" {
		psql = psql.Set("body", t.Body)
	}
	if !t.CreatedAt.IsZero() {
		psql = psql.Set("created_at", t.CreatedAt)
	}
	if float64(t.FavouritesCount) != float64(0) {
		psql = psql.Set("favourites_count", t.FavouritesCount)
	}
	if float64(t.ID) != float64(0) {
		psql = psql.Set("id", t.ID)
	}
	if t.Location != "" {
		psql = psql.Set("location", t.Location)
	}
	if t.PossiblySensitive != false {
		psql = psql.Set("possibly_sensitive", t.PossiblySensitive)
	}
	if float64(t.ProfileID) != float64(0) {
		psql = psql.Set("profile_id", t.ProfileID)
	}
	if float64(t.RetweetCount) != float64(0) {
		psql = psql.Set("retweet_count", t.RetweetCount)
	}
	return psql.Where("id = ?", t.ID).ToSql()
}

// GenerateDeleteSQL generates plain delete sql statement for the given Tweet
func (t *Tweet) GenerateDeleteSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(t.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if t.Body != "" {
		columns = append(columns, "body = ?")
		values = append(values, t.Body)
	}
	if !t.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, t.CreatedAt)
	}
	if float64(t.FavouritesCount) != float64(0) {
		columns = append(columns, "favourites_count = ?")
		values = append(values, t.FavouritesCount)
	}
	if float64(t.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, t.ID)
	}
	if t.Location != "" {
		columns = append(columns, "location = ?")
		values = append(values, t.Location)
	}
	if t.PossiblySensitive != false {
		columns = append(columns, "possibly_sensitive = ?")
		values = append(values, t.PossiblySensitive)
	}
	if float64(t.ProfileID) != float64(0) {
		columns = append(columns, "profile_id = ?")
		values = append(values, t.ProfileID)
	}
	if float64(t.RetweetCount) != float64(0) {
		columns = append(columns, "retweet_count = ?")
		values = append(values, t.RetweetCount)
	}
	return psql.Where(strings.Join(columns, " AND "), values...).ToSql()
}

// GenerateSelectSQL generates plain select sql statement for the given Tweet
func (t *Tweet) GenerateSelectSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From(t.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if t.Body != "" {
		columns = append(columns, "body = ?")
		values = append(values, t.Body)
	}
	if !t.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, t.CreatedAt)
	}
	if float64(t.FavouritesCount) != float64(0) {
		columns = append(columns, "favourites_count = ?")
		values = append(values, t.FavouritesCount)
	}
	if float64(t.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, t.ID)
	}
	if t.Location != "" {
		columns = append(columns, "location = ?")
		values = append(values, t.Location)
	}
	if t.PossiblySensitive != false {
		columns = append(columns, "possibly_sensitive = ?")
		values = append(values, t.PossiblySensitive)
	}
	if float64(t.ProfileID) != float64(0) {
		columns = append(columns, "profile_id = ?")
		values = append(values, t.ProfileID)
	}
	if float64(t.RetweetCount) != float64(0) {
		columns = append(columns, "retweet_count = ?")
		values = append(values, t.RetweetCount)
	}
	return psql.Where(strings.Join(columns, " AND "), values...).ToSql()
}

// TableName returns the table name for Tweet
func (t *Tweet) TableName() string {
	return "tweet"
}
