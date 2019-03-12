package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dolegi/uci"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	Token   string // done
	Botname string // done
	Url     string // done
	Engine  struct {
		Path    string // done
		Options struct {
			Threads int // done
			Hash    int // done
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
var eng *uci.Engine

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

	retries := 0
	for _ = range time.Tick(50 * time.Millisecond) {
		if eng.IsReady() {
			break
		}
		if retries > 10 {
			break
		}
		retries++
	}
	if retries > 10 {
		log.Fatal("Failed to find ready engine. Max retries exceeded", retries)
	}
	eng.SetOption("Threads", conf.Engine.Options.Threads)
	eng.SetOption("Hash", conf.Engine.Options.Hash)

	streamEvent()
}
