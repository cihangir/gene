package clients

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

const expected = `package testclient

import (
	"github.com/youtube/vitess/go/rpcplus"
	"golang.org/x/net/context"
)

// New creates a new {message test} rpc client
func NewMessage(client *rpcplus.Client) *Message {
	return &Message{
		client: client,
	}
}

// Message is for holding the api functions
type Message struct {
	client *rpcplus.Client
}

// generate this for all indexes
// func (m *Message) ById(ctx context.Context, id *int64, res *models.Message) error {
//   return m.client.Call(ctx, "Message.ById", id, res)
// }

// generate this for all indexes
// func (m *Message) ByIds(ctx context.Context, ids *[]int64, res *[]*models.Message) error {
//   return m.client.Call(ctx, "Message.ByIds", id, res)
// }

func (m *Message) One(ctx context.Context, req *models.Message, res *models.Message) error {
	return m.client.Call(ctx, "Message.One", req, res)
}

func (m *Message) Create(ctx context.Context, req *models.Message, res *models.Message) error {
	return m.client.Call(ctx, "Message.Create", req, res)
}

func (m *Message) Update(ctx context.Context, req *models.Message, res *models.Message) error {
	return m.client.Call(ctx, "Message.Update", req, res)
}

func (m *Message) Delete(ctx context.Context, req *models.Message, res *models.Message) error {
	return m.client.Call(ctx, "Message.Delete", req, res)
}

func (m *Message) Some(ctx context.Context, req *request.Options, res *[]*models.Message) error {
	return m.client.Call(ctx, "Message.Some", req, res)
}
`

func TestConstructors(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	a, err := generate("test", &s)
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
