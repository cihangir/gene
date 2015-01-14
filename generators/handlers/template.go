package handlers

// APITemplate holds the template for the handlers
var APITemplate = `package {{ToLower .ModuleName}}api

// New creates a new local {{ToUpperFirst .Name}} handler
func New{{ToUpperFirst .Name}}() *{{ToUpperFirst .Name}} { return &{{ToUpperFirst .Name}}{} }

// {{ToUpperFirst .Name}} is for holding the api functions
type {{ToUpperFirst .Name}} struct{}

// generate this for all indexes
// func (m *{{ToUpperFirst .Name}}) ById(ctx context.Context, id *int64, res *models.{{ToUpperFirst .Name}}) error {
//  return db.MustGetDB(ctx).ById(models.New{{ToUpperFirst .Name}}(), id, res)
// }

// generate this for all indexes
// func (m *{{ToUpperFirst .Name}}) ByIds(ctx context.Context, ids *[]int64, res *[]*models.{{ToUpperFirst .Name}}) error {
//  return db.MustGetDB(ctx).ByIds(models.New{{ToUpperFirst .Name}}(), ids, res)
// }

func (m *{{ToUpperFirst .Name}}) One(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return db.MustGetDB(ctx).One(models.New{{ToUpperFirst .Name}}(), req, res)
}

func (m *{{ToUpperFirst .Name}}) Create(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return db.MustGetDB(ctx).Create(models.New{{ToUpperFirst .Name}}(), req, res)
}

func (m *{{ToUpperFirst .Name}}) Update(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return db.MustGetDB(ctx).Update(models.New{{ToUpperFirst .Name}}(), req, req)
}

func (m *{{ToUpperFirst .Name}}) Delete(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return db.MustGetDB(ctx).Delete(models.New{{ToUpperFirst .Name}}(), req, req)
}

func (m *{{ToUpperFirst .Name}}) Some(ctx context.Context, req *request.Options, res *[]*models.{{ToUpperFirst .Name}}) error {
    return db.MustGetDB(ctx).Some(models.New{{ToUpperFirst .Name}}(), req, req)
}`
