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
	AccessTokenKey  = "access_token::"
	RefreshTokenKey = "refresh_token::"
	AccessTokenTTL  = 7000
	RefreshTokenTTL = 7000
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessToken(appId string, secret string, cache *cache.RedisCache) (string, error) {
	token, err2 := cache.Get(AccessTokenKey)
	if err2 == nil && token != "" {
		return token, nil
	}
	client := httputils.NewHTTPClient(10 * time.Second)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, secret)
	headers := map[string]string{}
	respStr, err := client.GET(url, headers)
	if err != nil {
		log.Println("GetAccessToken:", err)
		return "", err
	}
	var accessToken AccessToken
	err = json.Unmarshal([]byte(respStr), &accessToken)
	if err != nil {
		log.Println("json.Unmarshal AccessToken:", err)
		return "", err
	}
	err = cache.Set(AccessTokenKey, accessToken.AccessToken, int64(AccessTokenTTL))
	if err != nil {
		log.Println("cache AccessToken:", err)
		return "", err
	}
	return accessToken.AccessToken, nil
}

func GetUserInfo(accessToken string, openId string) (*UserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, openId)
	client := httputils.NewHTTPClient(10 * time.Second)
	headers := map[string]string{}
	respStr, err := client.GET(url, headers)
	if err != nil {
		log.Println("GetUserInfo:", err)
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
