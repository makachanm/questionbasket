package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ShareRangeType int

const (
	PUBLIC ShareRangeType = iota
	HOME
	PRIVATE

	DIRECT
)

type MessageBot interface {
	SendMessage(aid string, isnsfw bool, content string, visiblity ShareRangeType) error
}

type PlatformType int

const (
	MASTODON PlatformType = iota
	MISSKEY
)

func NewBot(t PlatformType, insurl string, token string, service string) MessageBot {
	switch t {
	case MASTODON:
		return newMastodonBot(insurl, token, service)
	case MISSKEY:
		return newMisskeyBot(insurl, token, service)
	}

	return nil
}

type MastodonBot struct {
	service_url  string
	instance_url string
	token        string
}

type mastodonPayload struct {
	Status     string `json:"status"`
	Visibility string `json:"visibility"`
	Spoiler    string `json:"spoiler_text"`
}

func newMastodonBot(insurl string, token string, service string) MastodonBot {
	return MastodonBot{
		instance_url: insurl,
		token:        token,
		service_url:  service,
	}
}

func switchMastodonVisblity(visiblity ShareRangeType) string {
	switch visiblity {
	case PUBLIC:
		return "public"
	case HOME:
		return "unlisted"
	case PRIVATE:
		return "private"
	case DIRECT:
		return "direct"
	}

	return ""
}

func (mb MastodonBot) SendMessage(aid string, isnsfw bool, content string, visiblity ShareRangeType) error {
	message := makeMessage(mb.service_url, aid, content)

	note := mastodonPayload{
		Status:     message,
		Visibility: switchMastodonVisblity(visiblity),
	}

	data, jerr := json.Marshal(note)
	if jerr != nil {
		return jerr
	}

	turl := fmt.Sprintf("https://%s/api/v1/statuses", mb.instance_url)
	furl, uerr := url.Parse(turl)
	if uerr != nil {
		return uerr
	}

	var req http.Request = http.Request{}
	req.Method = http.MethodPost
	req.URL = furl
	req.Body = io.NopCloser(bytes.NewBuffer(data))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mb.token))
	req.Header.Set("Content-Type", "application/json")

	_, err := http.DefaultClient.Do(&req)
	if err != nil {
		return err
	}

	return err
}

func newMisskeyBot(insurl string, token string, service string) MisskeyBot {
	return MisskeyBot{
		instance_url: insurl,
		token:        token,
		service_url:  service,
	}
}

type MisskeyBot struct {
	service_url  string
	instance_url string
	token        string
}

type misskeyPayload struct {
	Status     string `json:"status"`
	Visibility string `json:"visibility"`
	Spoiler    string `json:"spoiler_text"`
}

func (mb MisskeyBot) SendMessage(aid string, isnsfw bool, content string, visiblity ShareRangeType) error {
	return nil
}

func makeLink(service string, aid string) string {
	return fmt.Sprintf("https://%s/answer?id=%s", service, aid)
}

func makeMessage(service string, aid string, content string) string {
	message := strings.Builder{}
	message.WriteString("New Question Arrived:\n")

	runes := []rune(content)
	if len(runes) > 250 {
		cutted := string(runes[:250])
		message.WriteString(cutted)
		message.WriteString("...\n\n")
	} else {
		message.WriteString(content)
		message.WriteString("\n\n")
	}

	message.WriteString(makeLink(service, aid))

	return message.String()
}
