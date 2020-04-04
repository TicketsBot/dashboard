package table

import "github.com/TicketsBot/GoPanel/database"

type TicketNamingScheme struct {
	Guild        uint64  `gorm:"column:GUILDID;unique;primary_key"`
	NamingScheme string `gorm:"column:NAMINGSCHEME;type:VARCHAR(16)"`
}

type NamingScheme string

const (
	Id       NamingScheme = "id"
	Username NamingScheme = "username"
)

var Schemes = []NamingScheme{Id, Username}

func (TicketNamingScheme) TableName() string {
	return "TicketNamingScheme"
}

func GetTicketNamingScheme(guild uint64, ch chan NamingScheme) {
	var node TicketNamingScheme
	database.Database.Where(TicketNamingScheme{Guild: guild}).First(&node)
	namingScheme := node.NamingScheme

	if namingScheme == "" {
		ch <- Id
	} else {
		ch <- NamingScheme(namingScheme)
	}
}

func SetTicketNamingScheme(guild uint64, scheme NamingScheme) {
	database.Database.Where(&TicketNamingScheme{Guild: guild}).Assign(&TicketNamingScheme{NamingScheme: string(scheme)}).FirstOrCreate(&TicketNamingScheme{})
}
