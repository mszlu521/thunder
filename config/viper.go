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
	Pay    *Pay       `mapstructure:"pay"`
	Server *Server    `mapstructure:"server"`
	Cache  *Cache     `mapstructure:"cache"`
	Upload *Upload    `mapstructure:"upload"`
	Qiniu  *Qiniu     `mapstructure:"qiniu"`
	DB     *DB        `mapstructure:"db"`
	Auth   *Auth      `mapstructure:"auth"`
	Wx     *Wx        `mapstructure:"wx"`
	Jwt    *Jwt       `mapstructure:"jwt"`
	Email  *Email     `mapstructure:"email"`
	Log    *LogConfig `mapstructure:"log"`
}

type Email struct {
	Host     *string `mapstructure:"host"`
	Port     *int    `mapstructure:"port"`
	Username *string `mapstructure:"username"`
	Password *string `mapstructure:"password"`
	Identity *string `mapstructure:"identity"`
	From     *string `mapstructure:"from"`
	BaseURL  *string `mapstructure:"baseUrl"`
}

func (e *Email) GetHost() string {
	if e == nil || e.Host == nil {
		return "localhost"
	}
	return *e.Host
}

func (e *Email) GetPort() int {
	if e == nil || e.Port == nil {
		return 587
	}
	return *e.Port
}

func (e *Email) GetUsername() string {
	if e == nil || e.Username == nil {
		return ""
	}
	return *e.Username
}

func (e *Email) GetPassword() string {
	if e == nil || e.Password == nil {
		return ""
	}
	return *e.Password
}

func (e *Email) GetIdentity() string {
	if e == nil || e.Identity == nil {
		return ""
	}
	return *e.Identity
}

func (e *Email) GetFrom() string {
	if e == nil || e.From == nil {
		return ""
	}
	return *e.From
}

func (e *Email) GetBaseURL() string {
	if e == nil || e.BaseURL == nil {
		return ""
	}
	return *e.BaseURL
}

type Jwt struct {
	Secret  *string        `mapstructure:"secret"`
	Expire  *time.Duration `mapstructure:"expire"`
	Refresh *time.Duration `mapstructure:"refresh"`
}

func (j *Jwt) GetSecret() string {
	if j == nil || j.Secret == nil {
		return ""
	}
	return *j.Secret
}

func (j *Jwt) GetExpire() time.Duration {
	if j == nil || j.Expire == nil {
		return 24 * time.Hour
	}
	return *j.Expire
}

func (j *Jwt) GetRefresh() time.Duration {
	if j == nil || j.Refresh == nil {
		return 7 * 24 * time.Hour
	}
	return *j.Refresh
}

type Pay struct {
	WxPay *WxPay `mapstructure:"wxPay"`
}

type WxPay struct {
	AppId       *string `mapstructure:"appId"`
	MchId       *string `mapstructure:"mchId"`       //商户证书的证书序列号
	MchSerialNo *string `mapstructure:"mchSerialNo"` //商户证书的证书序列号
	ApiV3Key    *string `mapstructure:"apiV3Key"`    //apiV3Key，商户平台获取
	PrivateKey  *string `mapstructure:"privateKey"`  //私钥 apiclient_key.pem 读取后的内容
	AppSecret   *string `mapstructure:"appSecret"`
	NotifyUrl   *string `mapstructure:"notifyUrl"`
	MchCertPath *string `mapstructure:"mchCertPath"`
	MchKeyPath  *string `mapstructure:"mchKeyPath"`
}

type DB struct {
	Redis    *Redis    `mapstructure:"redis"`
	Mysql    *Mysql    `mapstructure:"mysql"`
	Postgres *Postgres `mapstructure:"postgres"`
}

type Server struct {
	Port         *int           `mapstructure:"port"`
	Cros         []string      `mapstructure:"cros"`
	AllowOrigins []string      `mapstructure:"allowOrigins"`
	Mode         *string        `mapstructure:"mode"`
	Name         *string        `mapstructure:"name"`
	Version      *string        `mapstructure:"version"`
	Host         *string        `mapstructure:"host"`
	ReadTimeout  *time.Duration `mapstructure:"readTimeout"`
	WriteTimeout *time.Duration `mapstructure:"writeTimeout"`
}

type LogConfig struct {
	Level      *string `mapstructure:"level"`
	Format     *string `mapstructure:"format"`
	AddSource  *bool   `mapstructure:"addSource"`
	Filename   *string `mapstructure:"filename"`
	MaxSize    *int    `mapstructure:"maxSize"`
	MaxAge     *int    `mapstructure:"maxAge"`
	MaxBackups *int    `mapstructure:"maxBackups"`
	Output     io.Writer `mapstructure:"output"`
}

func (l *LogConfig) GetLevel() string {
	if l == nil || l.Level == nil {
		return "info"
	}
	return *l.Level
}

func (l *LogConfig) GetFormat() string {
	if l == nil || l.Format == nil {
		return "json"
	}
	return *l.Format
}

func (l *LogConfig) GetAddSource() bool {
	if l == nil || l.AddSource == nil {
		return false
	}
	return *l.AddSource
}

func (l *LogConfig) GetFilename() string {
	if l == nil || l.Filename == nil {
		return ""
	}
	return *l.Filename
}

func (l *LogConfig) GetMaxSize() int {
	if l == nil || l.MaxSize == nil {
		return 100
	}
	return *l.MaxSize
}

func (l *LogConfig) GetMaxAge() int {
	if l == nil || l.MaxAge == nil {
		return 30
	}
	return *l.MaxAge
}

func (l *LogConfig) GetMaxBackups() int {
	if l == nil || l.MaxBackups == nil {
		return 3
	}
	return *l.MaxBackups
}

type Redis struct {
	Addr         *string `mapstructure:"addr"`
	Password     *string `mapstructure:"password"`
	DB           *int    `mapstructure:"db"`
	PoolSize     *int    `mapstructure:"poolSize"`
	IdleTimeout  *int    `mapstructure:"idleTimeout"`
	MaxOpenConns *int    `mapstructure:"maxOpenConns"`
	MaxIdleConns *int    `mapstructure:"maxIdleConns"`
}

type Mysql struct {
	Host         *string        `mapstructure:"host"`
	Port         *int           `mapstructure:"port"`
	User         *string        `mapstructure:"user"`
	Password     *string        `mapstructure:"password"`
	Database     *string        `mapstructure:"database"`
	MaxIdleConns *int           `mapstructure:"maxIdleConns"`
	PingTimeout  *time.Duration `mapstructure:"pingTimeout"`
	MaxOpenConns *int           `mapstructure:"maxOpenConns"`
}

type Cache struct {
	NeedCache []string `mapstructure:"needCache"`
	Expire *int64 `mapstructure:"expire"` //单位秒
}

func (c *Cache) GetExpire() int64 {
	if c == nil || c.Expire == nil {
		return 3600 // 默认1小时
	}
	return *c.Expire
}

func (c *Cache) GetNeedCache() []string {
	if c == nil || c.NeedCache == nil {
		return []string{}
	}
	return c.NeedCache
}

type Upload struct {
	Prefix *string `mapstructure:"prefix"`
}

type Qiniu struct {
	Bucket    *string `mapstructure:"bucket"`
	AccessKey *string `mapstructure:"accessKey"`
	SecretKey *string `mapstructure:"secretKey"`
	Region    *string `mapstructure:"region"`
}

// Aliyun 阿里云配置
type Aliyun struct {
	AccessKeyID     *string `mapstructure:"accessKeyId"`
	AccessKeySecret *string `mapstructure:"accessKeySecret"`
	Endpoint        *string `mapstructure:"endpoint"`
	Bucket          *string `mapstructure:"bucket"`
}

type Auth struct {
	IsAuth *bool `mapstructure:"isAuth"`
	Ignores    []string `mapstructure:"ignores"`
	NeedLogins []string `mapstructure:"needLogins"`
}

type Wx struct {
	AppId  *string `mapstructure:"appId"`
	Secret *string `mapstructure:"secret"`
	Token  *string `mapstructure:"token"`
	AesKey *string `mapstructure:"aesKey"`
}

type Postgres struct {
	Host         *string        `mapstructure:"host"`
	Port         *int           `mapstructure:"port"`
	User         *string        `mapstructure:"user"`
	Password     *string        `mapstructure:"password"`
	Database     *string        `mapstructure:"database"`
	SSLMode      *string        `mapstructure:"sslmode"`
	MaxIdleConns *int           `mapstructure:"maxIdleConns"`
	PingTimeout  *time.Duration `mapstructure:"pingTimeout"`
	MaxOpenConns *int           `mapstructure:"maxOpenConns"`
}

func (s *Server) GetHost() string {
	if s == nil || s.Host == nil {
		return "127.0.0.1"
	}
	return *s.Host
}

func (s *Server) GetPort() int {
	if s == nil || s.Port == nil {
		return 8080
	}
	return *s.Port
}

func (s *Server) GetMode() string {
	if s == nil || s.Mode == nil {
		return "release"
	}
	return *s.Mode
}

func (s *Server) GetReadTimeout() time.Duration {
	if s == nil || s.ReadTimeout == nil {
		return 5 * time.Second
	}
	return *s.ReadTimeout
}

func (s *Server) GetWriteTimeout() time.Duration {
	if s == nil || s.WriteTimeout == nil {
		return 5 * time.Second
	}
	return *s.WriteTimeout
}

func (s *Server) GetCros() []string {
	if s == nil || s.Cros == nil {
		return []string{}
	}
	return s.Cros
}

func (w *WxPay) GetAppId() string {
	if w == nil || w.AppId == nil {
		return ""
	}
	return *w.AppId
}

func (w *WxPay) GetMchId() string {
	if w == nil || w.MchId == nil {
		return ""
	}
	return *w.MchId
}

func (w *WxPay) GetNotifyUrl() string {
	if w == nil || w.NotifyUrl == nil {
		return ""
	}
	return *w.NotifyUrl
}

func (c *Config) GetDB() *DB {
	if c == nil {
		return nil
	}
	return c.DB
}

func (d *DB) GetPostgres() *Postgres {
	if d == nil {
		return nil
	}
	return d.Postgres
}

func (p *Postgres) GetHost() string {
	if p == nil || p.Host == nil {
		return "127.0.0.1"
	}
	return *p.Host
}

func (p *Postgres) GetPort() int {
	if p == nil || p.Port == nil {
		return 5432
	}
	return *p.Port
}

func (p *Postgres) GetDatabase() string {
	if p == nil || p.Database == nil {
		return ""
	}
	return *p.Database
}

func (p *Postgres) GetUser() string {
	if p == nil || p.User == nil {
		return ""
	}
	return *p.User
}

func (p *Postgres) GetPassword() string {
	if p == nil || p.Password == nil {
		return ""
	}
	return *p.Password
}

func (p *Postgres) GetSSLMode() string {
	if p == nil || p.SSLMode == nil {
		return "disable"
	}
	return *p.SSLMode
}

func (p *Postgres) GetMaxIdleConns() int {
	if p == nil || p.MaxIdleConns == nil {
		return 10
	}
	return *p.MaxIdleConns
}

func (p *Postgres) GetPingTimeout() time.Duration {
	if p == nil || p.PingTimeout == nil {
		return 5 * time.Second
	}
	return *p.PingTimeout
}

func (p *Postgres) GetMaxOpenConns() int {
	if p == nil || p.MaxOpenConns == nil {
		return 100
	}
	return *p.MaxOpenConns
}

func (d *DB) GetRedis() *Redis {
	if d == nil {
		return nil
	}
	return d.Redis
}

func (r *Redis) GetAddr() string {
	if r == nil || r.Addr == nil {
		return "127.0.0.1:6379"
	}
	return *r.Addr
}

func (r *Redis) GetDB() int {
	if r == nil || r.DB == nil {
		return 0
	}
	return *r.DB
}

func (r *Redis) GetPassword() string {
	if r == nil || r.Password == nil {
		return ""
	}
	return *r.Password
}

func (r *Redis) GetPoolSize() int {
	if r == nil || r.PoolSize == nil {
		return 100
	}
	return *r.PoolSize
}

func (r *Redis) GetMaxIdleConns() int {
	if r == nil || r.MaxIdleConns == nil {
		return 10
	}
	return *r.MaxIdleConns
}

func (r *Redis) GetMaxOpenConns() int {
	if r == nil || r.MaxOpenConns == nil {
		return 100
	}
	return *r.MaxOpenConns
}

func (c *Config) GetJwt() *Jwt {
	if c == nil {
		return nil
	}
	return c.Jwt
}

func (a *Auth) GetIsAuth() bool {
	if a == nil || a.IsAuth == nil {
		return false
	}
	return *a.IsAuth
}

func (a *Auth) GetIgnores() []string {
	if a == nil || a.Ignores == nil {
		return []string{}
	}
	return a.Ignores
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/etc")
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

func (m *Mysql) GetHost() string {
	if m == nil || m.Host == nil {
		return "127.0.0.1"
	}
	return *m.Host
}

func (m *Mysql) GetPort() int {
	if m == nil || m.Port == nil {
		return 3306
	}
	return *m.Port
}

func (m *Mysql) GetUser() string {
	if m == nil || m.User == nil {
		return ""
	}
	return *m.User
}

func (m *Mysql) GetPassword() string {
	if m == nil || m.Password == nil {
		return ""
	}
	return *m.Password
}

func (m *Mysql) GetDatabase() string {
	if m == nil || m.Database == nil {
		return ""
	}
	return *m.Database
}

func (m *Mysql) GetMaxIdleConns() int {
	if m == nil || m.MaxIdleConns == nil {
		return 10
	}
	return *m.MaxIdleConns
}

func (m *Mysql) GetPingTimeout() time.Duration {
	if m == nil || m.PingTimeout == nil {
		return 5 * time.Second
	}
	return *m.PingTimeout
}

func (m *Mysql) GetMaxOpenConns() int {
	if m == nil || m.MaxOpenConns == nil {
		return 100
	}
	return *m.MaxOpenConns
}
