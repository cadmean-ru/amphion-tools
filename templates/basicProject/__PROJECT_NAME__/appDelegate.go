package main

import (
	"github.com/cadmean-ru/amphion/engine"
)

type AppDelegate struct {
	engine.AppDelegateImpl
}

func (d *AppDelegate) OnAppLoaded() {
	engine.LogDebug("App loaded")
}