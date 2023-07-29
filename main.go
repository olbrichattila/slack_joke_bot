package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mileusna/crontab"
)

func main() {
	loadEnv()
	randomJoke := NewRandomJoke()
	messager := NewMessager()
	skipBefore, skipAfter, err := parseSkipParameters()
	if err != nil {
		fmt.Println(err)
		return
	}

	joker = newJoker(randomJoke, messager, skipBefore, skipAfter)
	sendJoke()

	if err := initCron(); err != nil {
		fmt.Println(err)
		return
	}

	waitForTermSignal()
}

func hourMinuteToMinutes(input string) (int, error) {
	timeComponents := strings.Split(input, ":")
	if len(timeComponents) != 2 {
		return 0, fmt.Errorf("invalid input format, expected 'hour:minute'")
	}

	hour, err := strconv.Atoi(timeComponents[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse hour: %v", err)
	}

	minute, err := strconv.Atoi(timeComponents[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse minute: %v", err)
	}

	return hour*60 + minute, nil
}

func minutesToString(minutes int) string {
	hour := int(math.Floor(float64(minutes) / 60))
	minute := minutes % 60

	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func initCron() error {
	cron := os.Getenv("CRON")
	ctab := crontab.New()

	return ctab.AddJob(cron, sendJoke)
}

func waitForTermSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Program is running. Press Ctrl+C (SIGINT) or send SIGTERM to stop.")
	<-sigChan
}

func sendJoke() {
	sent, err := joker.SendJoke()
	if err != nil {
		fmt.Println(err)
	}

	if sent == false {
		fmt.Printf(
			"Skip sending joke, skipBefore %s, skipAfter %s\n",
			minutesToString(joker.skipBefore),
			minutesToString(joker.skipAfter),
		)
	} else {
		fmt.Println("Joke successfully delivered to the channel")
	}
}

func parseSkipParameters() (int, int, error) {
	skipBefore, err := hourMinuteToMinutes(os.Getenv("SKIP_BEFORE"))
	if err != nil {
		return 0, 0, err
	}

	skipAfter, err := hourMinuteToMinutes(os.Getenv("SKIP_AFTER"))
	if err != nil {
		return 0, 0, err
	}

	return skipBefore, skipAfter, nil
}

func loadEnv() error {
	if fileExists("./.env") {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
