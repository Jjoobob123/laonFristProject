package config

import (
	"log"

	"github.com/go-ini/ini"
)

// ConfigData 구조체는 설정 데이터를 보유합니다.
type ConfigData struct {
	DB struct {
		Host     string
		User     string
		Password string
		Database string
		Port     int
	}
	Network struct {
		Method int
		Host   string
		Port   int
	}
}

// LoadConfig 함수는 ini 파일로부터 설정 값을 불러옵니다.
func LoadConfig(filename string) (*ConfigData, error) {
	cfg := new(ConfigData)

	file, err := ini.Load("config.ini")
	if err != nil {
		return nil, err
	}

	if err := file.Section("DB").MapTo(&cfg.DB); err != nil {
		log.Fatal("DB 설정 값을 불러오는 중 에러 발생:", err)
	}

	if err := file.Section("Network").MapTo(&cfg.Network); err != nil {
		log.Fatal("Network 설정 값을 불러오는 중 에러 발생:", err)
	}

	return cfg, nil
}
