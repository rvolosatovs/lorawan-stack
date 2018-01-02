// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package main

import (
	"fmt"
	"os"

	"github.com/TheThingsNetwork/ttn/cmd/shared"
	"github.com/TheThingsNetwork/ttn/pkg/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// Config contains the base configuration with shared flags.
type Config struct {
	config.Base `name:",squash"`
	Int         int    `name:"int" description:"An example int"`
	String      string `name:"string" description:"An example string"`
}

// SubConfig contains the flags for the sub command.
type SubConfig struct {
	Config `name:",squash"`
	Bar    string `name:"bar" description:"The bar config flag"`
}

var (
	// mgr is a configuration manager for the example program, initialized
	// with default config values.
	mgr = config.InitializeWithDefaults("example", &Config{
		Base:   shared.DefaultBaseConfig,
		Int:    42,
		String: "foo",
	})

	// cfg will be the config shared by all commands.
	cfg = new(Config)

	// root is the root command.
	root = &cobra.Command{
		// PersistentPreRunE runs always, so we can use it to read in config files
		// and parse the shared config.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := mgr.ReadInConfig()
			if err != nil {
				return err
			}

			return mgr.Unmarshal(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return printYAML(cfg)
		},
	}

	// sub is a sub command.
	sub = &cobra.Command{
		Use: "sub",

		// RunE parses the sub config that is only relevant for this command.
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := new(SubConfig)
			err := mgr.Unmarshal(cfg)
			if err != nil {
				return err
			}

			return printYAML(cfg)
		},
	}
)

func init() {
	root.AddCommand(sub)
	root.Flags().AddFlagSet(mgr.Flags())
	sub.Flags().AddFlagSet(mgr.WithConfig(&SubConfig{
		Bar: "baz",
	}))
}

func main() {
	// execute the commands and print resulting errors.
	err := root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// printYAML prints the nested config struct.
func printYAML(in interface{}) error {
	bs, err := yaml.Marshal(in)
	if err != nil {
		return err
	}

	fmt.Print(string(bs))
	return nil
}
