package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var config configuration

type buddy struct {
	Address   string `yaml:"address,omitempty"`
	Mac       string `yaml:"mac,omitempty"`
	Username  string `yaml:"username"`
	Pass      string `yaml:"pass"`
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Reachable bool
}

type einsy struct {
	Address   string `yaml:"address,omitempty"`
	Mac       string `yaml:"mac,omitempty"`
	Apikey    string `yaml:"apikey"`
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`
	Reachable bool
}

type configuration struct {
	Exporter struct {
		MetricsPort    int    `yaml:"metrics_port"`
		ScrapeTimeout  int    `yaml:"scrape_timeout"`
		ReloadInterval int    `yaml:"reload_interval"`
		LogLevel       string `yaml:"log_level"`
	} `yaml:"exporter"`
	Printers struct {
		Buddy []buddy `yaml:"buddy"`
		Einsy []einsy `yaml:"einsy"`
	} `yaml:"printers"`
}

func setLogLevel(level string) string {
	if level == "" {
		level = "info"
	}

	level = strings.ToLower(level)
	var zeroLogLevel zerolog.Level

	switch level {
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "trace":
		zeroLogLevel = zerolog.TraceLevel
	case "error":
		zeroLogLevel = zerolog.ErrorLevel
	case "panic":
		zeroLogLevel = zerolog.PanicLevel
	case "fatal":
		zeroLogLevel = zerolog.FatalLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zeroLogLevel)
	print(level)
	return level
}

func loadConfigFile() {
	config = probeConfigFile(parseConfig(getConfigPath()))
}

func getConfigPath() string {
	cfgFile := os.Getenv("BUDDY_EXPORTER_CONFIG")
	if cfgFile == "" {
		pwd, e := os.Getwd()
		if e != nil {
			fmt.Println(e)
			os.Exit(1)
		}
		cfgFile = pwd + "/buddy.yaml"
	}

	return cfgFile
}

func parseConfig(path string) configuration {
	f, e := os.ReadFile(path)
	if e != nil {
		log.Panic().Msg(e.Error())
	}
	var p configuration
	if e := yaml.Unmarshal(f, &p); e != nil {
		log.Error().Msg(e.Error())
	}
	return p
}

func configReloader() {
	ticker := time.NewTicker(time.Duration(config.Exporter.ReloadInterval) * time.Second)

	for t := range ticker.C {
		log.Info().Msg(fmt.Sprintf("Config reloaded at: %v\n", t.UTC()))
		loadConfigFile()
	}
}
