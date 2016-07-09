package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	ar "github.com/cwpearson/archaeology/arch"
)

const (
	ARCH_CFG_ENV = "ARCHAEOLOGY_CFG"
)

func main() {

	log.SetFlags(log.Lshortfile)

	// Get the configuration
	cfg_path := os.Getenv(ARCH_CFG_ENV)
	if cfg_path == "" {
		usr, _ := user.Current()
		cfg_path = filepath.Join(usr.HomeDir, ".archaeology", "config")
		log.Printf("Environment variable %s was not set. Using %s\n", ARCH_CFG_ENV, cfg_path)
	}
	cfg, err := ar.GetConfig(cfg_path)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Printf("Using config: %#v\n", cfg)

	// Open a connection to the index database
	index := new(ar.IndexDB)
	err = index.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	index.CreateTables()
	defer index.Close()

}
