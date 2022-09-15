package main

import (
	"github.com/PereRohit/util/server"
	"github.com/vatsal278/htmltopdfsvc/internal/config"
	"github.com/vatsal278/htmltopdfsvc/internal/router"
)

func main() {
	appContainer := config.GetAppContainer()
	r := router.Register(appContainer)

	//log.Fatal(http.ListenAndServe(":9080", r))
	server.Run(r)
}
