package dip

import (
	"errors"
	"fmt"
)

type Logger interface {
	Log(message string)
}

type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

func LogOutPut(message string) {
	fmt.Println(message)
}

type SimpleLogic struct {
	Logger
	DataStore
}

func (sl SimpleLogic) SayHello(userId string) (string, error) {
	sl.Log("In SayHello for " + userId)
	name, ok := sl.UserNameForId(userId)
	if !ok {
		return "", errors.New("user not found")
	}
	return "Hello, " + name, nil
}

func (sl SimpleLogic) SayGoodBye(userId string) (string, error) {
	sl.Log("inside SayGoodBye " + userId)
	name, ok := sl.UserNameForId(userId)
	if !ok {
		return "", errors.New("user not found")
	}

	return "GoodBye, " + name, nil
}
