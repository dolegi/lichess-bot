package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dolegi/uci"
	"io/ioutil"
	"log"
	"net/http"
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
			Contempt     int // done
			Threads      int // done
			Hash         int // done
			MoveOverhead int // done
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
var whiteFirst bool = true

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
	log.Printf("%v\n", getUsersStatus(conf.Botname))

	if len(os.Args) >= 3 && os.Args[2] == "upgrade" {
		resp := request("POST", "bot/account/upgrade")

		if resp.StatusCode == http.StatusOK {
			log.Println("Account upgraded to bot account")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Println("Failed to upgrade account")
			log.Println(string(body))
		}
		resp.Body.Close()
		return
	}

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
	eng.SetOption("Contempt", conf.Engine.Options.Contempt)
	if conf.Engine.Options.Threads > 0 {
		eng.SetOption("Threads", conf.Engine.Options.Threads)
	}
	if conf.Engine.Options.Hash > 0 {
		eng.SetOption("Hash", conf.Engine.Options.Hash)
	}
	eng.SetOption("Move Overhead", conf.Engine.Options.MoveOverhead)

	streamEvent(eng)
}
