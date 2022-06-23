package utils

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/database"
	"github.com/rxdn/gdl/objects/member"
	discordperms "github.com/rxdn/gdl/permission"
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

	// Check staff override
	staffOverride, err := dbclient.Client.StaffOverride.HasActiveOverride(guildId)
	if err != nil {
		return permission.Everyone, err
	}

	// If staff override enabled and the user is bot staff, grant admin permissions
	if staffOverride {
		isBotStaff, err := dbclient.Client.BotStaff.IsStaff(userId)
		if err != nil {
			return permission.Everyone, err
		}

		if isBotStaff {
			return permission.Admin, nil
		}
	}

	// get member
	member, err := botContext.GetGuildMember(guildId, userId)
	if err != nil {
		return permission.Everyone, err
	}

	return permission.GetPermissionLevel(botContext, member, guildId)
}

// TODO: Use this on the ticket list
func HasPermissionToViewTicket(guildId, userId uint64, ticket database.Ticket) (bool, error) {
	// If user opened the ticket, they will always have permission
	if ticket.UserId == userId && ticket.GuildId == guildId {
		return true, nil
	}

	// Admin override
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return false, err
	}

	if botContext.IsBotAdmin(userId) {
		return true, nil
	}

	// Check if server owner
	guild, err := botContext.GetGuild(guildId)
	if err != nil {
		return false, err
	}

	if guild.OwnerId == userId {
		return true, nil
	}

	member, err := botContext.GetGuildMember(guildId, userId)
	if err != nil {
		return false, err
	}

	// Admins should have access to all tickets
	isAdmin, err := dbclient.Client.Permissions.IsAdmin(guildId, userId)
	if err != nil {
		return false, err
	}

	if isAdmin {
		return true, nil
	}

	// TODO: Check in db
	adminRoles, err := dbclient.Client.RolePermissions.GetAdminRoles(guildId)
	if err != nil {
		return false, err
	}

	for _, roleId := range adminRoles {
		if member.HasRole(roleId) {
			return true, nil
		}
	}

	// Check Discord permissions
	hasAdminPermission := permission.HasPermissions(botContext, guildId, member.User.Id, discordperms.Administrator)
	if hasAdminPermission {
		return true, nil
	}

	// If ticket is not from a panel, we can use default team perms
	if ticket.PanelId == nil {
		canView, err := isOnDefaultTeam(guildId, userId, botContext, member)
		if err != nil {
			return false, err
		}

		return canView, nil
	} else {
		panel, err := dbclient.Client.Panel.GetById(*ticket.PanelId)
		if err != nil {
			return false, err
		}

		if panel.WithDefaultTeam {
			canView, err := isOnDefaultTeam(guildId, userId, botContext, member)
			if err != nil {
				return false, err
			}

			if canView {
				return true, nil
			}
		}

		// If panel does not use the default team, or the user is not assigned to it, check support teams
		supportTeams, err := dbclient.Client.PanelTeams.GetTeams(*ticket.PanelId)
		if err != nil {
			return false, err
		}

		if len(supportTeams) > 0 {
			var supportTeamIds []int
			for _, team := range supportTeams {
				supportTeamIds = append(supportTeamIds, team.Id)
			}

			// Check if user is added to support team directly
			isSupport, err := dbclient.Client.SupportTeamMembers.IsSupportSubset(guildId, userId, supportTeamIds)
			if err != nil {
				return false, err
			}

			if isSupport {
				return true, nil
			}

			// Check if user is added to support team via a role
			isSupport, err = dbclient.Client.SupportTeamRoles.IsSupportAnySubset(guildId, member.Roles, supportTeamIds)
			if err != nil {
				return false, err
			}

			if isSupport {
				return true, nil
			}
		}

		return false, nil
	}
}

func isOnDefaultTeam(guildId, userId uint64, ctx botcontext.BotContext, member member.Member) (bool, error) {
	// Admin perms are already checked straight away, so we don't need to check for them here
	// Check user perms for support
	if isSupport, err := dbclient.Client.Permissions.IsSupport(guildId, member.User.Id); err == nil {
		if isSupport {
			return true, nil
		}
	} else {
		return false, err
	}

	// Check DB for support roles
	supportRoles, err := dbclient.Client.RolePermissions.GetSupportRoles(guildId)
	if err != nil {
		return false, err
	}

	for _, supportRoleId := range supportRoles {
		if member.HasRole(supportRoleId) {
			return true, nil
		}
	}

	return false, nil
}
