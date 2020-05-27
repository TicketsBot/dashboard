package utils

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/common/permission"
)

func GetPermissionLevel(guildId, userId uint64) permission.PermissionLevel {
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return permission.Everyone
	}

	// get member
	member, err := botContext.GetGuildMember(guildId, userId)
	if err != nil {
		return permission.Everyone
	}

	return permission.GetPermissionLevel(botContext, member, guildId)
}