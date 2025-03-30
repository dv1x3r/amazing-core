package config

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	Servers struct {
		API  string `json:"api"`
		Game string `json:"game"`
	} `json:"servers"`
	Settings struct {
		AssetDeliveryURL string `json:"assetDeliveryURL"`
		SyncServerIP     string `json:"syncServerIP"`
		SyncServerPort   int    `json:"syncServerPort"`
	} `json:"settings"`
	Storage struct {
		Databases struct {
			Core string `json:"core"`
			Blob string `json:"blob"`
		} `json:"databases"`
	} `json:"storage"`
	Secure Secure `json:"secure,omitzero"`
}

type Secure struct {
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
	Session struct {
		Key    string `json:"key"`
		Secure bool   `json:"secure"`
	} `json:"session"`
	CSRF struct {
		Key            string   `json:"key"`
		Secure         bool     `json:"secure"`
		TrustedOrigins []string `json:"trustedOrigins"`
	} `json:"csrf"`
}

func (cfg Config) MarshalJSON() ([]byte, error) {
	type alias Config
	shadow := alias(cfg)
	shadow.Secure = Secure{}
	return json.Marshal(shadow)
}

func Get() Config {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "config.json"
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(data, &config)
		if err != nil {
			panic(err)
		}
	})

	return config
}
