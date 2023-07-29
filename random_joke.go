package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type payload struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Joke     string `json:"joke"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}

type RandomJokeInterface interface {
	RandomJoke() (string, error)
}

type RandomJoke struct {
}

func NewRandomJoke() *RandomJoke {
	return &RandomJoke{}
}

func (j *RandomJoke) RandomJoke() (string, error) {
	return j.laodFromAPI()
}

func (j *RandomJoke) laodFromAPI() (string, error) {
	url := "https://v2.jokeapi.dev/joke/Programming"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var joke payload
	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "", err
	}

	if joke.Type == "twopart" {
		return fmt.Sprintf("%s\n\n%s", joke.Setup, joke.Delivery), nil
	}

	return joke.Joke, nil
}
