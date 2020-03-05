package table

import "github.com/TicketsBot/GoPanel/database"

type TicketWebhook struct {
	Uuid     string `gorm:"column:UUID;type:varchar(36);unique;primary_key"`
	WebhookUrl   string `gorm:"column:CDNURL;type:varchar(200)"`
}

func (TicketWebhook) TableName() string {
	return "webhooks"
}

func (w *TicketWebhook) AddWebhook() {
	database.Database.Create(w)
}

func DeleteWebhookByUuid(uuid string) {
	database.Database.Where(TicketWebhook{Uuid: uuid}).Delete(TicketWebhook{})
}

func GetWebhookByUuid(uuid string, res chan *string) {
	var row TicketWebhook
	database.Database.Where(TicketWebhook{Uuid: uuid}).Take(&row)

	if row.WebhookUrl == "" {
		res <- nil
	} else {
		res <- &row.WebhookUrl
	}
}
