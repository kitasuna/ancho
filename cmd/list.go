package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const timeFormat = "2006-01-02"

var listPath string
var listDate string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List some stuff",
	Long:  "List aaaaaaaaaaall the stuff",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Args: %v\n", args)

		if _, err := checkDateFormat(listDate); err != nil {
			return fmt.Errorf("dates should be of format YYYY-MM-DD. Got input: %v", listDate)
		}

		appendSlash(&listPath)
		f, err := os.Open(listPath + listDate + ".ancho")
		if err != nil {
			return fmt.Errorf("No log file found for date: %v", listDate)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listDate, "date", "d", time.Now().Format(timeFormat), "The date of the log file you want to view")
	listCmd.Flags().StringVarP(&listPath, "path", "p", ".", "Path to use when looking for log files")
}
