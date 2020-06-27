package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const timeFormat = "2006-01-02"
const boxSubCmd = "box"
const listSubCmd = "list"

//var ErrorMessages map[string]string
//ErrorMessages["huh"] = "this wont work"

func main() {
	// Setup subcommands
	timerCmd := flag.NewFlagSet(boxSubCmd, flag.ExitOnError)
	taskSeconds := timerCmd.Int("seconds", 0, "How many seconds you want the timebox to last")
	taskMinutes := timerCmd.Int("minutes", 0, "How many minutes you want the timebox to last")
	taskName := timerCmd.String("task", "", "Description of the task you'll be working on during this timebox")

	listCmd := flag.NewFlagSet(listSubCmd, flag.ExitOnError)
	listDate := listCmd.String("date", time.Now().Format(timeFormat), "The date of the timeboxes you want to view")

	if len(os.Args) < 2 {
		log.Fatal("No valid subcommand found. Valid commands include `box` and `list`")
	}

	switch os.Args[1] {
	case boxSubCmd:
		timerCmd.Parse(os.Args[2:])
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		done := make(chan bool)

		startTime := time.Now()
		fmt.Printf("Starting at (RFC3339): %v\n", startTime.Format(time.RFC3339))
		s := time.Duration(*taskSeconds) * time.Second
		m := time.Duration(*taskMinutes) * time.Minute
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
				// Current date
				today := endTime.Format(timeFormat)
				f, err := os.OpenFile(today+".ancho", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}

				toWrite := []byte(startTime.Format(time.RFC3339) + "\t" + endTime.Format(time.RFC3339) + "\t" + *taskName + "\n")
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
		f, err := os.Open(*listDate + ".ancho")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

	default:
		log.Fatal("No valid subcommand found. Valid commands include `box` and `list`")
	}
}
