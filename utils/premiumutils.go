package utils

import (
	"encoding/json"
	"fmt"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/rpc/ratelimit"
	gocache "github.com/robfig/go-cache"
	"github.com/rxdn/gdl/rest"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ProxyResponse struct {
	Premium bool
	Tier int
}

var premiumCache = gocache.New(10 * time.Minute, 10 * time.Minute)

func IsPremiumGuild(guildId uint64, ch chan bool) {
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
		guild, found := cache.Instance.GetGuild(guildId, false)

		if !found {
			var err error
			guild, err = rest.GetGuild(config.Conf.Bot.Token, ratelimit.Ratelimiter, guildId)

			if err == nil { // cache
				go cache.Instance.StoreGuild(guild)
			}
		}

		// TODO: Find a  way to stop people using votes to exploit panels

		// Lookup Patreon
		client := &http.Client{
			Timeout: time.Second * 3,
		}

		url := fmt.Sprintf("%s/ispremium?key=%s&id=%d", config.Conf.Bot.PremiumLookupProxyUrl, config.Conf.Bot.PremiumLookupProxyKey, guild.OwnerId)
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