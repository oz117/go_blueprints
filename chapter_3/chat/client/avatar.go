package client

import "errors"

// ErrNoAvatarURL displayed when there is no avatar
var ErrNoAvatarURL = errors.New("Chat: unable to get an avatar URL.")

// Avatar is a type that abstracts the way we obtain the avatar
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar will contain the link to the Avatar url
type AuthAvatar struct{}

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
