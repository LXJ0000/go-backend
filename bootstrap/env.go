package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	ServerAddr     string `mapstructure:"SERVER_ADDR"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`

	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`

	RateLimit       int `mapstructure:"RATE_LIMIT"`
	RateLimitWindow int `mapstructure:"RATE_LIMIT_WINDOW"`

	SnowflakeStartTime string `mapstructure:"SNOWFLAKE_START_TIME"`
	SnowflakeMachineID int64  `mapstructure:"SNOWFLAKE_MACHINE_ID"`

	LocalStaticPath string `mapstructure:"LOCAL_STATIC_PATH"`
	UrlStaticPath   string `mapstructure:"URL_STATIC_PATH"`

	OpenIMServerDoamin string `mapstructure:"OPENIM_SERVER_DOMAIN"`

	SMSAppID      string `mapstructure:"SMS_APP_ID"`
	SMSSignName   string `mapstructure:"SMS_SIGN_NAME"`
	SMSTemplateID string `mapstructure:"SMS_TEMPLATE_ID"`
	SMSEndpoint   string `mapstructure:"SMS_ENDPOINT"`
	SMSRegionId   string `mapstructure:"SMS_REGION_ID"`

	MySQLHost     string `mapstructure:"MYSQL_HOST"`
	MySQLPORT     string `mapstructure:"MYSQL_PORT"`
	MySQLUsername string `mapstructure:"MYSQL_USERNAME"`
	MySQLPassword string `mapstructure:"MYSQL_PASSWORD"`
	MySQLDB       string `mapstructure:"MYSQL_DB"`

	RedisHost       string `mapstructure:"REDIS_HOST"`
	RedisPort       string `mapstructure:"REDIS_PORT"`
	RedisExpiration int    `mapstructure:"REDIS_EXPIRATION"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	RedisDB         int    `mapstructure:"REDIS_DB"`

	KafkaAddr      string `mapstructure:"KAFKA_ADDR"`
	PrometheusAddr string `mapstructure:"PROMETHEUS_ADDR"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
