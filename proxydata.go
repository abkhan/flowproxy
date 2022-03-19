package main

import (
	"io/ioutil"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ProxyData struct {
	DevToDest    map[string]string // list of device [key] to destination [value]
	Destinations []string          // list of destination IP address/fqdn
	Default      string            // default destination at a given time
	mapm         sync.Mutex
}

func NewProxyData(def, dest string) *ProxyData {
	rv := ProxyData{
		DevToDest: make(map[string]string),
		Default:   def,
	}

	dests := strings.Split(dest, ",")
	if len(dests) > 0 {
		rv.Destinations = dests
	}
	return &rv
}

func (p *ProxyData) GetDest(dev string) string {
	p.mapm.Lock()
	dest := p.DevToDest[dev]
	p.mapm.Unlock()

	if dest == "" {
		dest = p.Default
	}
	return dest
}

func (p *ProxyData) WriteToFile(f string) error {
	ydata, err := yaml.Marshal(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f, ydata, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyData) ReadFile(f string) error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyData) AddDef(d string) {
	there := false
	p.Default = d
	for _, sd := range p.Destinations {
		if d == sd {
			there = true
		}
	}
	if !there {
		p.Destinations = append(p.Destinations, d)
	}
}
