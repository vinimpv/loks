package loks

import (
	"fmt"
	"log"
	"path/filepath"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/deployment"
	"vinimpv/loks/pkg/docker"
	"vinimpv/loks/pkg/renderer"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var updateCommand = &cobra.Command{
	Use:   "update [component]",
	Short: "Updates the specified component",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, currentContextRootPath, err := config.LoadUserConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		componentName := args[0]
		component, err := cfg.GetComponent(componentName)
		if err != nil {
			log.Fatalf("error getting component: %v", err)
		}

		devTag := fmt.Sprintf("%s:dev", component.Name)
		randomTag := fmt.Sprintf("%s:%s", component.Name, uuid.New().String())

		if component.BuildDev == "" {
			// pull production image
			err = docker.PullImage(component.Image)
			if err != nil {
				log.Fatalf("error pulling image: %v", err)
			}
			err = docker.Tag(component.Image, randomTag)
			if err != nil {
				log.Fatalf("error tagging image: %v", err)
			}
			err = cluster.LoadImage(cfg.Name, component.Image)
			if err != nil {
				log.Fatalf("error pushing image: %v", err)
			}

		} else {
			// build dev image
			err = docker.BuildDev(filepath.Join(currentContextRootPath, component.Name), component.BuildDev)
			if err != nil {
				log.Fatalf("error building dev image: %v", err)
			}
			err = docker.Tag(devTag, randomTag)
			if err != nil {
				log.Fatalf("error tagging dev image: %v", err)
			}

		}

		err = cluster.LoadImage(cfg.Name, randomTag)
		if err != nil {
			log.Fatalf("error pushing dev image: %v", err)
		}

		renderedYaml, err := renderer.Render(filepath.Join(currentContextRootPath, "loks.yaml"), fmt.Sprintf("update_context.image=%s", randomTag), fmt.Sprintf("update_context.component=%s", component.Name))
		if err != nil {
			log.Fatalf("error rendering k8s manifests: %v", err)
		}

		err = deployment.DeployToCluster(cfg.Name, renderedYaml)
		if err != nil {
			log.Fatalf("error deploying to cluster: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCommand)
}
