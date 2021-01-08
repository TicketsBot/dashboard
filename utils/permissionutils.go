package utils

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/common/permission"
)

func GetPermissionLevel(guildId, userId uint64) (permission.PermissionLevel, error) {
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return permission.Everyone, err
	}

	// do this check here before trying to get the member
	if botContext.IsBotAdmin(userId) {
		return permission.Admin, nil
	}

	// get member
	member, err := botContext.GetGuildMember(guildId, userId)
	if err != nil {
		return permission.Everyone, err
	}

	return permission.GetPermissionLevel(botContext, member, guildId)
}