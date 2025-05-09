package cluster

import (
	"encoding/json"
	"os"
	"slices"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/logger"
)

const (
	clusterctag            = "[CONFIG]"
	defaultPort            = 8080
	defaultChSize          = 512
	defaultProtocolVersion = 1
)

type ClusterConfig struct {
	Port            int      `json:"port"`
	Servers         []Server `json:"servers"`
	ChannelSize     int      `json:"chsize"`
	ProtocolVersion int      `json:"protocolVersion"`
	Web             string   `json:"webPath"`
}

func defaultForErr(err error) ClusterConfig {
	assert.Assert(err != nil)

	logger.Printf("%v Error reading config: %v\n", clusterctag, err)
	cfg := DefaultConfig()
	logger.Printf("%v Using default config: %v\n", clusterctag, cfg)
	return cfg

}

func DefaultConfig() ClusterConfig {
	return ClusterConfig{
		Port:            defaultPort,
		ProtocolVersion: defaultProtocolVersion,
		ChannelSize:     defaultChSize,
		Servers:         []Server{},
	}
}

func ReadConfig(filepath string) ClusterConfig {
	content, err := os.ReadFile(filepath)

	if err != nil {
		return defaultForErr(err)
	}

	var cfg ClusterConfig
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return defaultForErr(err)
	}

	if cfg.Port == 0 {
		cfg.Port = defaultPort
	}

	if cfg.ChannelSize == 0 {
		cfg.ChannelSize = defaultChSize
	}

	if cfg.ProtocolVersion == 0 {
		cfg.ProtocolVersion = defaultProtocolVersion
	}

	if slices.ContainsFunc(cfg.Servers, func(s Server) bool {
		return s.Identifier == ""
	}) {
		logger.Fatalf("%v Error: Missing server tokens\n", clusterctag)
	}

	logger.Printf("%v Using config: %v\n", clusterctag, cfg)
	return cfg
}
