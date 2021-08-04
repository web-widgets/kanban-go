package main

type ConfigServer struct {
	URL  string
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
