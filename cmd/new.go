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

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a new microservice",
	Long:  `Creates a new microservice including stubs for main and service setup.`,
	RunE:  CreateNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	newCmd.PersistentFlags().String("name", "", "a name for the service")
	newCmd.MarkFlagRequired("name")
}

func CreateNew(cmd *cobra.Command, args []string) error {
	serviceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	packageName, err := project.ReadPackage()
	if err != nil {
		return fmt.Errorf("failed getting go package: %w", err)
	}

	lastPort, err := project.GetLastUsedPort()
	if err != nil {
		return err
	}

	sd, err := generator.NewServiceData(serviceName, lastPort+1, packageName)
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
