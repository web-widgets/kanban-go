package main

import "web-widgets/kanban-go/data"

type ConfigServer struct {
	URL  string
	Port string
	Cors []string
}

type AppConfig struct {
	Server     ConfigServer
	DB         data.DBConfig
	Features   data.Features
	BinaryData string
}
