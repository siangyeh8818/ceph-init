package main

import (
	"os"

	cephapi "github.com/pnetwork/sre.ceph.init/internal/api"
	conf "github.com/pnetwork/sre.ceph.init/internal/config"
)

func main() {
	configPath := os.Getenv("PATH_CONFIG")
	if configPath == "" {
		configPath = "config.json"
	}
	InitBucketList := os.Getenv("INIT_BUCKET_LIST")

	RemoveBucketList := os.Getenv("REMOVE_BUCKET_LIST")

	caphInitConfig := conf.Config{}
	caphInitConfig.InitBucketList(InitBucketList)
	caphInitConfig.RemoveBucketList(RemoveBucketList)
	caphInitConfig.BaseConfig.InitConfig(configPath)
	cephapi.CephAPI(&caphInitConfig)
	//caphInitConfig.InitConfig(configPath)

}
