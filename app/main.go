package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type RespImpl struct {
	Data string
}

var (
	HOME            string
	WERF_SECRET_KEY string
)

func init() {
	HOME = "/tmp/werf"
	WERF_SECRET_KEY = os.Getenv("WERF_SECRET_KEY")
}

func main() {
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("index.html"))

	resp := RespImpl{}

	data := r.FormValue("secret")
	if data != "" {
		rawDecodedText, err := base64.StdEncoding.DecodeString(data)
		if err == nil {
			resp.Data = werfEncrypt(string(rawDecodedText))
		} else {
			fmt.Println(err)
		}
	}

	tpl.Execute(w, resp)
}

func werfEncrypt(data string) string {
	var stdOut, stdErr bytes.Buffer

	cmd := exec.Command("werf", "helm", "secret", "encrypt")

	cmd.Env = append(cmd.Env, fmt.Sprintf("WERF_SECRET_KEY=%s", WERF_SECRET_KEY), fmt.Sprintf("HOME=%s", HOME))

	cmd.Stdin = strings.NewReader(data)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		fmt.Println(stdErr.String())
		return "enncrypt error"

	}

	return strings.TrimSuffix(stdOut.String(), "\n")
}
