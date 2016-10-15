package client

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

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

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)

	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("authAvatar.GetAvatarURL should return ErrNoAvatarURL when no" +
			"avatarURL is found")
	}
	// md5 yolo@lol.com
	client.userData = map[string]interface{}{"userId": "dcf6951bf7e8e10d994b50de86068994"}
	url, err = gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	if url != "//gravatar.com/avatar/dcf6951bf7e8e10d994b50de86068994" {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned [%s]", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := path.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()
	var fileSystemAvatar FileSystemAvatar
	client := new(client)

	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("authAvatar.GetAvatarURL should return ErrNoAvatarURL when no" +
			"avatarURL is found")
	}
	client.userData = map[string]interface{}{"userId": "abc"}
	url, err = fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("fileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("fileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
