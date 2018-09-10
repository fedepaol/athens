package credentials

import (
	"bytes"
	"testing"
)

func TestParseStdin(t *testing.T) {
	input := `protocol=https
host=example.com
path=foo.git`
	r := bytes.NewBufferString(input)
	params, err := ParseInput(r)

	if err != nil {
		t.Error(err)
	}

	if params.RepoURL != "example.com" {
		t.Errorf("%s != example.com", params.RepoURL)
	}

	if params.RepoProtocol != "https" {
		t.Errorf("%s != https", params.RepoProtocol)
	}
}

func TestJson(t *testing.T) {
	params := pullParams{
		RepoURL:      "github.com",
		RepoProtocol: "https",
	}
	json := `[
		{
			"url" : "github.com/fedepaol",
			"username": "fedepaol",
			"password" : "prova"
		},
		{
			"url" : "github.*",
			"username": "user",
			"password" : "password"
		}
	]`
	j := bytes.NewBufferString(json)

	credentials, err := FromJSON(j, &params)
	if err != nil {
		t.Error(err)
	}
	if credentials.Username != "user" {
		t.Errorf("%s != user", credentials.Username)
	}
}
