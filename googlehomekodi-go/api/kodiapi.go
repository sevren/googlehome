package api

import (
	log "github.com/sirupsen/logrus"
)

// Simple Kodi Api functions

func SyncLibrary() {
	log.Info("Executing syncLibrary")
}

func NavSelect() {
	log.Info("Selecting the current item")
}

// Text Kodi Api functions
func HandleTvShow(s string) {
	log.Infof("HandlingTVShow - Params: %s", s)
}

// Text and Number Kodi api functions

func PlayEpisode(s string, n int) {
	log.Infof("Attempting to play an episode of  %s %d", s, n)
}
