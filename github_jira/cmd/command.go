package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getNonEmptyString(command *cobra.Command, name string) (string, error) {
	str, err := command.Flags().GetString(name)
	if err != nil {
		return "", fmt.Errorf("invalid %s parameter", name)
	}
	if str == "" {
		return "", fmt.Errorf("expected %s to not be empty", name)
	}
	return str, nil
}

func getNonEmptyStringSlice(command *cobra.Command, name string) ([]string, error) {
	slice, err := command.Flags().GetStringSlice(name)
	if err != nil {
		return nil, fmt.Errorf("invalid %s parameter", name)
	}
	if len(slice) == 0 {
		return nil, fmt.Errorf("expected %s to not be empty", name)
	}
	return slice, nil
}
