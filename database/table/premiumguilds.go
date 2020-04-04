package table

import (
	"github.com/TicketsBot/GoPanel/database"
	"strings"
	"time"
)

type PremiumGuilds struct {
	Guild       uint64 `gorm:"column:GUILDID;unique;primary_key"`
	Expiry      int64 `gorm:"column:EXPIRY"`
	User        uint64 `gorm:"column:USERID"`
	ActivatedBy uint64 `gorm:"column:ACTIVATEDBY"`
	Keys        string `gorm:"column:KEYSUSED"`
}

func (PremiumGuilds) TableName() string {
	return "premiumguilds"
}

func IsPremium(guild uint64, ch chan bool) {
	var node PremiumGuilds
	database.Database.Where(PremiumGuilds{Guild: guild}).First(&node)

	if node.Expiry == 0 {
		ch <- false
		return
	}

	current := time.Now().UnixNano() / int64(time.Millisecond)
	ch <- node.Expiry > current
}

func AddPremium(key string, guild, userId uint64, length int64, activatedBy uint64) {
	var expiry int64

	hasPrem := make(chan bool)
	go IsPremium(guild, hasPrem)
	isPremium := <-hasPrem

	if isPremium {
		expiryChan := make(chan int64)
		go GetExpiry(guild, expiryChan)
		currentExpiry := <-expiryChan

		expiry = currentExpiry + length
	} else {
		current := time.Now().UnixNano() / int64(time.Millisecond)
		expiry = current + length
	}

	keysChan := make(chan []string)
	go GetKeysUsed(guild, keysChan)
	keys := <-keysChan
	keys = append(keys, key)
	keysStr := strings.Join(keys, ",")

	var node PremiumGuilds
	database.Database.Where(PremiumGuilds{Guild: guild}).Assign(PremiumGuilds{Expiry: expiry, User: userId, ActivatedBy: activatedBy, Keys: keysStr}).FirstOrCreate(&node)
}

func GetExpiry(guild uint64, ch chan int64) {
	var node PremiumGuilds
	database.Database.Where(PremiumGuilds{Guild: guild}).First(&node)
	ch <- node.Expiry
}

func GetKeysUsed(guild uint64, ch chan []string) {
	var node PremiumGuilds
	database.Database.Where(PremiumGuilds{Guild: guild}).First(&node)
	ch <- strings.Split(node.Keys, ",")
}
