/*
Copyright © 2025 Scott White
*/
package cmd

import (
	"fmt"

	"github.com/smw1218/sour/generator"
	"github.com/spf13/cobra"
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "create domain templates",
	Long:  `Creates a package within the service that includes a struct on which to add Handlers.`,
	RunE:  CreateDomain,
}

func init() {
	serviceCmd.AddCommand(domainCmd)

	serviceCmd.PersistentFlags().String("domain", "", "a name for the domain, it must be a valid exportable type in go. ex \"PetFood\"")
	serviceCmd.MarkFlagRequired("domain")
}

func CreateDomain(cmd *cobra.Command, args []string) error {
	serviceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	sd, err := NewServiceData(serviceName)
	if err != nil {
		return err
	}

	domainName, err := cmd.Flags().GetString("domain")
	if err != nil {
		return err
	}

	dd, err := generator.NewDomainData(sd, domainName)
	if err != nil {
		return err
	}

	err = dd.CreateDomain()
	if err != nil {
		return err
	}
	fmt.Println("Successfully created", domainName, "domain")
	return nil
}
