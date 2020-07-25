package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

func appendSlash(path *string) {
	if (*path)[len(*path)-1:] != "/" {
		*path += "/"
	}
}

func checkDateFormat(date string) (time.Time, error) {
	return time.Parse(timeFormat, date)
}

var rootCmd = &cobra.Command{
	Use:   "ancho",
	Short: "brief description of ancho",
	Long:  "nice sexy description of ancho",
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
