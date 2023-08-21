package loks

import (
	"log"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Destroy the exiting cluster specified in the config file",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _, err := config.LoadUserConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		err = cluster.DestroyCluster(cfg.Name)
		if err != nil {
			log.Fatalf("error destroying cluster: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
