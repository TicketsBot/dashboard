package template

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils"
)

type Layout struct {
	Name string
	Content string
}

func LoadLayout(name string) Layout {
	content, err := utils.ReadFile(fmt.Sprintf("./public/templates/layouts/%s.mustache", name)); if err != nil {
		panic(err)
	}

	return Layout{
		Name: name,
		Content: content,
	}
}
