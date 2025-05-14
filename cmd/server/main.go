package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	for {
		time.Sleep(1000)
    	logrus.Info("Hello, Logrus!")
	}
}