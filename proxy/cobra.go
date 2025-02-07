package proxy

import (
	"fmt"
	"net/http"

	"github.com/smw1218/sour/project"
	"github.com/smw1218/sour/service"
	"github.com/spf13/cobra"
)

type ProxyCLI struct {
	rootCmd  *cobra.Command
	services []service.ServiceInterface
}

func NewProxyCLI(services []service.ServiceInterface) *ProxyCLI {
	pc := &ProxyCLI{
		services: services,
	}
	pc.rootCmd = &cobra.Command{
		Use:   "local-proxy",
		Short: "runs microservice reverse proxy",
		Long:  `Runs microservice reverse proxy`,
		RunE:  pc.rootCommandRunner,
	}
	pc.rootCmd.Flags().Int("port", project.DefaultStartingPort, "port to run the proxy")
	for _, svc := range services {
		pc.rootCmd.Flags().String(svc.Name(), "", fmt.Sprintf("url to override destination for %v-service", svc.Name()))
	}
	return pc
}

func (pc *ProxyCLI) Execute() {
	pc.rootCmd.Execute()
}

func (pc *ProxyCLI) rootCommandRunner(cmd *cobra.Command, args []string) error {
	overrides := map[string]string{}
	for _, svc := range pc.services {
		val, err := cmd.Flags().GetString(svc.Name())
		if err != nil {
			return err
		}
		if val != "" {
			overrides[svc.Name()] = val
		}
	}
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	proxy, err := NewProxy(pc.services, overrides)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), proxy)
}
