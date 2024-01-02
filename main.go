package main

import (
	"fmt"
	"log"

	"github.com/MRdyRy/go-datagrid/config"
)

const (
	PROTOCOL    = "http"
	CACHE_NAME  = "test"
	CACHE_KEY   = "test"
	CACHE_VALUE = "test"
	url         = "127.0.0.1"
	port        = "11222"
	user        = "ryan"
	pass        = "changeme"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load env", err.Error())
	}
	fmt.Println(cfg)

	// dg := config.NewDatagridClient(PROTOCOL, url, port, user, pass)
	dg := config.NewDatagridClient(cfg.DatagridProtocol, cfg.DatagridUrl, cfg.DatagridPort, cfg.DatagridUser, cfg.DatagridPass)

	err = dg.AddToCache(CACHE_NAME, CACHE_KEY, CACHE_VALUE)
	if err != nil {
		log.Println(err)
	}

	res, err := dg.GetDataFromCache(CACHE_NAME, CACHE_KEY)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("data", res)

	res, err = dg.GetAllKeysFromCache("test")
	if err != nil {
		fmt.Println(err)
	}
	log.Println(res)
}
