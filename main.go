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
			Nodes    int // done
			Depth    int // done
			Movetime int // done
		}
	}
	Network struct {
		Latency int // done
	}
	Challenge struct {
		Variants []string // done
		Speeds   []string // done
		Modes    []string // done
	}
	Game struct {
		Aborttime int
	}
}

var conf Config
var white bool = true

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing path to config file. Usage `lichess-bot path/to/config`")
	}
	configFile, err := os.Open(os.Args[1])
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
		log.Fatal("Failed to create a new engine", err)
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
	if conf.Engine.Options.Threads > 0 {
		eng.SetOption("Threads", conf.Engine.Options.Threads)
	}
	if conf.Engine.Options.Hash > 0 {
		eng.SetOption("Hash", conf.Engine.Options.Hash)
	}

	streamEvent(eng)
}
