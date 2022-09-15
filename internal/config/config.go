package config

import (
	"github.com/vatsal278/go-redis-cache"
	"os"
)

type AppContainer struct {
	Cacher redis.Cacher
}

func GetAppContainer() *AppContainer {
	cacher := redis.NewCacher(redis.Config{Addr: os.Getenv("Address") + ":6379"})
	return &AppContainer{
		Cacher: cacher,
	}
}
