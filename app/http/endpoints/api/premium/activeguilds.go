package api

import (
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/common/model"
	"github.com/TicketsBot/common/permission"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
	"net/http"
)

type setActiveGuildsBody struct {
	SelectedGuilds types.UInt64StringSlice `json:"selected_guilds"`
}

func SetActiveGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	var body setActiveGuildsBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorJson(err))
		return
	}

	legacyEntitlement, err := dbclient.Client.LegacyPremiumEntitlements.GetUserTier(ctx, userId, premium.PatreonGracePeriod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	if legacyEntitlement == nil || legacyEntitlement.IsLegacy {
		ctx.JSON(http.StatusBadRequest, utils.ErrorStr("Not a premium user"))
		return
	}

	tx, err := dbclient.Client.BeginTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	defer tx.Rollback(ctx)

	// Validate under the limit
	limit, ok, err := dbclient.Client.MultiServerSkus.GetPermittedServerCount(ctx, tx, legacyEntitlement.SkuId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.ErrorStr("Not a multi-server subscription"))
		return
	}

	if len(body.SelectedGuilds) > limit {
		ctx.JSON(http.StatusBadRequest, utils.ErrorStr("Too many guilds selected"))
		return
	}

	// Validate has admin in each server
	for _, guildId := range body.SelectedGuilds {
		permissionLevel, err := utils.GetPermissionLevel(ctx, guildId, userId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
			return
		}

		if permissionLevel < permission.Admin {
			ctx.JSON(http.StatusForbidden, utils.ErrorStr("Missing permissions in guild %d", guildId))
			return
		}
	}

	existingGuildEntitlements, err := dbclient.Client.LegacyPremiumEntitlementGuilds.ListForUser(ctx, tx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	// Remove entitlements from guilds that are no longer selected
	for _, existingEntitlement := range existingGuildEntitlements {
		if !utils.Contains(body.SelectedGuilds, existingEntitlement.GuildId) {
			if err := dbclient.Client.LegacyPremiumEntitlementGuilds.DeleteByEntitlement(ctx, tx, existingEntitlement.EntitlementId); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
				return
			}

			if err := dbclient.Client.Entitlements.DeleteById(ctx, tx, existingEntitlement.EntitlementId); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
				return
			}
		}
	}

	// Create entitlements for guilds that were not previously selected, but now are
	existingGuildIds := make([]uint64, len(existingGuildEntitlements))
	for i, existingEntitlement := range existingGuildEntitlements {
		existingGuildIds[i] = existingEntitlement.GuildId
	}

	for _, guildId := range body.SelectedGuilds {
		if !utils.Contains(existingGuildIds, guildId) {
			created, err := dbclient.Client.Entitlements.Create(ctx, tx, &guildId, &userId, legacyEntitlement.SkuId, model.EntitlementSourcePatreon, nil)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
				return
			}

			if err := dbclient.Client.LegacyPremiumEntitlementGuilds.Insert(ctx, tx, userId, guildId, created.Id); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
				return
			}
		}
	}

	// Update entitlements for guilds that were previously selected and still are. This may involve recreating the
	// entitlement if the SKU has changed.
	for _, existingEntitlement := range existingGuildEntitlements {
		if utils.Contains(body.SelectedGuilds, existingEntitlement.GuildId) {
			entitlement, err := dbclient.Client.Entitlements.GetById(ctx, tx, existingEntitlement.EntitlementId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
				return
			}

			if entitlement == nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorStr("Entitlement %s not found", existingEntitlement.EntitlementId.String()))
				return
			}

			if entitlement.SkuId == legacyEntitlement.SkuId {
				continue
			} else {
				// If we need to switch the SKU, then delete and recreate the entitlement
				if err := dbclient.Client.LegacyPremiumEntitlementGuilds.DeleteByEntitlement(ctx, tx, existingEntitlement.EntitlementId); err != nil {
					ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
					return
				}

				if err := dbclient.Client.Entitlements.DeleteById(ctx, tx, existingEntitlement.EntitlementId); err != nil {
					ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
					return
				}

				if _, err := dbclient.Client.Entitlements.Create(ctx, tx, &existingEntitlement.GuildId, &userId, legacyEntitlement.SkuId, model.EntitlementSourcePatreon, entitlement.ExpiresAt); err != nil {
					ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
					return
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
