package utils

import (
	"encoding/json"
	"fmt"
	redis "github.com/TicketsBot/GoPanel/cache"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils/discord/endpoints/guild"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/robfig/go-cache"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ProxyResponse struct {
	Premium bool
	Tier int
}

var premiumCache = cache.New(10 * time.Minute, 10 * time.Minute)

func IsPremiumGuild(store sessions.Session, guildId uint64, ch chan bool) {
	guildIdRaw := strconv.FormatUint(guildId, 10)

	if premium, ok := premiumCache.Get(guildIdRaw); ok {
		ch<-premium.(bool)
		return
	}

	// First lookup by premium key, then votes, then patreon
	keyLookup := make(chan bool)
	go table.IsPremium(guildId, keyLookup)

	if <-keyLookup {
		if err := premiumCache.Add(guildIdRaw, true, 10 * time.Minute); err != nil {
			fmt.Println(err.Error())
		}

		ch<-true
	} else {
		// Get guild object
		guildChan := make(chan *objects.Guild)
		go redis.Client.GetGuildByID(guildIdRaw, guildChan)
		g := <-guildChan

		ownerIdRaw := ""
		if g == nil {
			var g objects.Guild
			endpoint := guild.GetGuild(int(guildId))

			endpoint.Request(store, nil, nil, &g)

			ownerIdRaw = g.OwnerId
			go redis.Client.StoreGuild(g)
		}

		// TODO: Find a  way to stop people using votes to exploit panels
		// Lookup votes
		/*ownerId, err := strconv.ParseInt(ownerIdRaw, 10, 64); if err != nil {
			fmt.Println(err.Error())
			ch <- false
			return
		}

		hasVoted := make(chan bool)
		go table.HasVoted(ownerId, hasVoted)
		if <-hasVoted {
			ch <- true

			if err := premiumCache.Add(guildIdRaw, true, 10 * time.Minute); err != nil {
				fmt.Println(err.Error())
			}

			return
		}*/

		// Lookup Patreon
		client := &http.Client{
			Timeout: time.Second * 3,
		}

		url := fmt.Sprintf("%s/ispremium?key=%s&id=%s", config.Conf.Bot.PremiumLookupProxyUrl, config.Conf.Bot.PremiumLookupProxyKey, ownerIdRaw)
		req, err := http.NewRequest("GET", url, nil)

		res, err := client.Do(req); if err != nil {
			fmt.Println(err.Error())
			ch<-false
			return
		}
		defer res.Body.Close()

		content, err := ioutil.ReadAll(res.Body); if err != nil {
			fmt.Println(err.Error())
			ch<-false
			return
		}

		var proxyResponse ProxyResponse
		if err = json.Unmarshal(content, &proxyResponse); err != nil {
			fmt.Println(err.Error())
			ch<-false
			return
		}

		if err := premiumCache.Add(guildIdRaw, proxyResponse.Premium, 10 * time.Minute); err != nil {
			fmt.Println(err.Error())
		}
		ch <-proxyResponse.Premium
	}
}