package main

import (
	"time"

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

	if j.skipBefore != 0 && currentMinute < j.skipBefore {
		return false
	}

	if j.skipAfter != 0 && currentMinute > j.skipAfter {
		return false
	}

	return true
}
