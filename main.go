package main

import (
	"im_cacher/src/aof"
	"im_cacher/src/cache"
	"im_cacher/src/config"
	"im_cacher/src/server"
	"time"
)

func main() {

	cnf := config.GetConfig()

	AOF, err := aof.OpenAOF(".aof")
	if err != nil {
		panic(err.Error())
	}

	dict := cache.NewDict(cnf.InitialSize, AOF, cnf.MaxKeysAmount)

	AOF.Scan(dict.Iterator)

	close := dict.StartEvictionLoop(time.Second * time.Duration(cnf.EvictionLoopInterval))
	defer close()

	server.StartListen(cnf.Addr, dict, AOF)
}
