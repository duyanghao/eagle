package main

import (
	"flag"
	"fmt"
	"github.com/duyanghao/eagle/seeder/bt"
	"github.com/duyanghao/eagle/seeder/muxconf"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

var (
	argRootDataDir string
	argOrigin      string
	argTrackers    string
	argPort        int
	argVerbose     bool
	argLimitSize   string
	seeder         *bt.Seeder
)

func main() {
	log.Infof("launch seeder on port: %d", argPort)

	// start seeder
	err := http.ListenAndServe(fmt.Sprintf(":%d", argPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.IntVar(&argPort, "port", 65005, "The port seeder listens to")
	flag.StringVar(&argRootDataDir, "rootdir", "/data/", "The root directory of seeder")
	flag.StringVar(&argOrigin, "origin", "", "The data origin of seeder")
	flag.StringVar(&argTrackers, "trackers", "", "The tracker list of seeder")
	flag.BoolVar(&argVerbose, "verbose", false, "verbose")
	flag.StringVar(&argLimitSize, "limitsize", "100G", "disk cache limit,format:xxxT/G")
	flag.Parse()
	if argVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	trackers := strings.Split(argTrackers, ",")
	c := &bt.Config{
		EnableUpload:  true,
		EnableSeeding: true,
		IncomingPort:  50017,
	}
	// transform ratelimiter
	switch argLimitSize[len(argLimitSize)-1:] {
	case "G":
		c.CacheLimitSize, _ = strconv.ParseInt(argLimitSize[:len(argLimitSize)-1], 10, 64)
		c.CacheLimitSize *= 1024 * 1024 * 1024
	case "T":
		c.CacheLimitSize, _ = strconv.ParseInt(argLimitSize[:len(argLimitSize)-1], 10, 64)
		c.CacheLimitSize *= 1024 * 1024 * 1024 * 1024
	}
	log.Debugf("cache limit size: %d", c.CacheLimitSize)
	seeder = bt.NewSeeder(argRootDataDir, argOrigin, trackers, c)
	err := seeder.Run()
	if err != nil {
		log.Fatal(err)
	}
	muxconf.InitMux(seeder)
}
