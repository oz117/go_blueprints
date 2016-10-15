package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func UploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userId")
	file, header, err := req.FormFile("avatarFile")

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	if err := ioutil.WriteFile(filename, data, 0777); err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "Successful")
}
