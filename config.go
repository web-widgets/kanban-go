package main

type ConfigServer struct {
	Port string
	Cors bool
}

type ConfigDB struct {
	Path         string
	ResetOnStart bool
}

type AppConfig struct {
	Server     ConfigServer
	DB         ConfigDB
	BinaryData string
}
