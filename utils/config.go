package utils

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type AppInfo struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
}

type Http struct {
	Port    string `toml:"port"`
	Timeout int    `toml:"timeout"`
}

type Etcd struct {
	BindAddrs []string `toml:"bindaddr"`
	Prefix    string   `toml:"prefix"`
}

type Alarm struct {
	Module  string
	Names   []string
	Url     string
	Timeout int
	IsOn    bool `toml:"Status"`
	Level   int
}

type LogConfig struct {
	Level string
}

type Config struct {
	AppConf   AppInfo   `toml:"app"`
	HTTPConf  Http      `toml:"http"`
	EtcdConf  Etcd      `toml:"etcd"`
	AlarmConf Alarm     `toml:"alarm"`
	LogConf   LogConfig `toml:"log"`
}

const confPath = "./conf/conf.toml"

var Conf Config

func init() {
	file, err := os.OpenFile(confPath, os.O_RDWR, 0666)
	if err != nil {
		panic(err.Error())
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	_, err = toml.Decode(string(content), &Conf)
	if err != nil {
		panic(err.Error())
	}
}
