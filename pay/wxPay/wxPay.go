package wxPay

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"thunder/config"
)

var Instance *WxPay

type WxPay struct {
	Client *wechat.ClientV3
}

func Init(pay config.WxPay) {
	client, err := wechat.NewClientV3(pay.MchId, pay.MchSerialNo, pay.ApiV3Key, pay.PrivateKey)
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
	Instance = &WxPay{
		Client: client,
	}
}
