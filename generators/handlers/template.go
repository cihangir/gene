package handlers

var APITemplate = `package {{ToLowerFirst .}}api

// New creates a new local {{.}} handler
func New() *{{.}} { return &{{.}}{} }

// {{.}} is for holding the api functions
type {{.}} struct{}

// generate this for all indexes
// func (m *{{.}}) ById(ctx context.Context, id *int64, res *models.{{.}}) error {
//  return db.MustGetDB(ctx).ById(models.New{{.}}(), id, res)
// }

// generate this for all indexes
// func (m *{{.}}) ByIds(ctx context.Context, ids *[]int64, res *[]*models.{{.}}) error {
//  return db.MustGetDB(ctx).ByIds(models.New{{.}}(), ids, res)
// }

func (m *{{.}}) One(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return db.MustGetDB(ctx).One(models.New{{.}}(), req, res)
}

func (m *{{.}}) Create(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return db.MustGetDB(ctx).Create(models.New{{.}}(), req, res)
}

func (m *{{.}}) Update(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return db.MustGetDB(ctx).Update(models.New{{.}}(), req, req)
}

func (m *{{.}}) Delete(ctx context.Context, req *request.{{.}}, res *models.{{.}}) error {
    return db.MustGetDB(ctx).Delete(models.New{{.}}(), req, req)
}

func (m *{{.}}) Some(ctx context.Context, req *request.Options, res *[]*models.{{.}}) error {
    return db.MustGetDB(ctx).Some(models.New{{.}}(), req, req)
}`

var HandlersTemplate = `package {{ToLowerFirst .}}handlers

func Update(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Delete(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Create(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Get(u *url.URL, h http.Header, _ interface{}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}
`
