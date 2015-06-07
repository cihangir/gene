package statements

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestStatements(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	sts, err := New().Generate(common.NewContext(), s)
	equals(t, nil, err)
	for _, s := range sts {
		if strings.HasSuffix(s.Path, "profile_statements.go") {
			equals(t, expected, string(s.Content))
		}
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}

const expected = `// Generated struct for models.
package models

func (p *Profile) GenerateCreateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at")
		values = append(values, p.CreatedAt)
	}
	if p.FirstName != "" {
		columns = append(columns, "first_name")
		values = append(values, p.FirstName)
	}
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id")
		values = append(values, p.ID)
	}
	if p.LastName != "" {
		columns = append(columns, "last_name")
		values = append(values, p.LastName)
	}
	if p.Nick != "" {
		columns = append(columns, "nick")
		values = append(values, p.Nick)
	}
	return psql.Columns(columns...).Values(values...).ToSql()
}
func (p *Profile) GenerateUpdateSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update(p.TableName())
	if p.AvatarURL != "" {
		psql = psql.Set("avatar_url", p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		psql = psql.Set("created_at", p.CreatedAt)
	}
	if p.FirstName != "" {
		psql = psql.Set("first_name", p.FirstName)
	}
	if float64(p.ID) != float64(0) {
		psql = psql.Set("id", p.ID)
	}
	if p.LastName != "" {
		psql = psql.Set("last_name", p.LastName)
	}
	if p.Nick != "" {
		psql = psql.Set("nick", p.Nick)
	}
	return psql.Where("id = ?", p.ID).ToSql()
}
func (p *Profile) GenerateDeleteSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Delete(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url = ?")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, p.CreatedAt)
	}
	if p.FirstName != "" {
		columns = append(columns, "first_name = ?")
		values = append(values, p.FirstName)
	}
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, p.ID)
	}
	if p.LastName != "" {
		columns = append(columns, "last_name = ?")
		values = append(values, p.LastName)
	}
	if p.Nick != "" {
		columns = append(columns, "nick = ?")
		values = append(values, p.Nick)
	}
	return psql.Where(strings.Join(columns, " AND "), values...).ToSql()
}
func (p *Profile) GenerateSelectSQL() (string, []interface{}, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From(p.TableName())
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	if p.AvatarURL != "" {
		columns = append(columns, "avatar_url = ?")
		values = append(values, p.AvatarURL)
	}
	if !p.CreatedAt.IsZero() {
		columns = append(columns, "created_at = ?")
		values = append(values, p.CreatedAt)
	}
	if p.FirstName != "" {
		columns = append(columns, "first_name = ?")
		values = append(values, p.FirstName)
	}
	if float64(p.ID) != float64(0) {
		columns = append(columns, "id = ?")
		values = append(values, p.ID)
	}
	if p.LastName != "" {
		columns = append(columns, "last_name = ?")
		values = append(values, p.LastName)
	}
	if p.Nick != "" {
		columns = append(columns, "nick = ?")
		values = append(values, p.Nick)
	}
	return psql.Where(strings.Join(columns, " AND "), values...).ToSql()
}
func (p *Profile) TableName() string {
	return "Profile"
}
`
