package cmd

import (
	"github.com/urfave/cli/v2"
)

// NewCli is a constructor will initialize cli.
func NewCli() *cli.App {
	c := cli.NewApp()

	return c
}
