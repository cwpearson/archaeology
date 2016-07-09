package archaeology

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

type ArchCfg struct {
	ArchMode     string `ini:"mode"`
	IndexDbPath  string `ini:"path"`
	Include_file string `ini:"include_file"`
	Exclude_file string `ini:"exclude_file"`
	Includes     []string
	Excludes     []string
}

func defaultCfg() ArchCfg {

	return ArchCfg{
		ArchMode:    "local",
		IndexDbPath: "/home/pearson/.archaeology/index.db"}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	lines := make([]string, 0)
	line, err := r.ReadString('\n')
	for err == nil {
		line = strings.TrimSuffix(line, "\n")
		lines = append(lines, line)
		line, err = r.ReadString('\n')
	}
	if err != io.EOF {
		return nil, err
	} else {
		return lines, nil
	}
}

func GetConfig(path string) (ArchCfg, error) {
	// Load the ini file
	cfg := new(ArchCfg)
	err := ini.MapTo(cfg, path)
	if err != nil {
		log.Fatal(err)
	}

	// load the include and exclude files, if they exist
	if fileExists(cfg.Include_file) {
		includes, err := readLines(cfg.Include_file)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Includes = includes
	} else {
		log.Print("includes file '", cfg.Include_file, "' not found")
	}
	if fileExists(cfg.Exclude_file) {
		excludes, err := readLines(cfg.Exclude_file)
		if err != nil {
			log.Fatal(err)
		}
		cfg.Excludes = excludes
	} else {
		log.Print("excludes file '", cfg.Exclude_file, "' not found")
	}
	return *cfg, nil
}
