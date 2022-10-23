package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var filterByName string
	var filterByTags []string
	rootCmd := cobra.Command{
		Use: os.Args[0],
		Run: func(cmd *cobra.Command, args []string) {
			tasks := buildTasksList()
			executeTasks(
				tasks,
				filterByName,
				filterByTags,
			)
		},
	}
	rootCmd.Root().PersistentFlags().StringVar(&filterByName, "name", "", "")
	rootCmd.Root().PersistentFlags().StringSliceVar(&filterByTags, "tags", []string{}, "")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
