package loks

import (
	"fmt"
	"log"
	"path/filepath"
	"vinimpv/loks/pkg/config"
	"vinimpv/loks/pkg/docker"
	"vinimpv/loks/pkg/gefyra"

	"github.com/spf13/cobra"
)

func getDeploymentAndComponentByDeploymentName(name string, cfg *config.Config) (config.Deployment, config.Component, error) {
	for _, c := range cfg.Components {
		for _, d := range c.Deployments {
			if d.Name == name {
				return d, c, nil
			}
		}
	}
	return config.Deployment{}, config.Component{}, fmt.Errorf("deployment %s not found", name)
}

var devCmd = &cobra.Command{
	Use:   "dev [deploymentName]",
	Short: "Enter development mode for a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadUserConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		deploymentName := args[0]
		dep, comp, err := getDeploymentAndComponentByDeploymentName(deploymentName, cfg)
		if err != nil {
			log.Fatalf("error getting deployment: %v", err)
		}
		if comp.BuildDev == "" {
			log.Fatalf("build_dev command not provided for %s", comp.Name)
		}
		currentContextRootPath, err := config.GetCurrentContextRootPath()
		if err != nil {
			log.Fatalf("error getting current context root path: %v", err)
		}
		if !docker.CheckImageExists(fmt.Sprintf("%s:dev", comp.Name)) {
			log.Println("Building dev image")
			err = docker.BuildDev(filepath.Join(currentContextRootPath, comp.Name), comp.BuildDev)
			if err != nil {
				log.Fatalf("error building dev image: %v", err)
			}
		}
		err = gefyra.Start(cfg.Name)
		if err != nil {
			log.Fatalf("error starting gefyra: %v", err)
		}
		devImageName := fmt.Sprintf("%s:dev", comp.Name)
		devContainerName := fmt.Sprintf("%s-%s", comp.Name, dep.Name)
		ports := []string{fmt.Sprintf("%d:%d", dep.Port, dep.Port)}
		volumes := []string{fmt.Sprintf("%s:/app", filepath.Join(currentContextRootPath, comp.Name))}

		err = gefyra.RunContainer(
			devImageName,
			devContainerName,
			"default",
			ports,
			volumes,
		)
		if err != nil {
			log.Fatalf("error running gefyra container: %v", err)
		}

		err = gefyra.Bridge(
			devImageName,
			devContainerName,
			dep.Name,
			"default",
			ports,
			volumes,
		)
		if err != nil {
			log.Fatalf("error running gefyra bridge: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
