package renderer

import (
	"bytes"
	"html/template"
	"log"

	"github.com/labstack/echo"
	"github.com/GeertJohan/go.rice"
)

var (
	templateBox *rice.Box
)

func init() {
	templateBox = mustLoadRiceBoxes()
}

// TODO: Need better way to set template path.  Rice will not use absolute path.  Maybe config or embedded rice box?  Embedded can be a pain when restarting often.
func mustLoadRiceBoxes() *rice.Box {
	/* _, filename, _, _ := runtime.Caller(1)
	dir := path.Dir(filename)
	fmt.Printf("filename: %s, dir: %s\n", filename, dir)
	path := path.Join(path.Dir(filename), "../../web/note-monger/templates")
	*/
	return rice.MustFindBox("../../../web/note-monger/templates")
}

func RenderWithLayout(c echo.Context, name string, layoutTemplate string, pageTemplate string, data interface{}) (bytes.Buffer, error) {
	var out bytes.Buffer

	t := template.New(name).Funcs(template.FuncMap{
		"cGetURL": func() string {
			return c.Request().URL.String()
		},
	})

	layout, err := templateBox.String(layoutTemplate)
	if err != nil {
		log.Printf("Could not load layout template '%s': %v", layoutTemplate, err)
		return out, err
	}

	page, err := templateBox.String(pageTemplate)
	if err != nil {
		log.Printf("Could not load page template '%s': %v", pageTemplate, err)
		return out, err
	}

	parsed, err := t.Parse(layout + page)
	if err != nil {
		log.Printf("Could not parse combined templates: '%s' + '%s': %v", layoutTemplate, pageTemplate, err)
		return out, err
	}

	err = parsed.ExecuteTemplate(&out, "base", data)
	if err != nil {
		log.Printf("Could not execute combined templates: '%s' + '%s': %v", layoutTemplate, pageTemplate, err)
		return out, err
	}

	return out, nil
}
