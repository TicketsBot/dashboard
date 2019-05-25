package http

import "github.com/TicketsBot/GoPanel/app/http/template"

type Template string

const(
	Index Template = "index"
)

var(
	templates = map[string]template.Template{
		Index.ToString(): template.LoadTemplate(Main.GetInstance(), Index.ToString()),
	}
)

func (t Template) ToString() string {
	return string(t)
}

func (t Template) GetInstance() template.Template {
	return templates[t.ToString()]
}

func (t Template) Render(context ...interface{}) string {
	temp := t.GetInstance()
	return temp.Render(context)
}
