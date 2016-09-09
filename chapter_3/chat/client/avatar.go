package client

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL displayed when there is no avatar
var ErrNoAvatarURL = errors.New("Chat: unable to get an avatar URL.")

// Avatar is a type that abstracts the way we obtain the avatar
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar will contain the link to the Avatar url
type AuthAvatar struct{}

// GravatarAvatar will contain the link given by gravatar
type GravatarAvatar struct{}

// UseAuthAvatar for later
var UseAuthAvatar AuthAvatar

// GetAvatarURL returns the url of the avatart for a specific client
// else returns ErrNoAvatarURL
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL returns the url of the avatar given by gravatar
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
