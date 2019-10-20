package table

import "github.com/TicketsBot/GoPanel/database"

type Ticket struct {
	Uuid string `gorm:"column:UUID;type:varchar(36);primary_key"`
	TicketId int `gorm:"column:ID"`
	Guild int64 `gorm:"column:GUILDID"`
	Channel int64 `gorm:"column:CHANNELID"`
	Owner int64 `gorm:"column:OWNERID"`
	Members string `gorm:"column:MEMBERS;type:text"`
	IsOpen bool `gorm:"column:OPEN"`
	OpenTime int64 `gorm:"column:OPENTIME"`
}

func (Ticket) TableName() string {
	return "tickets"
}

func GetTickets(guild int64) []Ticket {
	var tickets []Ticket
	database.Database.Where(&Ticket{Guild: guild}).Order("ID asc").Find(&tickets)
	return tickets
}

func GetOpenTickets(guild int64) []Ticket {
	var tickets []Ticket
	database.Database.Where(&Ticket{Guild: guild, IsOpen: true}).Order("ID asc").Find(&tickets)
	return tickets
}

func GetTicket(uuid string, ch chan Ticket) {
	var ticket Ticket
	database.Database.Where(&Ticket{Uuid: uuid}).First(&ticket)
	ch <- ticket
}

func GetTicketById(guild int64, id int, ch chan Ticket) {
	var ticket Ticket
	database.Database.Where(&Ticket{Guild: guild, TicketId: id}).First(&ticket)
	ch <- ticket
}
