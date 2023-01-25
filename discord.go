package OAuthClientDiscord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DISCORD SCOPES
const (
	ScopeIdentify    = "identify"
	ScopeEmail       = "email"
	ScopeGuilds      = "guilds"
	ScopeConnections = "connections"
	ScopeGuildsJoin  = "guilds.join"
)

type DiscordOAuth2 struct {
	ClientID     string
	ClientSecret string
	Token        string
	RedirectURL  string
	Scopes       []string
}

type DiscordOAuthUser struct {
	Code         string
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error"`
}

func (d DiscordOAuth2) GetAuthURL() string {
	return "https://discord.com/api/oauth2/authorize?client_id=" + d.ClientID + "&redirect_uri=" + url.QueryEscape(d.RedirectURL) + "&response_type=code&scope=" + strings.Join(d.Scopes, "%20")
}

func (d DiscordOAuth2) RequestAccessToken(DiscordConfig *DiscordOAuthUser) error {
	if DiscordConfig.Code == "" || DiscordConfig.RefreshToken == "" {
		return errors.New("code and refresh_token is empty")
	}

	var reqBody io.Reader

	if DiscordConfig.Code != "" {
		reqBody = strings.NewReader("client_id=" + d.ClientID + "&client_secret=" + d.ClientSecret + "&grant_type=authorization_code&code=" + DiscordConfig.Code + "&redirect_uri=" + url.QueryEscape(d.RedirectURL))
	} else if DiscordConfig.RefreshToken != "" {
		reqBody = strings.NewReader("client_id=" + d.ClientID + "&client_secret=" + d.ClientSecret + "&grant_type=refresh_token&refresh_token=" + DiscordConfig.RefreshToken)
	} else {
		return errors.New("code and refresh_token is empty")
	}

	resp, err := http.Post("https://discord.com/api/oauth2/token", "application/x-www-form-urlencoded", reqBody)

	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&DiscordConfig)

	if err != nil {
		return err
	}

	if DiscordConfig.Error != "" {
		return errors.New(DiscordConfig.Error)
	}

	return nil
}

func (DiscordConfig *DiscordOAuthUser) RequestDiscordUser() (user DiscordUser, err error) {
	if DiscordConfig.AccessToken == "" {
		return DiscordUser{}, errors.New("access_token is empty")
	}

	client := http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", "https://discord.com/api/users/@me", nil)

	if err != nil {
		return DiscordUser{}, err
	}

	req.Header.Add("Authorization", DiscordConfig.TokenType+" "+DiscordConfig.AccessToken)
	var resp *http.Response

	resp, err = client.Do(req)

	if err != nil {
		return DiscordUser{}, err
	}

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return DiscordUser{}, err
	}

	return user, nil
}

func (DiscordConfig *DiscordOAuthUser) RequestDiscordGuilds() (guilds []DiscordGuild, err error) {
	if DiscordConfig.AccessToken == "" {
		return nil, errors.New("access_token is empty")
	}

	client := http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", "https://discord.com/api/users/@me/guilds", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", DiscordConfig.TokenType+" "+DiscordConfig.AccessToken)
	var resp *http.Response

	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&guilds)

	if err != nil {
		return nil, err
	}

	return guilds, nil
}
func (DiscordConfig *DiscordOAuthUser) RequestDiscordConnections() (connections []DiscordConnection, err error) {
	if DiscordConfig.AccessToken == "" {
		return nil, errors.New("access_token is empty")
	}

	client := http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", "https://discord.com/api/users/@me/connections", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", DiscordConfig.TokenType+" "+DiscordConfig.AccessToken)
	var resp *http.Response

	resp, err = client.Do(req)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&connections)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (d DiscordOAuth2) AddUserToGuild(DiscordConfig DiscordOAuthUser, guildId string, userId string, options struct {
	nick  string
	roles []string
	mute  bool
	deaf  bool
}) (err error) {
	if DiscordConfig.AccessToken == "" {
		return errors.New("access_token is empty")
	}
	if guildId == "" {
		return errors.New("guildId is empty")
	}
	if userId == "" {
		return errors.New("userId is empty")
	}
	if d.Token == "" {
		return errors.New("bot token is empty")
	}

	client := http.Client{}
	var req *http.Request

	reqBody := "access_token=" + DiscordConfig.AccessToken

	if options.nick != "" {
		reqBody += "&nick=" + options.nick
	}
	if options.roles != nil {
		reqBody += "&roles=" + strings.Join(options.roles, ",")
	}
	if options.mute {
		reqBody += "&mute=true"
	}
	if options.deaf {
		reqBody += "&deaf=true"
	}

	req, err = http.NewRequest("PUT", "https://discord.com/api/guilds/"+guildId+"/members/"+userId, strings.NewReader(reqBody))

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bot "+d.Token)

	var resp *http.Response

	resp, err = client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New(fmt.Sprint("error with adding user to guild, status code: ", resp.StatusCode))
	}

	return nil
}
