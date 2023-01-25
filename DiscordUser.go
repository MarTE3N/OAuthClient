package OAuthClientDiscord

import (
	"fmt"
	"strconv"
	"time"
)

const DiscordEpoch = 1420070400000

// Discord User

type DiscordUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
	Locale        string `json:"locale"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}

func (user DiscordUser) GetAvatarUrl() string {
	if user.Avatar == "" {
		if user.Discriminator == "" {
			return ""
		}
		t, _ := strconv.Atoi(user.Discriminator)
		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", t%5)
	} else {
		if user.Id == "" {
			return ""
		}
		if user.Avatar[:2] == "a_" {
			return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.gif", user.Id, user.Avatar)
		} else {
			return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", user.Id, user.Avatar)
		}
	}
}

func (user DiscordUser) GetName() string {
	return user.Username + "#" + user.Discriminator
}

func (user DiscordUser) CreatedAt() time.Time {
	if user.Id == "" {
		return time.Time{}
	}

	t, _ := strconv.ParseInt(user.Id, 10, 64)
	return time.Unix((t>>22)+DiscordEpoch/1000, 0)
}

// Discord User Guild

type DiscordGuild struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Icon           string `json:"icon"`
	Owner          bool   `json:"owner"`
	Permissions    int    `json:"permissions"`
	Features       []int  `json:"features"`
	PermissionsNew string `json:"permissions_new"`
}

func (guild DiscordGuild) GetIconUrl() string {
	if guild.Icon == "" {
		return ""
	} else {
		if guild.Id == "" {
			return ""
		}
		return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", guild.Id, guild.Icon)
	}
}

func (guild DiscordGuild) CreatedAt() time.Time {
	if guild.Id == "" {
		return time.Time{}
	}

	t, _ := strconv.ParseInt(guild.Id, 10, 64)
	return time.Unix((t>>22)+DiscordEpoch/1000, 0)
}

func (guild DiscordGuild) HasPermission(permission int) bool {
	return guild.Permissions&permission != 0
}

// Discord User Connection

type DiscordConnection struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Revoked      bool          `json:"revoked"`
	Integrations []interface{} `json:"integrations"`
	Verified     bool          `json:"verified"`
	FriendSync   bool          `json:"friend_sync"`
	ShowActivity bool          `json:"show_activity"`
	Visibility   int           `json:"visibility"`
}
