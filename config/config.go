//	This file is part of Fwew.
//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package config handles... the configuration file stuff. Probably.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tirea/fwew/util"
)

// Config is a struct designed to hold the values of the configuration file when loaded
type Config struct {
	Language   string `json:"language"`
	PosFilter  string `json:"posFilter"`
	UseAffixes bool   `json:"useAffixes"`
}

// ReadConfig reads a configuration file and puts the data into Config struct
func ReadConfig() Config {
	configfile, e := ioutil.ReadFile(util.Text("config"))
	if e != nil {
		fmt.Printf("%s: %v\n", util.Text("fileError"), e)
	}

	var config Config
	json.Unmarshal(configfile, &config)

	return config
}

func (c Config) String() string {
	// this string only doesn't get translated or called from util.Text() because they're var names
	return fmt.Sprintf("Language: %s\nPosFilter: %s\nUseAffixes: %t\n", c.Language, c.PosFilter, c.UseAffixes)
}
