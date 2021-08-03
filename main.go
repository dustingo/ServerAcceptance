package main

import (
	"flag"
	"fmt"

	"github.com/dustingo/ServerAcceptance/work"
)

var (
	precheck  = flag.Bool("precheck", false, "precheck by using config file  which you have edited it before")
	lastcheck = flag.Bool("lastcheck", false, "lastcheck by using config file  which you have edited it before")
	printjson = flag.Bool("printjson", false, "print all information as json")
	config    = flag.String("config", "", "config of server info or software info")
)

func main() {
	flag.Parse()
	if *precheck {
		if *config == "" {
			fmt.Println("config missed!")
			return
		}
		work.PreCheck(*config)
		return
	}
	if *printjson {
		all := new(work.AllInfo)
		all.PrintJSON()
		return
	}
	if *lastcheck {
		if *config == "" {
			fmt.Println("config missed")
			return
		}
		work.LastCheck(*config)
		return
	}
}
