package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cli = &cobra.Command{
	Use:     "web",
	Short:   "web renders markdown files",
	Version: "0.0.1",
}

func init() {
	cli.PersistentFlags().StringP("site", "s", "", "top directory for site")
	cli.PersistentFlags().StringP("config", "c", "config.yaml", "config file (optional)")
}

func main() {
	cli.ParseFlags(os.Args)
	site := cli.Flag("site").Value.String()
	if site == "" {
		log.Fatalf("--site is required and should point to the root directory of the website")
	}
	if err := os.Chdir(site); err != nil {
		log.Fatalf("Error changing directory to %q: %v", site, err)
	}

	cfg := cli.Flag("config").Value.String()
	viper.SetConfigName(cfg)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found")
		} else {
			// Config file was found but another error was produced
			log.Fatal(err)
		}
	}

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
