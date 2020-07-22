package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const timeFormat = "2006-01-02"
const boxSubCmd = "box"
const listSubCmd = "list"
const helpSubCmd = "help"

var commands = []string{boxSubCmd, listSubCmd, helpSubCmd}

func appendSlash(path *string) {
	if (*path)[len(*path)-1:] != "/" {
		*path += "/"
	}
}

func checkDateFormat(date string) (time.Time, error) {
	return time.Parse(timeFormat, date)
}

func main() {
	// Setup subcommands
	boxCmd := flag.NewFlagSet(boxSubCmd, flag.ExitOnError)
	boxSeconds := boxCmd.Int("seconds", 0, "How many seconds you want the timebox to last")
	boxMinutes := boxCmd.Int("minutes", 0, "How many minutes you want the timebox to last")
	boxTaskName := boxCmd.String("task", "", "Description of the task you'll be working on during this timebox")
	boxLogPath := boxCmd.String("path", ".", "Path to use when writing the log file")
	boxHelp := boxCmd.Bool("help", false, "Display help text")
	boxCmd.Usage = func() {
		fmt.Println("usage: ancho box {--seconds int|--minutes int|--seconds int --minutes int} [--task string] [--path string]")
		boxCmd.PrintDefaults()
	}

	listCmd := flag.NewFlagSet(listSubCmd, flag.ExitOnError)
	listDate := listCmd.String("date", time.Now().Format(timeFormat), "The date of the timeboxes you want to view")
	listLogPath := listCmd.String("path", ".", "Path to use when looking for log files")
	listHelp := listCmd.Bool("help", false, "Display help text")
	listCmd.Usage = func() {
		fmt.Println("usage: ancho list [--date YYYY-MM-DD]")
		listCmd.PrintDefaults()
	}

	helpCmd := flag.NewFlagSet(helpSubCmd, flag.ExitOnError)
	helpCmd.Usage = func() {
		fmt.Println("usage: ancho <command> [<args>]")
		fmt.Printf("\tbox\t\tStart a new timebox with an optional task name\n")
		fmt.Printf("\tlist\t\tGet a list of timeboxes from a given date (defaults to today's date)\n")
		fmt.Println("See `ancho <command> --help` for command-specific help.")
	}

	subcommand := helpSubCmd
	if len(os.Args) >= 2 {
		subcommand = os.Args[1]
	}

	switch subcommand {
	case boxSubCmd:
		boxCmd.Parse(os.Args[2:])
		appendSlash(boxLogPath)

		if *boxHelp {
			boxCmd.Usage()
			os.Exit(1)
		}

		if *boxSeconds <= 0 && *boxMinutes <= 0 {
			boxCmd.PrintDefaults()
			os.Exit(1)
		}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		done := make(chan bool)

		startTime := time.Now()
		s := time.Duration(*boxSeconds) * time.Second
		m := time.Duration(*boxMinutes) * time.Minute
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
				f, err := os.OpenFile(*boxLogPath+today+".ancho", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}

				toWrite := []byte(startTime.Format(time.RFC3339) + "\t" + endTime.Format(time.RFC3339) + "\t" + *boxTaskName + "\n")
				if _, err := f.Write(toWrite); err != nil {
					f.Close()
					fmt.Println("write error")
					log.Fatal(err)
				}

				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
				return
			case _ = <-ticker.C:
				fmt.Printf(".")
			}
		}
	case listSubCmd:
		listCmd.Parse(os.Args[2:])
		if *listHelp {
			listCmd.Usage()
			os.Exit(1)
		}

		appendSlash(listLogPath)

		if _, err := checkDateFormat(*listDate); err != nil {
			fmt.Printf("Dates should be of format YYYY-MM-DD. Got input: %v\n", *listDate)
			os.Exit(1)
		}

		f, err := os.Open(*listLogPath + *listDate + ".ancho")
		if err != nil {
			fmt.Printf("No log file found for date %v\n", *listDate)
			os.Exit(1)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

	case helpSubCmd:
		helpCmd.Usage()

	default:
		fmt.Printf("No valid subcommand found. Valid commands include %v\n", strings.Join(commands, ","))
		os.Exit(1)
	}
}
