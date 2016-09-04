package client

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)

	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("authAvatar.GetAvatarURL should return ErrNoAvatarURL when no" +
			"avatarURL is found")
	}
	testURL := "yolo.com/test"
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("authAvatar.GetAvatarURL should return no error when avatarURL has" +
			"been set.")
	}
	if url != testURL {
		t.Error("The url is not correct")
	}
}
