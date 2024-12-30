package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var Conf *Config

type Config struct {
	Pay    Pay    `mapstructure:"pay"`
	Jwt    Jwt    `mapstructure:"jwt"`
	Wx     Wx     `mapstructure:"wx"`
	Auth   Auth   `mapstructure:"auth"`
	Qiniu  Qiniu  `mapstructure:"qiniu"`
	Upload Upload `mapstructure:"upload"`
	Cache  Cache  `mapstructure:"cache"`
	DB     DB     `mapstructure:"db"`
	Server Server `mapstructure:"server"`
}
type Server struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Cros string `mapstructure:"cros"`
}
type DB struct {
	Redis Redis `mapstructure:"redis"`
	Mysql Mysql `mapstructure:"mysql"`
}
type Redis struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"poolSize"`
	IdleTimeout  int    `mapstructure:"idleTimeout"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
}
type Mysql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
}
type Cache struct {
	NeedCache []string `mapstructure:"needCache"`
}
type Upload struct {
	Prefix string `mapstructure:"prefix"`
}
type Qiniu struct {
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Region    string `mapstructure:"region"`
}
type Auth struct {
	Ignores    []string `mapstructure:"ignores"`
	NeedLogins []string `mapstructure:"needLogins"`
}
type Wx struct {
	Gzh       Gzh    `mapstructure:"gzh"`
	Appid     string `mapstructure:"appid"`
	AppSecret string `mapstructure:"appSecret"`
	NotifyUrl string `mapstructure:"notifyUrl"`
}
type Gzh struct {
	Appid     string `mapstructure:"appid"`
	AppSecret string `mapstructure:"appSecret"`
}
type Jwt struct {
	Secret  string `mapstructure:"secret"`
	Expire  int64  `mapstructure:"expire"`
	Refresh int64  `mapstructure:"refresh"`
}
type Pay struct {
	WxPay WxPay `mapstructure:"wxPay"`
}
type WxPay struct {
	AppId       string `mapstructure:"appId"`
	MchId       string `mapstructure:"mchId"`       //商户证书的证书序列号
	MchSerialNo string `mapstructure:"mchSerialNo"` //商户证书的证书序列号
	ApiV3Key    string `mapstructure:"apiV3Key"`    //apiV3Key，商户平台获取
	PrivateKey  string `mapstructure:"privateKey"`  //私钥 apiclient_key.pem 读取后的内容
	AppSecret   string `mapstructure:"appSecret"`
	NotifyUrl   string `mapstructure:"notifyUrl"`
	MchCertPath string `mapstructure:"mchCertPath"`
	MchKeyPath  string `mapstructure:"mchKeyPath"`
}

func Init(confFile string) {
	Conf = new(Config)
	v := viper.New()
	v.SetConfigFile(confFile)
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件被修改了")
		err := v.Unmarshal(&Conf)
		if err != nil {
			panic(fmt.Errorf("Unmarshal change config data,err:%v \n", err))
		}
	})
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件出错,err:%v \n", err))
	}
	//解析
	err = v.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("Unmarshal config data,err:%v \n", err))
	}
}
