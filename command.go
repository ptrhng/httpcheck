package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewCommand creates a new httpcheck command.
func NewCommand() *cobra.Command {
	opts := NewDefaultOptions()

	cmd := &cobra.Command{
		Use:   `httpcheck [METHOD] URL [REQUEST ITEM...]`,
		Short: "Measuring HTTP performance",
		Example: `httpcheck www.example.com
httpcheck POST www.example.com colors:='["red", "green", "blue"]'`,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetLevel(logrus.FatalLevel)

			if err := ParseArgs(args, opts); err != nil {
				return err
			}

			r, err := Trace(cmd.Context(), opts)
			if err != nil {
				return err
			}

			return PrintResult(r, WithShowBody(opts.ShowBody), WithMaxBodySize(opts.maxBodySize))
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.ShowBody, "body", "b", false, "print response body")
	flags.BoolVarP(&opts.FollowRedirect, "follow", "F", false, "follow redirects")

	return cmd
}
