package options

import (
	"flag"
)

type Options struct {
	Version bool

	ListenAddr   string
	Debug        bool
	ColorizeLogs bool

	ConfigPath string
}

func Init() *Options {
	opts := new(Options)
	flag.BoolVar(&opts.Version, "v", false, "prints the version")
	flag.BoolVar(&opts.Version, "version", false, "prints the version")
	flag.BoolVar(&opts.Debug, "debug", false, "enable debug logging")
	flag.BoolVar(&opts.ColorizeLogs, "colorize-logs", false, "colorize log messages")
	flag.StringVar(&opts.ListenAddr, "address", ":8080", "server/bind address in format [host]:port")
	flag.StringVar(&opts.ConfigPath, "config", "config.yaml", "path to YAML config file")
	flag.Parse()
	return opts
}
