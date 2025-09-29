package wx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mszlu521/thunder/cache"
	"github.com/mszlu521/thunder/tools/httputils"
	"log"
	"time"
)

var (
	WebAccessTokenKey  = "web_access_token::"
	WebAccessTokenTTL  = 7000
	WebRefreshTokenKey = "web_refresh_token::"
	WebRefreshTokenTTL = 3600 * 24 * 29
)

type WebAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

func GetAccessTokenByCode(appId string, secret string, code string, cache *cache.RedisCache) (string, string, error) {
	client := httputils.NewHTTPClient(10 * time.Second)
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appId, secret, code)
	headers := map[string]string{}
	respStr, err := client.GET(url, headers)
	if err != nil {
		log.Println("GetAccessToken:", err)
		return "", "", err
	}
	return dealAccessTokenResponse(respStr, cache)
}
func GetAccessTokenByCache(appId string, cache *cache.RedisCache) (string, string, error) {
	client := httputils.NewHTTPClient(10 * time.Second)
	token, err2 := cache.Get(AccessTokenKey)
	if err2 == nil && token != "" {
		return token, "", nil
	}
	refreshToken, err2 := cache.Get(WebRefreshTokenKey)
	if err2 == nil && refreshToken != "" {
		url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s", appId, refreshToken)
		headers := map[string]string{}
		respStr, err := client.GET(url, headers)
		if err != nil {
			log.Println("Get Refresh AccessToken:", err)
			return "", "", err
		}
		return dealAccessTokenResponse(respStr, cache)
	}
	return "", "", errors.New("no accessToken")
}

func dealAccessTokenResponse(respStr string, cache *cache.RedisCache) (string, string, error) {
	var accessToken WebAccessToken
	err := json.Unmarshal([]byte(respStr), &accessToken)
	if err != nil {
		log.Println("json.Unmarshal AccessToken:", err)
		return "", "", err
	}
	if accessToken.Errcode != 0 {
		//有错误
		return "", "", errors.New(accessToken.Errmsg)
	}
	//在调用其他接口的时候 就不需要重新授权了，直接从缓存获取accessToken
	err = cache.Set(WebAccessTokenKey, accessToken.AccessToken, int64(AccessTokenTTL))
	if err != nil {
		log.Println("cache AccessToken:", err)
		return "", "", err
	}
	err = cache.Set(WebRefreshTokenKey, accessToken.RefreshToken, int64(RefreshTokenTTL))
	if err != nil {
		log.Println("cache refresh AccessToken:", err)
		return "", "", err
	}
	return accessToken.AccessToken, accessToken.Openid, nil
}

type UserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
	Errcode    int      `json:"errcode"`
	Errmsg     string   `json:"errmsg"`
}

func GetWebUserInfo(token string, openid string) (*UserInfo, error) {
	client := httputils.NewHTTPClient(10 * time.Second)
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", token, openid)
	headers := map[string]string{}
	respStr, err := client.GET(url, headers)
	if err != nil {
		log.Println("Get Refresh AccessToken:", err)
		return nil, err
	}
	var userInfo UserInfo
	err = json.Unmarshal([]byte(respStr), &userInfo)
	if err != nil {
		log.Println("json.Unmarshal UserInfo:", err)
		return nil, err
	}
	if userInfo.Errcode != 0 {
		log.Println("GetUserInfo:", userInfo.Errmsg)
		return nil, errors.New(userInfo.Errmsg)
	}
	return &userInfo, nil
}
