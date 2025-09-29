package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var conf = new(Config)

type Config struct {
	Pay    Pay    `mapstructure:"pay"`
	Server Server    `mapstructure:"server"`
	Cache  Cache     `mapstructure:"cache"`
	Upload Upload `mapstructure:"upload"`
	Qiniu  Qiniu     `mapstructure:"qiniu"`
	DB     DB     `mapstructure:"db"`
	Auth   Auth      `mapstructure:"auth"`
	Wx     Wx        `mapstructure:"wx"`
	Jwt    Jwt       `mapstructure:"jwt"`
	Log    LogConfig `mapstructure:"log"`
}
type Jwt struct {
	Secret  string        `mapstructure:"secret"`
	Expire  time.Duration `mapstructure:"expire"`
	Refresh time.Duration `mapstructure:"refresh"`
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
type DB struct {
	Redis    Redis    `mapstructure:"redis"`
	Mysql    Mysql    `mapstructure:"mysql"`
	Postgres Postgres `mapstructure:"postgres"`
}
type Server struct {
	Port         int           `mapstructure:"port"`
	Cros         []string      `mapstructure:"cros"`
	AllowOrigins []string      `mapstructure:"allowOrigins"`
	Mode         string        `mapstructure:"mode"`
	Name         string        `mapstructure:"name"`
	Version      string        `mapstructure:"version"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

type LogConfig struct {
	Level      string    `mapstructure:"level"`
	Format     string    `mapstructure:"format"`
	AddSource  bool      `mapstructure:"addSource"`
	Filename   string    `mapstructure:"filename"`
	MaxSize    int       `mapstructure:"maxSize"`
	MaxAge     int       `mapstructure:"maxAge"`
	MaxBackups int       `mapstructure:"maxBackups"`
	Output     io.Writer `mapstructure:"output"`
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
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	MaxIdleConns int           `mapstructure:"maxIdleConns"`
	PingTimeout  time.Duration `mapstructure:"pingTimeout"`
	MaxOpenConns int           `mapstructure:"maxOpenConns"`
}
type Cache struct {
	NeedCache []string `mapstructure:"needCache"`
	Expire int64 `mapstructure:"expire"` //单位秒
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

// Aliyun 阿里云配置
type Aliyun struct {
	AccessKeyID     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	Bucket          string `mapstructure:"bucket"`
}

type Auth struct {
	IsAuth bool `mapstructure:"isAuth"`
	Ignores    []string `mapstructure:"ignores"`
	NeedLogins []string `mapstructure:"needLogins"`
}
type Wx struct {
	AppId  string `mapstructure:"appId"`
	Secret string `mapstructure:"secret"`
	Token  string `mapstructure:"token"`
	AesKey string `mapstructure:"aesKey"`
}

type Postgres struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	Database     string        `mapstructure:"database"`
	SSLMode      string        `mapstructure:"sslmode"`
	MaxIdleConns int           `mapstructure:"maxIdleConns"`
	PingTimeout  time.Duration `mapstructure:"pingTimeout"`
	MaxOpenConns int           `mapstructure:"maxOpenConns"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/app/etc")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("config unmarshal failed, err:%v", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err = viper.Unmarshal(conf)
		if err != nil {
			log.Fatalf("config unmarshal failed, err:%v", err)
		}
	})
}

// Init 函数负责初始化配置
// 它会解析命令行参数，读取配置文件，并反序列化到 Config 结构体中
func Init() *viper.Viper {
	// 1. 设置命令行参数
	// 我们可以通过 -c 或 --config 来指定配置文件
	var configFile = pflag.StringP("config", "c", "etc/config.yml", "Path to the config file (e.g., etc/config.yml)")
	pflag.Parse()

	// 2. 初始化 Viper
	v := viper.New()
	v.SetDefault("app.name", "MyDefaultAppName")

	v.SetDefault("server.mode", "release")
	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readTimeout", "5s")
	v.SetDefault("server.writeTimeout", "5s")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.addSource", false)

	// 如果命令行指定了配置文件，则使用它
	if *configFile != "" {
		v.SetConfigFile(*configFile)
		log.Printf("Using config file from command line: %s", *configFile)
	} else {
		// 否则，按默认规则查找
		v.AddConfigPath("etc")    // 在etc目录查找
		v.SetConfigName("config") // 默认配置文件名（不带后缀）
		v.SetConfigType("yaml")   // 配置文件类型
		log.Println("Searching for 'config.yml' in the current directory...")
	}

	// 3. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 4. 将配置反序列化到全局的 conf 变量中
	if err := v.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	// 5. 开启配置热加载
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s. Reloading...", e.Name)
		if err := v.Unmarshal(&conf); err != nil {
			log.Printf("Error reloading config: %v", err)
		} else {
			log.Println("Config reloaded successfully.")
			// 这里可以加入回调函数，通知其他模块配置已更新
			// 比如，重新设置日志级别
		}
	})
	return v
}

// GetConfig 返回已加载的配置单例
// 在调用此函数前，必须先调用 Init()
func GetConfig() *Config {
	if conf == nil {
		// 确保即使有人忘记调用 Init，程序也会以明确的方式失败
		panic("config not initialized, please call config.Init() first")
	}
	return conf
}
