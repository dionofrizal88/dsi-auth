package cmd_test

import (
	"fmt"
	"github.com/dionofrizal88/dsi/auth/interfaces/cmd"
	"github.com/dionofrizal88/dsi/auth/pkg/tests"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"os"
	"testing"
)

func TestCMDAppCli(t *testing.T) {
	testSuite := tests.InitTestSuite()

	t.Cleanup(func() {
		testSuite.Clean()
	})

	t.Run("positive case to test cmd app, expected no error", func(t *testing.T) {
		t.Run("positive case while use func init cmd app, expected no error", func(t *testing.T) {
			os.Args = []string{"test", "--name", "Jeremy"}
			app := cmd.NewCli()
			app.Name = "test"
			app.Flags = []cli.Flag{
				&cli.StringFlag{Name: "name", Value: "dio", Usage: "a name to say"},
			}
			app.Action = func(c *cli.Context) error {
				fmt.Printf("Hello %v\n", c.String("name"))
				return nil
			}

			app.UsageText = "app [first_arg] [second_arg]"
			app.Authors = []*cli.Author{{Name: "Dio Agus Nofrizal", Email: "dio@gmail.com"}}

			err := app.Run(os.Args)

			require.NoError(t, err)
		})
	})

}
