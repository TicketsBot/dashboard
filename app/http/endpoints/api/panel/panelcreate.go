package api

import (
	"errors"
	"github.com/TicketsBot/GoPanel/botcontext"
	dbclient "github.com/TicketsBot/GoPanel/database"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/GoPanel/rpc/cache"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/types"
	"github.com/TicketsBot/common/collections"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/rest/request"
	"regexp"
	"strconv"
	"strings"
)

const freePanelLimit = 3

var (
	placeholderPattern = regexp.MustCompile(`%(\w+)%`)
	channelNamePattern = regexp.MustCompile(`^[\w\d-_\x{00a9}\x{00ae}\x{2000}-\x{3300}\x{d83c}\x{d000}-\x{dfff}\x{d83d}\x{d000}-\x{dfff}\x{d83e}\x{d000}-\x{dfff}]+$`)
)

type panelBody struct {
	ChannelId       uint64                `json:"channel_id,string"`
	MessageId       uint64                `json:"message_id,string"`
	Title           string                `json:"title"`
	Content         string                `json:"content"`
	Colour          uint32                `json:"colour"`
	CategoryId      uint64                `json:"category_id,string"`
	Emoji           types.Emoji           `json:"emote"`
	WelcomeMessage  *string               `json:"welcome_message"`
	Mentions        []string              `json:"mentions"`
	WithDefaultTeam bool                  `json:"default_team"`
	Teams           []int                 `json:"teams"`
	ImageUrl        *string               `json:"image_url,omitempty"`
	ThumbnailUrl    *string               `json:"thumbnail_url,omitempty"`
	ButtonStyle     component.ButtonStyle `json:"button_style,string"`
	ButtonLabel     string                `json:"button_label"`
	FormId          *int                  `json:"form_id"`
	NamingScheme    *string               `json:"naming_scheme"`
}

func (p *panelBody) IntoPanelMessageData(customId string, isPremium bool) panelMessageData {
	return panelMessageData{
		ChannelId:    p.ChannelId,
		Title:        p.Title,
		Content:      p.Content,
		CustomId:     customId,
		Colour:       int(p.Colour),
		ImageUrl:     p.ImageUrl,
		ThumbnailUrl: p.ThumbnailUrl,
		Emoji:        p.getEmoji(),
		ButtonStyle:  p.ButtonStyle,
		ButtonLabel:  p.ButtonLabel,
		IsPremium:    isPremium,
	}
}

func CreatePanel(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var data panelBody

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, utils.ErrorJson(err))
		return
	}

	data.MessageId = 0

	// Check panel quota
	premiumTier, err := rpc.PremiumClient.GetTierByGuildId(guildId, false, botContext.Token, botContext.RateLimiter)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	if premiumTier == premium.None {
		panels, err := dbclient.Client.Panel.GetByGuild(guildId)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}

		if len(panels) >= freePanelLimit {
			ctx.AbortWithStatusJSON(402, utils.ErrorStr("You have exceeded your panel quota. Purchase premium to unlock more panels."))
			return
		}
	}

	if !data.doValidations(ctx, guildId) {
		return
	}

	customId := utils.RandString(80)

	messageData := data.IntoPanelMessageData(customId, premiumTier > premium.None)
	msgId, err := messageData.send(&botContext)
	if err != nil {
		var unwrapped request.RestError
		if errors.As(err, &unwrapped) && unwrapped.StatusCode == 403 {
			ctx.JSON(500, utils.ErrorStr("I do not have permission to send messages in the specified channel"))
		} else {
			// TODO: Most appropriate error?
			ctx.JSON(500, utils.ErrorJson(err))
		}

		return
	}

	var emojiId *uint64
	var emojiName *string
	{
		emoji := data.getEmoji()
		if emoji != nil {
			emojiName = &emoji.Name

			if emoji.Id.Value != 0 {
				emojiId = &emoji.Id.Value
			}
		}
	}

	// Store in DB
	panel := database.Panel{
		MessageId:       msgId,
		ChannelId:       data.ChannelId,
		GuildId:         guildId,
		Title:           data.Title,
		Content:         data.Content,
		Colour:          int32(data.Colour),
		TargetCategory:  data.CategoryId,
		EmojiId:         emojiId,
		EmojiName:       emojiName,
		WelcomeMessage:  data.WelcomeMessage,
		WithDefaultTeam: data.WithDefaultTeam,
		CustomId:        customId,
		ImageUrl:        data.ImageUrl,
		ThumbnailUrl:    data.ThumbnailUrl,
		ButtonStyle:     int(data.ButtonStyle),
		ButtonLabel:     data.ButtonLabel,
		FormId:          data.FormId,
		NamingScheme:    data.NamingScheme,
	}

	panelId, err := dbclient.Client.Panel.Create(panel)
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// insert role mention data
	// string is role ID or "user" to mention the ticket opener
	validRoles, err := getRoleHashSet(guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	var roleMentions []uint64
	for _, mention := range data.Mentions {
		if mention == "user" {
			if err = dbclient.Client.PanelUserMention.Set(panelId, true); err != nil {
				ctx.JSON(500, utils.ErrorJson(err))
				return
			}
		} else {
			roleId, err := strconv.ParseUint(mention, 10, 64)
			if err != nil {
				ctx.JSON(400, utils.ErrorStr("Invalid role ID"))
				return
			}

			if validRoles.Contains(roleId) {
				roleMentions = append(roleMentions, roleId)
			}
		}
	}

	if err := dbclient.Client.PanelRoleMentions.Replace(panelId, roleMentions); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	// Already validated, we are safe to insert
	if err := dbclient.Client.PanelTeams.Replace(panelId, data.Teams); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success":  true,
		"panel_id": panelId,
	})
}

var urlRegex = regexp.MustCompile(`^https?://([-a-zA-Z0-9@:%._+~#=]{1,256})\.[a-zA-Z0-9()]{1,63}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)$`)

func (p *panelBody) doValidations(ctx *gin.Context, guildId uint64) bool {
	botContext, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return false // TODO: Log error
	}

	if !p.verifyTitle() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel titles must be between 1 - 80 characters in length",
		})
		return false
	}

	if !p.verifyContent() {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"error":   "Panel content must be between 1 - 4096 characters in length",
		})
		return false
	}

	channels := cache.Instance.GetGuildChannels(guildId)

	if !p.verifyChannel(channels) {
		ctx.JSON(400, utils.ErrorStr("Invalid channel"))
		return false
	}

	if !p.verifyCategory(channels) {
		ctx.JSON(400, utils.ErrorStr("Invalid channel category"))
		return false
	}

	if !p.verifyEmoji(botContext, guildId) {
		ctx.JSON(400, utils.ErrorStr("Invalid emoji"))
		return false
	}

	if !p.verifyWelcomeMessage() {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Welcome message must be blank or between 1 - 4096 characters",
		})
		return false
	}

	if !p.verifyImageUrl() || !p.verifyThumbnailUrl() {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Image URL must be between 1 - 255 characters and a valid URL",
		})
		return false
	}

	if !p.verifyButtonStyle() {
		ctx.JSON(400, gin.H{
			"success": false,
			"error":   "Invalid button style",
		})
		return false
	}

	if !p.verifyButtonLabel() {
		ctx.JSON(400, utils.ErrorStr("Button labels cannot be longer than 80 characters"))
		return false
	}

	{
		valid, err := p.verifyTeams(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return false
		}

		if !valid {
			ctx.JSON(400, utils.ErrorStr("Invalid teams provided"))
			return false
		}
	}

	{
		ok, err := p.verifyFormId(guildId)
		if err != nil {
			ctx.JSON(500, utils.ErrorJson(err))
			return false
		}

		if !ok {
			ctx.JSON(400, utils.ErrorStr("Guild ID for form does not match"))
			return false
		}
	}

	if !p.verifyNamingScheme() {
		ctx.JSON(400, utils.ErrorStr("Invalid naming scheme: ensure that the naming scheme is less than 100 characters and	 the placeholders you have used are valid"))
		return false
	}

	return true
}

func (p *panelBody) verifyTitle() bool {
	if len(p.Title) > 80 {
		return false
	} else if len(p.Title) == 0 { // Fill default
		p.Title = "Open a ticket!"
	}

	return true
}

func (p *panelBody) verifyContent() bool {
	if len(p.Content) > 4096 {
		return false
	} else if len(p.Content) == 0 { // Fill default
		p.Content = "By clicking the button, a ticket will be opened for you."
	}

	return true
}

// Data must be validated before calling this function
func (p *panelBody) getEmoji() *emoji.Emoji {
	return p.Emoji.IntoGdl()
}

func (p *panelBody) verifyChannel(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.ChannelId && ch.Type == channel.ChannelTypeGuildText {
			valid = true
			break
		}
	}

	return valid
}

func (p *panelBody) verifyCategory(channels []channel.Channel) bool {
	var valid bool
	for _, ch := range channels {
		if ch.Id == p.CategoryId && ch.Type == channel.ChannelTypeGuildCategory {
			valid = true
			break
		}
	}

	return valid
}

func (p *panelBody) verifyEmoji(ctx botcontext.BotContext, guildId uint64) bool {
	if p.Emoji.IsCustomEmoji {
		if p.Emoji.Id == nil {
			return false
		}

		emoji, err := ctx.GetGuildEmoji(guildId, *p.Emoji.Id)
		if err != nil { // TODO: Log
			return false
		}

		if emoji.Id.Value == 0 {
			return false
		}

		if emoji.Name != p.Emoji.Name {
			return false
		}

		return true
	} else {
		if len(p.Emoji.Name) == 0 {
			return true
		}

		// Convert from :emoji: to unicode if we need to
		name := strings.TrimSpace(p.Emoji.Name)
		name = strings.Replace(name, ":", "", -1)

		unicode, ok := utils.GetEmoji(name)
		if !ok {
			return false
		}

		p.Emoji.Name = unicode
		return true
	}
}

func (p *panelBody) verifyWelcomeMessage() bool {
	return p.WelcomeMessage == nil || (len(*p.WelcomeMessage) > 0 && len(*p.WelcomeMessage) <= 4096)
}

func (p *panelBody) verifyImageUrl() bool {
	if p.ImageUrl != nil && len(*p.ImageUrl) == 0 {
		p.ImageUrl = nil
	}

	return p.ImageUrl == nil || (len(*p.ImageUrl) <= 255 && urlRegex.MatchString(*p.ImageUrl))
}

func (p *panelBody) verifyThumbnailUrl() bool {
	if p.ThumbnailUrl != nil && len(*p.ThumbnailUrl) == 0 {
		p.ThumbnailUrl = nil
	}

	return p.ThumbnailUrl == nil || (len(*p.ThumbnailUrl) <= 255 && urlRegex.MatchString(*p.ThumbnailUrl))
}

func (p *panelBody) verifyButtonStyle() bool {
	return p.ButtonStyle >= component.ButtonStylePrimary && p.ButtonStyle <= component.ButtonStyleDanger
}

func (p *panelBody) verifyButtonLabel() bool {
	if len(p.ButtonLabel) == 0 {
		p.ButtonLabel = p.Title // Title already checked for 80 char max
	}

	return len(p.ButtonLabel) > 0 && len(p.ButtonLabel) <= 80
}

func (p *panelBody) verifyFormId(guildId uint64) (bool, error) {
	if p.FormId == nil {
		return true, nil
	} else {
		form, ok, err := dbclient.Client.Forms.Get(*p.FormId)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}

		if form.GuildId != guildId {
			return false, nil
		}

		return true, nil
	}
}

func (p *panelBody) verifyTeams(guildId uint64) (bool, error) {
	// Query does not work nicely if there are no teams created in the guild, but if the user submits no teams,
	// then the input is guaranteed to be valid.
	if len(p.Teams) == 0 {
		return true, nil
	}

	return dbclient.Client.SupportTeam.AllTeamsExistForGuild(guildId, p.Teams)
}

func (p *panelBody) verifyNamingScheme() bool {
	if p.NamingScheme == nil {
		return true
	}

	if len(*p.NamingScheme) == 0 {
		p.NamingScheme = nil
		return true
	}

	// Substitute out {} users may use by mistake, spaces for dashes and convert to lowercase
	p.NamingScheme = utils.Ptr(strings.ReplaceAll(*p.NamingScheme, "{", "%"))
	p.NamingScheme = utils.Ptr(strings.ReplaceAll(*p.NamingScheme, "}", "%"))
	p.NamingScheme = utils.Ptr(strings.ReplaceAll(*p.NamingScheme, " ", "-"))
	p.NamingScheme = utils.Ptr(strings.ToLower(*p.NamingScheme))

	if len(*p.NamingScheme) > 100 {
		return false
	}

	// We must remove all placeholders from the string to check whether the rest of the string is legal
	noPlaceholders := *p.NamingScheme

	// Validate placeholders used
	validPlaceholders := []string{"id", "username", "nickname"}
	for _, match := range placeholderPattern.FindAllStringSubmatch(*p.NamingScheme, -1) {
		if len(match) < 2 { // Infallible
			return false
		}

		placeholder := match[1]
		if !utils.Contains(validPlaceholders, placeholder) {
			return false
		}

		noPlaceholders = strings.Replace(noPlaceholders, match[0], "", -1) // match[0] = "%placeholder%"
	}

	return channelNamePattern.MatchString(noPlaceholders)
}

func getRoleHashSet(guildId uint64) (*collections.Set[uint64], error) {
	ctx, err := botcontext.ContextForGuild(guildId)
	if err != nil {
		return nil, err
	}

	roles, err := ctx.GetGuildRoles(guildId)
	if err != nil {
		return nil, err
	}

	set := collections.NewSet[uint64]()

	for _, role := range roles {
		set.Add(role.Id)
	}

	return set, nil
}
