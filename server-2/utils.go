package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Salt struct {
	Salt string `json:"salt"`
}

const (
	errMessageInternalServerError = "Some errors while marshalling"
	errMessageNoDocument          = "No data with given parameters"
	errMessageDuplicationEmail    = "Already exists email"
)

func GetMD5Hash(password, salt string) string {
	hash := md5.New()
	io.WriteString(hash, salt)
	io.WriteString(hash, password)

	return hex.EncodeToString(hash.Sum(nil))
}

func attachSalt() (string, error) {
	url := "http://127.0.0.1:8001/generate-salt"
	res, err := http.PostForm(url, nil)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var salt Salt
	if err := json.Unmarshal(body, &salt); err != nil {
		return "", err
	}
	return salt.Salt, nil
}

func SendValidatorMessage(writer http.ResponseWriter, v Validator) {
	writer.WriteHeader(http.StatusBadRequest)
	output := struct {
		Messages []string `json:"messages"`
	}{
		Messages: v.Errors,
	}
	jsonResp, err := json.Marshal(output)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(errMessageInternalServerError))
	}
	writer.Write(jsonResp)
}

func SendErrorMessage(writer http.ResponseWriter, status int, err error, message string) {
	writer.WriteHeader(status)
	output := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	jsonResp, err := json.Marshal(output)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(errMessageInternalServerError))
		return
	}
	writer.Write(jsonResp)
}
