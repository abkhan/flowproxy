package main

import (
	"flag"
	"log"
)

// read config file

// 2 sections
//  1. Targets (where to send the traffic to)
//		a. Default for all new sources
//  2. Learned sources - and the destination for them

// Flow;
// read the config file, must have one destination and one default
// any input source, send to default and add to config under learned sources

// change of default - verify new defaults is in destinations, otherwise add
//  - keep all the previous device to destination as is
// - if there is a new device, that does not have a destination, add to the dev-to-dest
// list

func main() {
	confinit := flag.Bool("init", false, "create new config, must provide default destination")
	defdest := flag.String("d", "", "IP or FQDN of default destination")
	destlist := flag.String("dlist", "", "comma separated IP or FQDN of destination(s)")
	ipport := flag.Int("p", 9995, "UDP listen port")
	confile := flag.String("f", "./config.yml", "config file name")
	flag.Parse()

	if *confile == "" {
		log.Fatalf("must provide config file name")
	}
	if *confinit && (*defdest == "" || *confile == "") {
		log.Fatalf("when init flag is there, need default destination and config file name to write to")
	}

	var pd *ProxyData

	// if init flag is there, build a conf struct and wrtie to file
	if *confinit {
		pd = NewProxyData(*defdest, *destlist)
		if err := pd.WriteToFile(*confile); err != nil {
			log.Fatalf("saving to conf file [%s] error: %+v", *confile, err)
		}
	} else {
		// else, read from config file, if def is given, add to destlist if not there already
		// and set as default
		pd = NewProxyData(*defdest, *destlist)
		if err := pd.ReadFile(*confile); err != nil {
			log.Fatalf("reading conf file [%s] error: %+v", *confile, err)
		}
		if *defdest != "" {
			pd.AddDef(*defdest) // over-written by readFile, has to written again
		}
	}

	StartProxy(pd, *ipport)

	// for now stop here,
	// todo, add graceful shutdown
	select {}
}
