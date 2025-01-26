package api

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database"
	database2 "github.com/TicketsBot/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ExportHandler(config config.Config) gin.HandlerFunc {
	// parse ed25519 pem key
	key, err := x509.ParseECPrivateKey([]byte(config.Export.PrivateKey))
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		guildId := c.Keys["guildid"].(uint64)

		panels, err := database.Client.Panel.GetByGuild(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		multiPanels, err := database.Client.MultiPanels.GetByGuild(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		multiPanelTargets := make(map[int][]int)
		for _, multiPanel := range multiPanels {
			subPanels, err := database.Client.MultiPanelTargets.GetPanels(ctx, multiPanel.Id)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			subPanelIds := make([]int, 0, len(subPanels))
			for _, subPanel := range subPanels {
				subPanelIds = append(subPanelIds, subPanel.PanelId)
			}

			multiPanelTargets[multiPanel.Id] = subPanelIds
		}

		teams := make(map[string]any)

		teamMetadata, err := database.Client.SupportTeam.Get(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		for _, team := range teamMetadata {
			users, err := database.Client.SupportTeamMembers.Get(ctx, team.Id)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			roles, err := database.Client.SupportTeamRoles.Get(ctx, team.Id)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			teams[team.Name] = gin.H{
				"id":           team.Id,
				"name":         team.Name,
				"on_call_role": team.OnCallRole,
				"users":        users,
				"roles":        roles,
			}
		}

		formMetadata, err := database.Client.Forms.GetForms(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		questions, err := database.Client.FormInput.GetInputsForGuild(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		forms := make([]map[string]any, 0, len(formMetadata))
		for _, form := range formMetadata {
			formQuestions := make([]database2.FormInput, 0)
			if tmp, ok := questions[form.Id]; ok {
				formQuestions = tmp
			}

			forms = append(forms, gin.H{
				"id":        form.Id,
				"title":     form.Title,
				"custom_id": form.CustomId,
				"questions": formQuestions,
			})
		}

		tags, err := database.Client.Tag.GetByGuild(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		blacklistedUsers, err := database.Client.Blacklist.GetBlacklistedUsers(ctx, guildId, 10_000, 0)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		blacklistedRoles, err := database.Client.RoleBlacklist.GetBlacklistedRoles(ctx, guildId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		res, err := buildZip(gin.H{
			"panels":                 panels,
			"multi_panels":           multiPanels,
			"multi_panel_sub_panels": multiPanelTargets,
			"teams":                  teams,
			"forms":                  forms,
			"tags":                   tags,
			"blacklist": gin.H{
				"users": blacklistedUsers,
				"roles": blacklistedRoles,
			},
		})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Data(http.StatusOK, "application/zip", res)
	}
}

func buildZip(key ed25519.PrivateKey, obj any) ([]byte, error) {
	marshalled, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	hash.Write(marshalled)
	checksum := hex.EncodeToString(hash.Sum(nil))

	key.Sign(rand.Reader, []byte(), crypto.SHA256)

	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	if err := writeToZip(writer, "data.json", marshalled); err != nil {
		return nil, err
	}

	if err := writeToZip(writer, "data.json.sha256sum", []byte(checksum)); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeToZip(writer *zip.Writer, fileName string, data []byte) error {
	fileWriter, err := writer.Create(fileName)
	if err != nil {
		return err
	}

	_, err = fileWriter.Write(data)
	return err
}
