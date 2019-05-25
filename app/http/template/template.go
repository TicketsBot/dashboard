package template

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/hoisie/mustache"
)

type Template struct {
	Layout Layout
	Content string
}

func LoadTemplate(layout Layout, name string) Template {
	content, err := utils.ReadFile(fmt.Sprintf("./public/templates/views/%s.mustache", name)); if err != nil {
		panic(err)
	}

	return Template{
		Layout: layout,
		Content: content,
	}
}

func (t *Template) Render(context ...interface{}) string {
	return mustache.RenderInLayout(t.Content, t.Layout.Content, context)
}
