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
	Layout   Layout
}

var (
	LayoutMain Layout
	LayoutManage Layout

	TemplateIndex    Template
	TemplateLogs     Template
	TemplateSettings Template
	TemplateBlacklist Template
	TemplateTicketList Template
	TemplateTicketView Template
)

func (t *Template) Render(context ...interface{}) string {
	return t.compiled.RenderInLayout(t.Layout.compiled, context[0])
}

func LoadLayouts() {
	LayoutMain = Layout{
		compiled: loadLayout("main"),
	}
	LayoutManage = Layout{
		compiled: loadLayout("manage"),
	}
}

func LoadTemplates() {
	TemplateIndex = Template{
		compiled: loadTemplate("index"),
		Layout:   LayoutMain,
	}
	TemplateLogs = Template{
		compiled: loadTemplate("logs"),
		Layout:   LayoutManage,
	}
	TemplateSettings = Template{
		compiled: loadTemplate("settings"),
		Layout:   LayoutManage,
	}
	TemplateBlacklist = Template{
		compiled: loadTemplate("blacklist"),
		Layout:   LayoutManage,
	}
	TemplateTicketList = Template{
		compiled: loadTemplate("ticketlist"),
		Layout:   LayoutManage,
	}
	TemplateTicketView = Template{
		compiled: loadTemplate("ticketview"),
		Layout:   LayoutManage,
	}
}

func loadLayout(name string) *mustache.Template {
	tmpl, err := mustache.ParseFile(fmt.Sprintf("./public/templates/layouts/%s.mustache", name))
	if err != nil {
		panic(err)
	}

	return tmpl
}

func loadTemplate(name string) *mustache.Template {
	tmpl, err := mustache.ParseFile(fmt.Sprintf("./public/templates/views/%s.mustache", name))
	if err != nil {
		panic(err)
	}

	return tmpl
}
