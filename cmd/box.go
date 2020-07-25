package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	boxSeconds int
	boxMinutes int
	boxLabel   string
	boxPath    string
)

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "Start a timebox",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Args: %v\n", args)

		if boxSeconds <= 0 && boxMinutes <= 0 {
			return errors.New("Either minutes or seconds must be greater than 0")
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		done := make(chan bool)

		startTime := time.Now()
		s := time.Duration(boxSeconds) * time.Second
		m := time.Duration(boxMinutes) * time.Minute
		endTime := startTime.Add(s).Add(m)

		fmt.Printf("Starting at (RFC3339): %v\n", startTime.Format(time.RFC3339))
		fmt.Printf("Scheduled end at (RFC3339): %v\n", endTime.Format(time.RFC3339))

		go func() {
			time.Sleep(s + m)
			done <- true
		}()

		for {
			select {
			case <-done:
				endTime := time.Now()
				fmt.Printf("\nEnding at (RFC3339): %v\n", endTime.Format(time.RFC3339))

				// Open file for logging
				today := endTime.Format(timeFormat)
				appendSlash(&boxPath)
				f, err := os.OpenFile(boxPath+today+".ancho", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}

				toWrite := []byte(startTime.Format(time.RFC3339) + "\t" + endTime.Format(time.RFC3339) + "\t" + boxLabel + "\n")
				if _, err := f.Write(toWrite); err != nil {
					f.Close()
					return err
				}

				if err := f.Close(); err != nil {
					return err
				}
				return nil

			case _ = <-ticker.C:
				fmt.Printf(".")
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(boxCmd)
	boxCmd.Flags().IntVarP(&boxSeconds, "seconds", "s", 0, "Number of seconds you want the timebox to last. Can be combined with `--minutes`")
	boxCmd.Flags().IntVarP(&boxMinutes, "minutes", "m", 0, "Number of minutes you want the timebox to last. Can be combined with `--seconds`")
	boxCmd.Flags().StringVarP(&boxLabel, "label", "l", "", "Label for the task you'll be working on during this timebox")
	boxCmd.Flags().StringVarP(&boxPath, "path", "p", ".", "Path to use when writing the log file")
}
