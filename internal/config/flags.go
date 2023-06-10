package config

import (
	"github.com/90amper/metmon/internal/models"
	pflag "github.com/spf13/pflag"
)

var CmdFlags models.CmdFlags

func init() {
	pflag.StringVarP(&CmdFlags.ServerUrl, "address", "a", "localhost:8080", "server URL")
	pflag.IntVarP(&CmdFlags.ReportInterval, "report", "r", 10, "metrics report interval")
	pflag.IntVarP(&CmdFlags.PollInterval, "poll", "p", 2, "metrics poll interval")

	pflag.Parse()
}
