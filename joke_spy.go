package main

import "fmt"

type RandomJokeSpy struct {
	calledRandomJoke int
	callText         string
	err              error
}

type MessagerSpy struct {
	calledSend int
	sentText   string
	err        error
}

func newRandomJokeSpy() *RandomJokeSpy {
	return &RandomJokeSpy{}
}

func newMessageSpy() *MessagerSpy {
	return &MessagerSpy{}
}

func (s *RandomJokeSpy) RandomJoke() (string, error) {
	s.calledRandomJoke++
	return s.callText, s.err
}

func (s *RandomJokeSpy) calledWith(text string) *RandomJokeSpy {
	s.callText = text
	return s
}

func (s *RandomJokeSpy) WithError() *RandomJokeSpy {
	s.err = fmt.Errorf("Cutom error")
	return s
}

func (s *MessagerSpy) Send(text string) error {
	s.calledSend++
	s.sentText = text
	return s.err
}

func (s *MessagerSpy) WithError() *MessagerSpy {
	s.err = fmt.Errorf("Cutom error")
	return s
}
