package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/user"
	"strconv"
)

type (
	response struct {
		PageLimit int               `json:"page_limit"`
		Users     []blacklistedUser `json:"users"`
		Roles     []blacklistedRole `json:"roles"`
	}

	blacklistedUser struct {
		UserId        uint64             `json:"id,string"`
		Username      string             `json:"username"`
		Discriminator user.Discriminator `json:"discriminator"`
	}

	blacklistedRole struct {
		RoleId uint64 `json:"id,string"`
		Name   string `json:"name"`
	}
)

const pageLimit = 30

// TODO: Paginate
func GetBlacklistHandler(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := pageLimit * (page - 1)

	blacklistedUsers, err := database.Client.Blacklist.GetBlacklistedUsers(guildId, pageLimit, offset)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	userObjects, err := cache.Instance.GetUsers(blacklistedUsers)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Build struct with user_id, name and discriminator
	users := make([]blacklistedUser, len(blacklistedUsers))
	for i, userId := range blacklistedUsers {
		userData := blacklistedUser{
			UserId: userId,
		}

		user, ok := userObjects[userId]
		if ok {
			userData.Username = user.Username
			userData.Discriminator = user.Discriminator
		}

		users[i] = userData
	}

	blacklistedRoles, err := database.Client.RoleBlacklist.GetBlacklistedRoles(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	roleObjects, err := cache.Instance.GetRoles(guildId, blacklistedRoles)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Build struct with role_id and name
	roles := make([]blacklistedRole, len(blacklistedRoles))
	for i, roleId := range blacklistedRoles {
		roleData := blacklistedRole{
			RoleId: roleId,
		}

		role, ok := roleObjects[roleId]
		if ok {
			roleData.Name = role.Name
		}

		roles[i] = roleData
	}

	ctx.JSON(200, response{
		PageLimit: pageLimit,
		Users:     users,
		Roles:     roles,
	})
}
