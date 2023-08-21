package loks

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"vinimpv/loks/pkg/cluster"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/deployment"
	"vinimpv/loks/pkg/docker"
	"vinimpv/loks/pkg/renderer"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bootstraps the cluster and deploys services",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadUserConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		currentContextRootPath, err := config.GetCurrentContextRootPath()
		if err != nil {
			log.Fatalf("error getting current context root path: %v", err)
		}

		portsToExpose := []int{}
		for _, component := range cfg.Components {
			for _, deployment := range component.Deployments {
				for _, port := range deployment.Ports {
					if port.HostPort != 0 {
						portsToExpose = append(portsToExpose, port.HostPort)
					}
				}
			}
		}

		err = cluster.CreateCluster(cfg.Name, currentContextRootPath, portsToExpose)
		if err != nil {
			log.Fatalf("error creating cluster: %v", err)
		}

		wg := sync.WaitGroup{}
		for _, component := range cfg.Components {
			wg.Add(1)
			go func(component config.Component) {
				if component.SkipPullImage {
					// for the cases where we're skipping pulling the image, we'll build the dev image
					// and load it to the cluster to use it for deployment instead of the production image
					// the dev image will be tagged as localhost/<component>:dev.
					// the user is responsible for providing the build_dev command if necessary
					// if the user doesn't provide the build_dev command, we'll build without any aditional
					// instructions

					devTag := fmt.Sprintf("%s:dev", component.Name)

					if component.BuildDev == "" {
						docker.Build(filepath.Join(currentContextRootPath, component.Name), devTag)
						fmt.Println("build command: ", fmt.Sprintf("docker build %s -t %s", filepath.Join(currentContextRootPath, component.Name), devTag))
					} else {
						docker.BuildDev(filepath.Join(currentContextRootPath, component.Name), component.BuildDev)
						fmt.Println("build dev command: ", component.BuildDev)
					}
					err = cluster.LoadImage(cfg.Name, devTag)
					if err != nil {
						log.Fatalf("error loading dev image to cluster: %s, %v", devTag, err)
					}
				} else {
					err = docker.PullImage(component.Image)
					if err != nil {
						log.Fatalf("error pulling image: %s, %v", component.Image, err)
					}
					err = cluster.LoadImage(cfg.Name, component.Image)
					if err != nil {
						log.Fatalf("error loading dev image to cluster: %s, %v", component.Image, err)
					}
				}
				wg.Done()
			}(component)
		}
		wg.Wait()

		renderedYaml, err := renderer.Render(filepath.Join(currentContextRootPath, "loks.yaml"))
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
	rootCmd.AddCommand(upCmd)
}
