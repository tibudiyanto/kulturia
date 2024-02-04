package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	R2AccountId       string
	R2AccessKeyId     string
	R2AccessKeySecret string
	R2BucketName      string
	R2PublicURL       string
}

func GetConfig() (Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	return Config{
		R2AccountId:       viper.GetString("R2_ACCOUNT_ID"),
		R2AccessKeyId:     viper.GetString("R2_ACCESS_KEY_ID"),
		R2AccessKeySecret: viper.GetString("R2_ACCESS_KEY_SECRET"),
		R2BucketName:      viper.GetString("R2_BUCKET_NAME"),
		R2PublicURL:       viper.GetString("R2_PUBLIC_URL"),
	}, nil

}
