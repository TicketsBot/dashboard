package customisation

import (
	"context"
	"fmt"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/worker/bot/customisation"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// UpdateColours TODO: Don't depend on worker
func UpdateColours(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
        ctx.JSON(500, utils.ErrorJson(err))
        return
    }

	// Allow votes
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, true, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if premiumTier < premium.Premium {
		ctx.JSON(402, utils.ErrorStr("You must have premium to customise message appearance"))
		return
	}

	var data map[customisation.Colour]utils.HexColour
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	if len(data) > len(customisation.DefaultColours) {
		ctx.JSON(400, utils.ErrorStr("Invalid colour"))
		return
	}

	for colourCode, hex := range customisation.DefaultColours {
		if _, ok := data[colourCode]; !ok {
			data[colourCode] = utils.HexColour(hex)
		}
	}

	// TODO: Single query
	group, _ := errgroup.WithContext(context.Background())
	for colourCode, hex := range data {
		colourCode := colourCode
		hex := hex

		fmt.Printf("%d: %d\n", colourCode, hex)

		group.Go(func() error {
			fmt.Printf("%d: %d\n", colourCode, hex.Int())
			return dbclient.Client.CustomColours.Set(guildId, colourCode.Int16(), hex.Int())
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
        return
	}

	ctx.JSON(200, utils.SuccessResponse)
}
