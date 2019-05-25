package http

import (
	"github.com/TicketsBot/GoPanel/app/http/template"
)

type Layout string

const(
	Main Layout = "main"
)

var(
	layouts = map[string]template.Layout{
		Main.ToString(): template.LoadLayout(Main.ToString()),
	}
)

func (l Layout) ToString() string {
	return string(l)
}

func (l Layout) GetInstance() template.Layout {
	return layouts[l.ToString()]
}
