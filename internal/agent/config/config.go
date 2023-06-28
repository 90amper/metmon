package config

import (
	"github.com/90amper/metmon/internal/models"
	"github.com/caarlos0/env/v6"
	pflag "github.com/spf13/pflag"
)

var Config models.Config

func init() {
	pflag.StringVarP(&Config.ServerURL, "address", "a", "localhost:8080", "server URL")
	pflag.IntVarP(&Config.ReportInterval, "report", "r", 10, "metrics report interval")
	pflag.IntVarP(&Config.PollInterval, "poll", "p", 2, "metrics poll interval")

	env.Parse(&Config)
}
