package cmd

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/gobuffalo/buffalo-plugins/genny/add"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var addOptions = struct {
	dryRun bool
}{}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds plugins to config/buffalo-plugins.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		run := genny.WetRunner(context.Background())
		if addOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil && (errors.Cause(err) != plugdeps.ErrMissingConfig) {
			return errors.WithStack(err)
		}

		for _, a := range args {
			a = strings.TrimSpace(a)
			bin := path.Base(a)
			plug := plugdeps.Plugin{
				Binary: bin,
				GoGet:  a,
			}
			if _, err := os.Stat(a); err == nil {
				plug.Local = a
				plug.GoGet = ""
			}
			plugs.Add(plug)
		}
		g, err := add.New(&add.Options{
			App:     app,
			Plugins: plugs.List(),
		})
		if err != nil {
			return errors.WithStack(err)
		}
		run.With(g)

		return run.Run()
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addOptions.dryRun, "dry-run", "d", false, "dry run")
}
