package botstaff

import (
	"context"
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
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
			user, ok := cache.Instance.GetUser(userId)

			data := userData{
				Id: userId,
			}

			if ok {
				data.Username = user.Username
				data.Discriminator = user.Discriminator
			} else {
				data.Username = "Unknown User"
			}

			users[i] = data

			return nil
		})
	}

	_ = group.Wait() // error not possible

	ctx.JSON(200, users)
}
