package main

import "web-widgets/kanban-go/data"

type ConfigServer struct {
	URL  string
	Port string
	Cors bool
}

type AppConfig struct {
	Server     ConfigServer
	DB         data.DBConfig
	BinaryData string
}
