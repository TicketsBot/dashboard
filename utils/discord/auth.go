package discord

import (
	"bytes"
	"encoding/json"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/pasztorpisti/qs"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type (
	TokenData struct {
		ClientId     string `qs:"client_id"`
		ClientSecret string `qs:"client_secret"`
		GrantType    string `qs:"grant_type"`
		Code         string `qs:"code"`
		RedirectUri  string `qs:"redirect_uri"`
		Scope        string `qs:"scope"`
	}

	RefreshData struct {
		ClientId     string `qs:"client_id"`
		ClientSecret string `qs:"client_secret"`
		GrantType    string `qs:"grant_type"`
		RefreshToken string `qs:"refresh_token"`
		RedirectUri  string `qs:"redirect_uri"`
		Scope        string `qs:"scope"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}
)

const TokenEndpoint = "https://discordapp.com/api/oauth2/token"

func AccessToken(code string) (TokenResponse, error) {
	data := TokenData{
		ClientId:     strconv.Itoa(int(config.Conf.Oauth.Id)),
		ClientSecret: config.Conf.Oauth.Secret,
		GrantType:    "authorization_code",
		Code:         code,
		RedirectUri:  config.Conf.Oauth.RedirectUri,
		Scope:        "identify guilds",
	}

	res, err := tokenPost(data)
	if err != nil {
		return TokenResponse{}, err
	}

	var unmarshalled TokenResponse
	if err = json.Unmarshal(res, &unmarshalled); err != nil {
		return TokenResponse{}, err
	}

	return unmarshalled, nil
}

func RefreshToken(refreshToken string) (TokenResponse, error) {
	data := RefreshData{
		ClientId:     strconv.Itoa(int(config.Conf.Oauth.Id)),
		ClientSecret: config.Conf.Oauth.Secret,
		GrantType:    "refresh_token",
		RefreshToken: refreshToken,
		RedirectUri:  config.Conf.Oauth.RedirectUri,
		Scope:        "identify guilds",
	}

	res, err := tokenPost(data)
	if err != nil {
		return TokenResponse{}, err
	}

	var unmarshalled TokenResponse
	if err = json.Unmarshal(res, &unmarshalled); err != nil {
		return TokenResponse{}, err
	}

	return unmarshalled, nil
}

func tokenPost(body ...interface{}) ([]byte, error) {
	str, err := qs.Marshal(body[0])
	if err != nil {
		return nil, err
	}
	encoded := []byte(str)

	req, err := http.NewRequest("POST", TokenEndpoint, bytes.NewBuffer([]byte(encoded)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", string(ApplicationFormUrlEncoded))

	client := &http.Client{}
	client.Timeout = 3 * time.Second

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
