package api

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
	"xiaomi_disk/log"
)

const (
	xiaomi       = "https://account.xiaomi.com"
	securityHome = xiaomi + "/pass2/security/home?userId=%s"
	userLogin    = "https://i.mi.com/api/user/login?followUp=https%3A%2F%2Fi.mi.com%2Fdrive%2Fh5%23%2Fall&_locale=zh_CN&ts="
	longPolling  = "https://c3.lp.account.xiaomi.com/longPolling/loginUrl" +
		"?sid=passport&callback=https%3A%2F%2Faccount.xiaomi.com&serviceParam=&qs=%253Fsid%253Dpassport&_qrsize=240&_="
)

type User struct {
	HttpClient   *http.Client
	IsLogin      bool
	NickName     string
	UserId       string
	ServiceToken string
	DeviceId     string
}

func NewUser() *User {
	var jar, _ = cookiejar.New(nil)
	user := User{
		HttpClient: &http.Client{
			Transport: log.GetHttpTransport(),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: jar,
		},
	}
	return &user
}

// GetQrImage 获取登录二维码
func (r *User) GetQrImage() ([]byte, error) {
	resp, err := r.HttpClient.Get(longPolling + fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := strings.TrimLeft(string(all), "&&&START&&&")
	qrUrl := gjson.Get(result, "qr").String()

	// 获取QR图片
	resp, err = r.HttpClient.Get(qrUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	qr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 开启长轮询，获取扫码结果
	lpUrl := gjson.Get(result, "lp").String()
	go r.LongPolling(lpUrl)

	return qr, nil
}

// LongPolling 处理扫描二维码后的回调
func (r *User) LongPolling(lpUrl string) {
	resp, err := r.HttpClient.Get(lpUrl)
	if err != nil {
		panic(err)
	}

	r.updateCookies(resp.Cookies())

	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var (
		result  = strings.TrimLeft(string(all), "&&&START&&&")
		cookies = []*http.Cookie{{
			Name:   "passInfo",
			Value:  "login-end",
			Domain: "account.xiaomi.com",
		}}
		fields = []string{"passToken", "userId", "cUserId"}
	)
	for _, v := range fields {
		cookies = append(cookies, &http.Cookie{
			Name:   v,
			Value:  gjson.Get(result, v).String(),
			Domain: "account.xiaomi.com",
		})
	}
	r.updateCookies(cookies)
	// 后面有多次跳转
	err = r.location(gjson.Get(result, "location").String())
	if err != nil {
		panic(err)
	}
	r.IsLogin = r.checkIfLoginIn()
}

// SecondLoginCheck 二次安全验证
func (r *User) SecondLoginCheck() error {
	resp, err := r.HttpClient.Get(userLogin + fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	loginUrl := gjson.Get(string(all), "data.loginUrl").String()
	if loginUrl == "" {
		return errors.New("获取loginUrl失败")
	}
	// 访问loginUrl
	resp, err = r.HttpClient.Get(loginUrl)
	if err != nil {
		return err
	}
	location := resp.Header.Get("Location")
	if location == "" {
		return errors.New("获取signUrl失败")
	}
	// 访问signUrl
	resp, err = r.HttpClient.Get(location)
	if err != nil {
		return err
	}
	if resp.Header.Get("Location") != "https://i.mi.com/drive/h5#/all" {
		return errors.New("安全验证失败")
	}
	return nil
}

func (r *User) location(location string) error {
	resp, err := r.HttpClient.Get(location)
	if err != nil {
		return err
	}
	r.updateCookies(resp.Cookies())
	l := resp.Header.Get("Location")
	if l != "" {
		time.Sleep(time.Millisecond * 100)
		return r.location(l)
	}
	return nil
}

func (r *User) checkIfLoginIn() bool {
	resp, err := r.HttpClient.Get(fmt.Sprintf(securityHome, r.UserId))
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	result := strings.TrimLeft(string(all), "&&&START&&&")
	r.NickName = gjson.Get(result, "data.nickName").String()
	return true
}

// 更新Cookies
func (r *User) updateCookies(newCookies []*http.Cookie) {
	parseUrl, err := url.Parse(xiaomi)
	if err != nil {
		return
	}
	jar, _ := cookiejar.New(nil)
	oldCookies := r.HttpClient.Jar.Cookies(parseUrl)
	for _, n := range newCookies {
		var existed = false
		for _, old := range oldCookies {
			if old.Name == n.Name {
				old = n
				existed = true
			}
		}
		if !existed {
			oldCookies = append(oldCookies, n)
		}
		if n.Name == "userId" {
			r.UserId = n.Value
		}
		if n.Name == "serviceToken" {
			r.ServiceToken = n.Value
		}
		if n.Name == "deviceId" {
			r.DeviceId = n.Value
		}
	}
	jar.SetCookies(parseUrl, oldCookies)
	r.HttpClient.Jar = jar
}
