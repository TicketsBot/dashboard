package botstaff

import (
	"context"
	"errors"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	cache2 "github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/user"
	"golang.org/x/sync/errgroup"
)

type userData struct {
	Id            uint64             `json:"id,string"`
	Username      string             `json:"username"`
	Discriminator user.Discriminator `json:"discriminator"`
}

func ListBotStaffHandler(ctx *gin.Context) {
	staff, err := database.Client.BotStaff.GetAll()
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Get usernames
	group, _ := errgroup.WithContext(context.Background())

	users := make([]userData, len(staff))
	for i, userId := range staff {
		i := i
		userId := userId

		group.Go(func() error {
			data := userData{
				Id: userId,
			}

			user, err := cache.Instance.GetUser(context.Background(), userId)
			if err == nil {
				data.Username = user.Username
				data.Discriminator = user.Discriminator
			} else if errors.Is(err, cache2.ErrNotFound) {
				data.Username = "Unknown User"
			} else {
				return err
			}

			users[i] = data

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, users)
}
