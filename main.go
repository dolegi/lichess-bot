package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token   string // done
	Botname string // done
	Url     string // done
	Engine  struct {
		Path    string
		Options struct {
			Threads int
			Hash    int
		}
		Go struct {
			Nodes    int
			Depth    int
			Movetime int
		}
	}
	Network struct {
		Latency int
	}
	Challenge struct {
		Variants []string
		Speeds   []string
		Modes    []string
	}
	Game struct {
		Aborttime int
	}
}

var conf Config
var eng *Engine

func main() {
	configFile, err := os.Open("config.toml")
	if err != nil {
		log.Fatal("Failed to open config.toml", err)
	}
	configData, _ := ioutil.ReadAll(configFile)
	configFile.Close()
	if _, err := toml.Decode(string(configData), &conf); err != nil {
		log.Fatal("Failed to decode config.toml", err)
	}
	log.Println(conf)

	eng, err := uci.NewEngine(conf.Engine.Path)
	if err != nil {
		log.Fatal("Failed to create new engine", err)
	}

	streamEvent()
}
