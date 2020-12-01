package cmd

import (
	"errors"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

const SecondsFlag = "seconds"
const MinutesFlag = "minutes"
const LabelFlag = "label"
const PathFlag = "path"

var plan = "ridge"

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "Start a timebox",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		boxSeconds := viper.GetInt("plans." + plan + "." + SecondsFlag)
		boxMinutes := viper.GetInt("plans." + plan + "." + MinutesFlag)
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
				boxPath := viper.GetString("path")
				appendSlash(&boxPath)
				f, err := os.OpenFile(boxPath+today+".ancho", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}

				label, _ := cmd.Flags().GetString("label")
				toWrite := []byte(startTime.Format(time.RFC3339) + "\t" + endTime.Format(time.RFC3339) + "\t" + label + "\n")
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

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir := home + "/.config/ancho"
	configFileName := "ancho"
	configFileExt := "yaml"
	fullConfigPath := fmt.Sprintf("%v/%v.%v", configDir, configFileName, configFileExt)

	// Create config directory if it doesn't already exist
	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, os.ModeDir|0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	// Create (touch) config file if it doesn't already exist
	// TODO add some defaults here later
	if _, err = os.Stat(fullConfigPath); os.IsNotExist(err) {
		fmt.Printf("... Creating config file at %v\n", fullConfigPath)
		file, err := os.OpenFile(fullConfigPath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		file.Close()
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigType(configFileExt)
	viper.SetConfigName(configFileName)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(boxCmd)
	boxCmd.Flags().IntP(SecondsFlag, "s", 0, fmt.Sprintf("Number of seconds you want the timebox to last. Can be combined with --%v", MinutesFlag))
	viper.BindPFlag("plans."+plan+"."+SecondsFlag, boxCmd.Flags().Lookup(SecondsFlag))

	boxCmd.Flags().IntP(MinutesFlag, "m", 0, fmt.Sprintf("Number of minutes you want the timebox to last. Can be combined with --%v", SecondsFlag))
	viper.BindPFlag("plans."+plan+"."+MinutesFlag, boxCmd.Flags().Lookup(MinutesFlag))

	boxCmd.Flags().StringP(PathFlag, "p", ".", "Path to use when writing the log file")
	viper.BindPFlag(PathFlag, boxCmd.Flags().Lookup(PathFlag))

	boxCmd.Flags().StringP(LabelFlag, "l", "", "Label for the task you'll be working on during this timebox")
}
