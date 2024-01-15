package poeclient

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const TxPoeHost = "poe.game.qq.com"

const getCharactersPath = "/character-window/get-characters"
const getPassiveSkillsPath = "/character-window/get-passive-skills"
const getItemsPath = "/character-window/get-items"

const poeSessIdName = "POESESSID"

var ErrUnauthorized = errors.New("POESESSID已失效，请更新")
var ErrGetCharactersForbidden = errors.New("你查看的用户不存在或已隐藏")
var ErrRateLimit = errors.New("请求过于频繁，请稍后再试")
var ErrUnknown = errors.New("未预期的错误")

type PoeClient struct {
	client              http.Client
	poeUrl              *url.URL
	getCharactersUrl    *url.URL
	getPassiveSkillsUrl *url.URL
	getItemsUrl         *url.URL
}

func NewPoeClient(poeHost string, poeSessId string) (*PoeClient, error) {
	poeUrl, err := url.Parse("https://" + poeHost)
	if err != nil {
		return nil, err
	}
	getCharactersUrl := poeUrl.JoinPath(getCharactersPath)
	getPassiveSkillsUrl := poeUrl.JoinPath(getPassiveSkillsPath)
	getItemsUrl := poeUrl.JoinPath(getItemsPath)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	poeClient := &PoeClient{
		client: http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		poeUrl:              poeUrl,
		getCharactersUrl:    getCharactersUrl,
		getPassiveSkillsUrl: getPassiveSkillsUrl,
		getItemsUrl:         getItemsUrl,
	}

	cookies := []*http.Cookie{{Name: poeSessIdName, Value: poeSessId}}
	poeClient.client.Jar.SetCookies(poeUrl, cookies)

	return poeClient, nil
}

func (c *PoeClient) GetCharacters(accountName, realm string) (string, error) {
	form := url.Values{}
	form.Add("accountName", accountName)
	form.Add("realm", realm)
	resp, err := c.client.PostForm(c.getCharactersUrl.String(), form)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(data), err
	}

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return "", err
	}

	return string(data), err
}

func checkStatusCode(code int) error {
	if code == 401 {
		return ErrUnauthorized
	}
	if code == 403 {
		return ErrGetCharactersForbidden
	}
	if code == 429 {
		return ErrRateLimit
	}
	if code != 200 {
		return ErrUnknown
	}
	return nil
}

func (c *PoeClient) GetItems(accountName, character, realm string) (string, error) {
	form := url.Values{}
	form.Add("accountName", accountName)
	form.Add("character", character)
	form.Add("realm", realm)
	resp, err := c.client.PostForm(c.getItemsUrl.String(), form)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return string(data), err
	}

	return string(data), nil
}

func (c *PoeClient) GetPassiveSkills(accountName, character, realm string) (string, error) {
	form := url.Values{}
	form.Add("accountName", accountName)
	form.Add("character", character)
	form.Add("realm", realm)
	resp, err := c.client.PostForm(c.getPassiveSkillsUrl.String(), form)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return string(data), err
	}

	return string(data), nil
}
