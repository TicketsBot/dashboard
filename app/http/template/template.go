package template

import (
	"fmt"
	"github.com/hoisie/mustache"
)

type Layout struct {
	compiled *mustache.Template
}

type Template struct {
	compiled *mustache.Template
	Layout Layout
}

var(
	LayoutMain Layout

	TemplateIndex    Template
	TemplateLogs     Template
	TemplateSettings Template
)

func (t *Template) Render(context ...interface{}) string {
	return t.compiled.RenderInLayout(t.Layout.compiled, context[0])
}

func LoadLayouts() {
	LayoutMain = Layout{
		compiled: loadLayout("main"),
	}
}

func LoadTemplates() {
	TemplateIndex = Template{
		compiled: loadTemplate("index"),
		Layout: LayoutMain,
	}
	TemplateLogs = Template{
		compiled: loadTemplate("logs"),
		Layout: LayoutMain,
	}
	TemplateSettings = Template{
		compiled: loadTemplate("settings"),
		Layout: LayoutMain,
	}
}

func loadLayout(name string) *mustache.Template {
	tmpl, err := mustache.ParseFile(fmt.Sprintf("./public/templates/layouts/%s.mustache", name)); if err != nil {
		panic(err)
	}

	return tmpl
}

func loadTemplate(name string) *mustache.Template {
	tmpl, err := mustache.ParseFile(fmt.Sprintf("./public/templates/views/%s.mustache", name)); if err != nil {
		panic(err)
	}

	return tmpl
}
