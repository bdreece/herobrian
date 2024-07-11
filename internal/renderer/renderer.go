package renderer

import (
	"html/template"
	"io"
	"io/fs"

	"github.com/Masterminds/sprig/v3"
	"github.com/bdreece/herobrian/web"
	"github.com/labstack/echo/v4"
)

type Renderer struct {
	tmpl  *template.Template
	dirfs fs.FS
}

func (r *Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {
    t, err := r.tmpl.Clone()
    if err != nil {
        return err
    }

    _, err = t.ParseFS(r.dirfs, name)
    if err != nil {
        return err
    }

    return t.ExecuteTemplate(w, name, data)
}

func New() (*Renderer, error) {
	tmpl, err := template.New("").
        Funcs(sprig.FuncMap()).
        ParseFS(web.Templates, "templates/*.gotmpl")

    if err != nil {
        return nil, err
    }

    dirfs, err := fs.Sub(web.Templates, "templates")
    if err != nil {
        return nil, err
    }

	return &Renderer{
		tmpl:  tmpl,
		dirfs: dirfs,
	}, nil
}
