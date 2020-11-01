package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

var cfgFile string

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
	Short: "A configurable pomodoro / timeboxing CLI application",
	Long:  "A configurable pomodoro / timeboxing CLI application",
}

// Executes the root command
func Execute() error {
	return rootCmd.Execute()
}
