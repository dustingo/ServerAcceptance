package main

import (
	"flag"

	"github.com/dustingo/ServerAcceptance/work"
)

var dryrun = flag.Bool("dryrun", false, "just Parse config file")
var precheck = flag.Bool("precheck", false, "precheck by using config.toml which you have edited it before")
var printjson = flag.Bool("printjson", false, "print all information as json")

func main() {
	flag.Parse()
	if *precheck {
		work.PreCheck()
		return
	}
	if *printjson {
		all := new(work.AllInfo)
		all.PrintJSON()
		return
	}
	work.LastCheck()
}
