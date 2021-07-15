package api

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	BaseUri    = "https://i.mi.com/drive/v2"
	GetFolders = BaseUri + "/user/folders/children?ts=%d&parentId=%s&pageNo=%d&limit=%d"
)

var UnauthorizedError = errors.New("登录授权失败")

type Api struct {
	User *User
}

func NewApi(user *User) *Api {
	api := Api{
		User: user,
	}
	return &api
}

func (api *Api) Get(url string) ([]byte, error) {
	result, err := api.User.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if result.StatusCode == http.StatusFound {
		return api.Get(result.Header.Get("Location"))
	}
	if result.StatusCode == http.StatusUnauthorized {
		return nil, UnauthorizedError
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if gjson.Get(string(body), "R").Int() == 401 {
		return nil, UnauthorizedError
	}
	return body, nil
}

func (api *Api) PostForm(url string, values url.Values) ([]byte, error) {
	result, err := api.User.HttpClient.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	if result.StatusCode == http.StatusUnauthorized {
		return nil, UnauthorizedError
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	if gjson.Get(string(body), "R").Int() == 401 {
		return nil, UnauthorizedError
	}
	return body, nil
}
