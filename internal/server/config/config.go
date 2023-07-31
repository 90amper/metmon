package config

import (
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/90amper/metmon/internal/models"
	"github.com/caarlos0/env/v6"
	pflag "github.com/spf13/pflag"
)

var Config models.Config

func init() {
	user, _ := user.Current()
	wdPath, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		Config.PathSeparator = "\\"
		pflag.StringVarP(&Config.FileStoragePath, "storefile", "f", strings.Join([]string{user.HomeDir, "metrics-db.json"}, Config.PathSeparator), "metrics store file path")
	} else {
		Config.PathSeparator = "/"
		pflag.StringVarP(&Config.FileStoragePath, "storefile", "f", "/tmp/metrics-db.json", "metrics store file path")
	}
	Config.ProjPath = wdPath
	pflag.StringVarP(&Config.ServerURL, "address", "a", "localhost:8080", "server URL")
	pflag.IntVarP(&Config.StoreInterval, "storeint", "i", 300, "metrics store interval")
	pflag.BoolVarP(&Config.Restore, "restore", "r", true, "restore metrics after startup")
	pflag.StringVarP(&Config.DatabaseDsn, "database-dsn", "d", "", "database dsn")
	pflag.BoolVarP(&Config.Cleanup, "cleanup", "c", false, "recreate schema at startup")
	// postgresql://postgres:postgres@localhost:5454/store
	pflag.StringVarP(&Config.HashKey, "hash-key", "k", "$ecretKey", "hash secret key")
	// pflag.StringVarP(&Config.HashAlg, "hash-alg", "h", "SHA256", "hash algorithm")

	pflag.Parse()

	env.Parse(&Config)

	Config.HashAlg = "SHA256"
}
