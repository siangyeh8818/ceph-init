package config

import (
	"io/ioutil"
	"os"

	"github.com/buger/jsonparser"

	//lib "github.com/pnetwork/sre.monitor.metrics/internal/app/pentium_api/lib"
	log "github.com/sirupsen/logrus"
)

// Config struct holds all of the runtime configuration
type Config struct {
	BaseConfig         BaseConfig
	INIT_BUCKET_LIST   string
	ROMOVE_BUCKET_LIST string
}

type BaseConfig struct {
	PN_GLOBAL_PORTAL             string
	PN_GLOBAL_JWT_PASSPHRASE     string
	PN_GLOBAL_STORAGE_SECRET_ID  string
	PN_GLOBAL_STORAGE_SECRET_KEY string
	PN_GLOBAL_STORAGE_ENDPOINT   string
	PN_GLOBAL_REDIS              string
	MY_POD_NAMESPACE             string
}

// LoadConfig load configuration from mounted K8s Secret
func LoadConfig(configPath string) ([]byte, error) {
	config, err := ioutil.ReadFile(configPath)
	return config, err
}

// InitConfig Set the BaseConfig with loaded contents
func (baseCfg *BaseConfig) InitConfig(configPath string) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	baseCfg.PN_GLOBAL_PORTAL, err = jsonparser.GetString(config, "PN_GLOBAL_PORTAL")
	if err != nil {
		log.Fatalf("PN_GLOBAL_PORTAL %v:", err)
	}
	baseCfg.PN_GLOBAL_JWT_PASSPHRASE, err = jsonparser.GetString(config, "PN_GLOBAL_JWT_PASSPHRASE")
	if err != nil {
		log.Fatalf("PN_GLOBAL_JWT_PASSPHRASE %v:", err)
	}
	baseCfg.PN_GLOBAL_STORAGE_SECRET_ID, err = jsonparser.GetString(config, "PN_GLOBAL_STORAGE_SECRET_ID")
	if err != nil {
		log.Fatalf("PN_GLOBAL_STORAGE_SECRET_ID %v:", err)
	}
	baseCfg.PN_GLOBAL_STORAGE_SECRET_KEY, err = jsonparser.GetString(config, "PN_GLOBAL_STORAGE_SECRET_KEY")
	if err != nil {
		log.Fatalf("PN_GLOBAL_STORAGE_SECRET_KEY %v:", err)
	}
	baseCfg.PN_GLOBAL_STORAGE_ENDPOINT, err = jsonparser.GetString(config, "PN_GLOBAL_STORAGE_ENDPOINT")
	if err != nil {
		log.Fatalf("PN_GLOBAL_STORAGE_ENDPOINT %v:", err)
	}

	baseCfg.PN_GLOBAL_REDIS, err = jsonparser.GetString(config, "PN_GLOBAL_REDIS")
	if err != nil {
		log.Fatalf("PN_GLOBAL_REDIS %v:", err)
	}

	baseCfg.MY_POD_NAMESPACE = os.Getenv("MY_POD_NAMESPACE")

	log.Printf("SECRET_CONTENT.PN_GLOBAL_PORTAL %v:", baseCfg.PN_GLOBAL_PORTAL)
	log.Printf("SECRET_CONTENT.PN_GLOBAL_JWT_PASSPHRASE %v:", baseCfg.PN_GLOBAL_JWT_PASSPHRASE)
	log.Printf("SECRET_CONTENT.PN_GLOBAL_STORAGE_SECRET_ID %v:", baseCfg.PN_GLOBAL_STORAGE_SECRET_ID)
	log.Printf("SECRET_CONTENT.PN_GLOBAL_STORAGE_SECRET_KEY %v:", baseCfg.PN_GLOBAL_STORAGE_SECRET_KEY)
	log.Printf("SECRET_CONTENT.PN_GLOBAL_STORAGE_ENDPOINT %v:", baseCfg.PN_GLOBAL_STORAGE_ENDPOINT)
	log.Printf("SECRET_CONTENT.PN_GLOBAL_REDIS %v:", baseCfg.PN_GLOBAL_REDIS)
	log.Printf("ENV.MY_POD_NAMESPACE %v:", baseCfg.MY_POD_NAMESPACE)
}

func (config *Config) InitBucketList(list string) {
	config.INIT_BUCKET_LIST = list
}

func (config *Config) RemoveBucketList(list string) {
	config.ROMOVE_BUCKET_LIST = list
}
