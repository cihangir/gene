package handlers

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"testing"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/testdata"
)

const expected = `package testapi

// New creates a new local Message handler
func NewMessage() *Message { return &Message{} }

// Message is for holding the api functions
type Message struct{}

// generate this for all indexes
// func (m *Message) ById(ctx context.Context, id *int64, res *models.Message) error {
//  return db.MustGetDB(ctx).ById(models.NewMessage(), id, res)
// }

// generate this for all indexes
// func (m *Message) ByIds(ctx context.Context, ids *[]int64, res *[]*models.Message) error {
//  return db.MustGetDB(ctx).ByIds(models.NewMessage(), ids, res)
// }

func (m *Message) One(ctx context.Context, req *models.Message, res *models.Message) error {
	return db.MustGetDB(ctx).One(models.NewMessage(), req, res)
}

func (m *Message) Create(ctx context.Context, req *models.Message, res *models.Message) error {
	return db.MustGetDB(ctx).Create(models.NewMessage(), req, res)
}

func (m *Message) Update(ctx context.Context, req *models.Message, res *models.Message) error {
	return db.MustGetDB(ctx).Update(models.NewMessage(), req, req)
}

func (m *Message) Delete(ctx context.Context, req *models.Message, res *models.Message) error {
	return db.MustGetDB(ctx).Delete(models.NewMessage(), req, req)
}

func (m *Message) Some(ctx context.Context, req *request.Options, res *[]*models.Message) error {
	return db.MustGetDB(ctx).Some(models.NewMessage(), req, req)
}
`

func TestConstructors(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	a, err := generate("test", s.Title)
	equals(t, nil, err)
	equals(t, expected, string(a))
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}
