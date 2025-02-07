/*
Copyright © 2025 Scott White
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/smw1218/sour/generator"
	"github.com/smw1218/sour/project"
	"github.com/spf13/cobra"
)

// serviceCmd represents the new command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "creates and modifies a microservice",
	Long:  `Creates a new microservice including stubs for main and service setup.`,
	RunE:  CreateNew,
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serviceCmd.PersistentFlags().String("name", "", "a name for the service")
	serviceCmd.MarkFlagRequired("name")
}

func CreateNew(cmd *cobra.Command, args []string) error {
	serviceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	sd, err := NewServiceData(serviceName)
	if err != nil {
		return err
	}

	// create service dir
	err = os.MkdirAll(sd.ServiceDirectory(), 0755)
	if err != nil {
		return fmt.Errorf("failed creating dir %v: %w", sd.ServiceDirectory(), err)
	}
	err = sd.CreateService()
	if err != nil {
		return err
	}

	fmt.Println("Successfully created", serviceName, "service")
	return nil
}

func NewServiceData(serviceName string) (*generator.ServiceData, error) {
	packageName, err := project.ReadPackage()
	if err != nil {
		return nil, fmt.Errorf("failed getting go package: %w", err)
	}

	lastPort, err := project.GetLastUsedPort()
	if err != nil {
		return nil, err
	}

	return generator.NewServiceData(serviceName, lastPort+1, packageName)
}
