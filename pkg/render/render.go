package render

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/moorada/OauthProject/pkg/template"
	"github.com/moorada/OauthProject/views"
	"io"
)

type Template struct{}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	tmpl := template.New("", &template.BinData{
		Asset:      views.Asset,
		AssetDir:   views.AssetDir,
		AssetNames: views.AssetNames,
	})

	tpl, err := tmpl.Parse(fmt.Sprintf("views/fe/%s", name))
	if err != nil {
		return err
	}
	return tpl.Execute(w, data)
}

