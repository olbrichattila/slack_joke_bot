package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mileusna/crontab"
)

type jokerStruct struct {
	ctab       *crontab.Crontab
	randomJoke RandomJokeInterface
	messanger  MessagerInterface
	skipBefore int
	skipAfter  int
}

var joker *jokerStruct

func newJoker(
	randomJoke RandomJokeInterface,
	messager MessagerInterface,
	skipBefore int,
	skipAfter int,
) *jokerStruct {
	return &jokerStruct{
		randomJoke: randomJoke,
		messanger:  messager,
		skipBefore: skipBefore,
		skipAfter:  skipAfter,
	}
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

func (j *jokerStruct) SendJoke() (bool, error) {
	if j.isJokeAllowed() == false {
		return false, nil
	}
	joke, err := j.NewJoke()
	if err != nil {
		return false, err
	}

	err = j.messanger.Send(*joke)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (j *jokerStruct) NewJoke() (*string, error) {
	joke, err := j.randomJoke.RandomJoke()
	if err != nil {
		return nil, err
	}
	return &joke, nil
}

func (j *jokerStruct) isJokeAllowed() bool {
	currentTime := time.Now()
	currentMinute := currentTime.Hour()*60 + currentTime.Minute()

	fmt.Println(currentMinute)
	if j.skipBefore != 0 && currentMinute < j.skipBefore {
		return false
	}

	if j.skipAfter != 0 && currentMinute > j.skipAfter {
		return false
	}

	return true
}

func sendJoke() {
	sent, err := joker.SendJoke()
	if err != nil {
		fmt.Println(err)
	}

	if sent == false {
		fmt.Printf("Skip sending joke, skipBefore %d, skipAfter %d\n", joker.skipBefore, joker.skipAfter)
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
