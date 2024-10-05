package config

import "github.com/pelletier/go-toml"

var Mysql struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Dbname   string `toml:"dbname"`
}

var configFile = "config.toml"

func init() {
	conf, err := toml.LoadFile(configFile)
	if err != nil {
		panic("load config file failed: " + err.Error())
	}

	if err := conf.Get("mysql").(*toml.Tree).Unmarshal(&Mysql); err != nil {
		panic("unmarshal config file failed: " + err.Error())
	}
}
