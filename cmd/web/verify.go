package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify the website",
		RunE:  verify,
	}
	cli.AddCommand(cmd)
}

func verify(cmd *cobra.Command, args []string) error {
	return errors.New("verify not implemented")
}
