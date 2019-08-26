package cmd

import (
	"context"
	"net/url"

	"github.com/google/wire"
	"github.com/int128/kauthproxy/pkg/logger"
	"github.com/int128/kauthproxy/pkg/usecases"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var Set = wire.NewSet(
	wire.Struct(new(Cmd), "*"),
	wire.Bind(new(Interface), new(*Cmd)),
)

type Interface interface {
	Run(ctx context.Context, osArgs []string, version string) int
}

type Cmd struct {
	AuthProxy usecases.AuthProxyInterface
	Logger    logger.Interface
}

func (cmd *Cmd) Run(ctx context.Context, osArgs []string, version string) int {
	rootCmd := cmd.newRootCmd(ctx)
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.Version = version

	rootCmd.SetArgs(osArgs[1:])
	if err := rootCmd.Execute(); err != nil {
		cmd.Logger.Printf("error: %s", err)
		return 1
	}
	return 0
}

type rootCmdOptions struct {
	k8sOptions *genericclioptions.ConfigFlags
	address    string
}

func (o *rootCmdOptions) addFlags(f *pflag.FlagSet) {
	o.k8sOptions.AddFlags(f)
	f.StringVar(&o.address, "address", "localhost:0", "The address on which to run the proxy. Default to a random port of localhost.")
}

func (cmd *Cmd) newRootCmd(ctx context.Context) *cobra.Command {
	var o rootCmdOptions
	o.k8sOptions = genericclioptions.NewConfigFlags(false)
	c := &cobra.Command{
		Use:   "kubectl auth-proxy POD_OR_SERVICE_URL",
		Short: "Forward a local port to a pod or service via authentication proxy",
		Long: `Forward a local port to a pod or service via authentication proxy.
To forward a local port to a service, set a service name with .svc suffix. e.g. http://service-name.svc
To forward a local port to a pod, set a pod name. e.g. http://pod-name`,
		Example: `  kubectl auth-proxy https://kubernetes-dashboard.svc`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.runRootCmd(ctx, o, args)
		},
	}
	o.addFlags(c.Flags())
	cmd.Logger.AddFlags(c.PersistentFlags())
	return c
}

func (cmd *Cmd) runRootCmd(ctx context.Context, o rootCmdOptions, args []string) error {
	remoteURL, err := url.Parse(args[0])
	if err != nil {
		return xerrors.Errorf("invalid remote URL: %w", err)
	}
	config, err := o.k8sOptions.ToRESTConfig()
	if err != nil {
		return xerrors.Errorf("could not load the config: %w", err)
	}
	namespace, _, err := o.k8sOptions.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return xerrors.Errorf("could not determine the namespace: %w", err)
	}
	authProxyOptions := usecases.AuthProxyOptions{
		Config:      config,
		Namespace:   namespace,
		TargetURL:   remoteURL,
		BindAddress: o.address,
	}
	if err := cmd.AuthProxy.Do(ctx, authProxyOptions); err != nil {
		return xerrors.Errorf("could not run an authentication proxy: %w", err)
	}
	return nil
}