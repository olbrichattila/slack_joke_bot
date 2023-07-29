package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) TestJokeSent() {
	expectedJoke := "This is the joke"
	jokeSpy := newRandomJokeSpy().calledWith(expectedJoke)
	messageSpy := newMessageSpy()

	joker := newJoker(jokeSpy, messageSpy, 0, 0)

	joker.SendJoke()

	t.Equal(1, jokeSpy.calledRandomJoke)
	t.Equal(1, messageSpy.calledSend)
	t.Equal(expectedJoke, messageSpy.sentText)

	// Assert new joke sent
	expectedNewJoke := "This is the new joke"
	jokeSpy.calledWith(expectedNewJoke)
	joker.SendJoke()

	t.Equal(2, jokeSpy.calledRandomJoke)
	t.Equal(2, messageSpy.calledSend)
	t.Equal(expectedNewJoke, messageSpy.sentText)
}

func (t *TestSuite) TestJokeSkipped() {
	jokeSpy := newRandomJokeSpy()
	messageSpy := newMessageSpy()

	// Note: this test can be flaky if it is executed on x hour 0 minutes
	now := time.Now()
	skipAfterStr := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()-1)
	skipAfter, _ := hourMinuteToMinutes(skipAfterStr)
	joker := newJoker(jokeSpy, messageSpy, 0, skipAfter)

	sent, _ := joker.SendJoke()

	t.False(sent)
}

func (t *TestSuite) TestErrorReportedIfCannotGetJoke() {
	jokeSpy := newRandomJokeSpy().WithError()
	messageSpy := newMessageSpy()

	joker := newJoker(jokeSpy, messageSpy, 0, 0)

	sent, err := joker.SendJoke()

	t.False(sent)
	t.NotNil(err)
}

func (t *TestSuite) TestErrorReportedIfCannotSendJoke() {
	jokeSpy := newRandomJokeSpy()
	messageSpy := newMessageSpy().WithError()

	joker := newJoker(jokeSpy, messageSpy, 0, 0)

	sent, err := joker.SendJoke()

	t.False(sent)
	t.NotNil(err)
}
