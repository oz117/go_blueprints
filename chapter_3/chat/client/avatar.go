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

// GravatarAvatar will contain the link given by gravatar
type GravatarAvatar struct{}

type FileSystemAvatar struct{}

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

// UseGravatarAvatar select gravatar as source for the avatar
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL returns the url of the avatar given by gravatar
func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return "//gravatar.com/avatar/" + userIdStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userId"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return "/avatars/" + userIdStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
