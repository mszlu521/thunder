package wxPay

import (
	"context"
	"errors"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/mszlu521/thunder/config"
)

var instance *WxPay

type WxPay struct {
	Client *wechat.ClientV3
}

type PayBody struct {
	Description string  `json:"description"`
	OutTradeNo  string  `json:"out_trade_no"`
	TimeExpire  string  `json:"time_expire"`
	Amount      float64 `json:"amount"`
	OpenId      string  `json:"openId"`
	ClientIp    string  `json:"clientIp"`
}

func Init(pay *config.WxPay) {
	if pay == nil {
		return
	}
	if pay.MchId == nil {
		panic("WxPay mchId is nil")
	}
	if pay.MchSerialNo == nil {
		panic("WxPay appId is nil")
	}
	if pay.ApiV3Key == nil {
		panic("WxPay apiV3Key is nil")
	}
	if pay.PrivateKey == nil {
		panic("WxPay privateKey is nil")
	}
	client, err := wechat.NewClientV3(*pay.MchId, *pay.MchSerialNo, *pay.ApiV3Key, *pay.PrivateKey)
	if err != nil {
		panic(err)
	}
	err = client.AutoVerifySign()
	if err != nil {
		panic(err)
	}
	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn
	//全局实例化 便于使用
	instance = &WxPay{
		Client: client,
	}
}

func Native(ctx context.Context, body *PayBody) (string, error) {
	conf := config.GetConfig()
	if conf.Pay == nil {
		panic("WxPay config is required")
	}
	pay := conf.Pay.WxPay
	if pay == nil {
		return "", errors.New("WxPay config is required")
	}
	bm := make(gopay.BodyMap)
	bm.Set("appid", pay.GetAppId()).
		Set("mchid", pay.GetMchId()).
		Set("description", body.Description).
		Set("out_trade_no", body.OutTradeNo).
		Set("time_expire", body.TimeExpire).
		Set("notify_url", pay.GetNotifyUrl()).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", body.Amount).
				Set("currency", "CNY")
		})
	rsp, err := instance.Client.V3TransactionNative(ctx, bm)
	if err != nil {
		return "", err
	}
	return rsp.Response.CodeUrl, nil
}

func JsApi(ctx context.Context, body *PayBody) (*wechat.JSAPIPayParams, error) {
	conf := config.GetConfig()
	if conf.Pay == nil {
		panic("WxPay config is required")
	}
	pay := conf.Pay.WxPay
	if pay == nil {
		return nil, errors.New("WxPay config is required")
	}

	bm := make(gopay.BodyMap)
	bm.Set("appid", pay.GetAppId()).
		Set("mchid", pay.GetMchId()).
		Set("description", body.Description).
		Set("out_trade_no", body.OutTradeNo).
		Set("time_expire", body.TimeExpire).
		Set("notify_url", pay.GetNotifyUrl()).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", body.Amount).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", body.OpenId)
		})
	rsp, err := instance.Client.V3TransactionJsapi(ctx, bm)
	if err != nil {
		return nil, err
	}
	if rsp.Code == 0 {
		jsapi, err := instance.Client.PaySignOfJSAPI(pay.GetAppId(), rsp.Response.PrepayId)
		if err != nil {
			return nil, err
		}
		return jsapi, nil
	}
	return nil, errors.New("fail")
}

func H5Pay(ctx context.Context, body *PayBody) (string, error) {
	conf := config.GetConfig()
	if conf.Pay == nil {
		panic("WxPay config is required")
	}
	pay := conf.Pay.WxPay
	if pay == nil {
		return "", errors.New("WxPay config is required")
	}

	bm := make(gopay.BodyMap)
	bm.Set("appid", pay.GetAppId()).
		Set("mchid", pay.GetMchId()).
		Set("description", body.Description).
		Set("out_trade_no", body.OutTradeNo).
		Set("time_expire", body.TimeExpire).
		Set("notify_url", pay.GetNotifyUrl()).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", body.Amount).
				Set("currency", "CNY")
		}).
		SetBodyMap("scene_info", func(bm gopay.BodyMap) {
			bm.Set("payer_client_ip", body.ClientIp).
				SetBodyMap("h5_info", func(bm gopay.BodyMap) {
					bm.Set("type", "Wap").
						Set("app_name", "码神之路")
				})
		})
	rsp, err := instance.Client.V3TransactionH5(ctx, bm)
	if err != nil {
		return "", err
	}
	if rsp.Code == 0 {
		if err != nil {
			return "", err
		}
		return rsp.Response.H5Url, nil
	}
	return "", errors.New("fail")
}