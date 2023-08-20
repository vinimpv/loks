// main.go

package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/deployment"
	"vinimpv/loks/pkg/docker"
	"vinimpv/loks/pkg/renderer"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "mycli"}

	var renderCmd = &cobra.Command{
		Use:   "render [configPath]",
		Short: "Render Kubernetes resources",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configPath := args[0]
			result, err := renderer.Render(configPath)
			if err != nil {
				log.Fatalf("Error rendering: %v", err)
			}
			fmt.Println(result)
		},
	}

	var deployCmd = &cobra.Command{
		Use:   "bootstrap [configPath]",
		Short: "Deploy rendered Kubernetes resources using kapp",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configPath := args[0]
			cfg, err := config.LoadConfig(configPath)
			if err != nil {
				log.Fatalf("error loading config: %v", err)
			}

			err = cluster.CreateCluster(cfg.Name)
			if err != nil {
				log.Fatalf("error creating cluster: %v", err)
			}
			wg := sync.WaitGroup{}
			for _, component := range cfg.Components {
				wg.Add(1)
				go func(component config.Component) {
					err = docker.PullImage(component.Image)
					if err != nil {
						log.Fatalf("error pulling image: %s, %v", component.Image, err)
					}
					err = cluster.LoadImage(cfg.Name, component.Image)
					if err != nil {
						log.Fatalf("error loading image to cluster: %s, %v", component.Image, err)
					}
					wg.Done()
				}(component)
			}
			wg.Wait()

			renderedYaml, err := renderer.Render(configPath)
			if err != nil {
				log.Fatalf("error rendering k8s manifests: %v", err)
			}

			err = deployment.DeployToCluster(cfg.Name, renderedYaml)
			if err != nil {
				log.Fatalf("error deploying to cluster: %v", err)
			}

		},
	}

	rootCmd.AddCommand(renderCmd, deployCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
