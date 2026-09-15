package main

import (
	_ "embed"
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/andybrewer/mack"
	"go.mattglei.ch/scripts/internal/logger"
	"go.mattglei.ch/timber"
)

var playlists = []string{
	"chill",
}

//go:embed read.applescript
var readScript string

//go:embed replace.applescript
var replaceScript string

type scriptRunner func(string) (string, error)

type playlistState struct {
	sourceTracks  []string
	destinationID string
	tracks        []string
}

func main() {
	logger.Setup()
	if runtime.GOOS != "darwin" {
		timber.FatalMsg("newtracks requires macOS and Apple Music")
	}
	err := syncPlaylists(playlists, runAppleScript)
	if err != nil {
		timber.Fatal(err, "failed to update playlists")
	}
}

func runAppleScript(script string) (string, error) {
	return mack.Tell("Music", "with timeout of 300 seconds", script, "end timeout")
}

func syncPlaylists(names []string, run scriptRunner) error {
	for _, name := range names {
		err := syncPlaylist(name, run)
		if err != nil {
			return fmt.Errorf("update %q: %w", "new "+name, err)
		}
	}
	return nil
}

func syncPlaylist(name string, run scriptRunner) error {
	destination := "new " + name
	script := "set sourceName to " + appleScriptString(name) + "\n" +
		"set destinationName to " + appleScriptString(destination) + "\n" + readScript
	output, err := run(script + "\nreturn snapshot")
	if err != nil {
		return fmt.Errorf("read playlists: %w", err)
	}
	state, err := parseState(output)
	if err != nil {
		return err
	}
	latest := latestTracks(state.sourceTracks)
	if state.destinationID != "" && slices.Equal(latest, state.tracks) {
		timber.Infof("No updates to %s needed", destination)
		return nil
	}
	timber.Donef("Found latest tracks from %s", name)

	quotedIDs := make([]string, len(latest))
	for i, id := range latest {
		quotedIDs[i] = appleScriptString(id)
	}
	script += "\nset expectedSnapshot to " + appleScriptString(output) + "\n" +
		"set desiredIDs to {" + strings.Join(quotedIDs, ", ") + "}\n" + replaceScript
	output, err = run(script)
	if err != nil {
		return fmt.Errorf("replace playlist contents (the destination may be partially updated; rerun to retry): %w", err)
	}
	if output != strings.Join(latest, ",") {
		return fmt.Errorf("destination tracks do not match the requested order; rerun to retry")
	}
	timber.Donef("Removed %d of old tracks", len(state.tracks))
	timber.Donef("Add %d of new tracks", len(latest))
	return nil
}

func latestTracks(tracks []string) []string {
	latest := slices.Clone(tracks[max(0, len(tracks)-50):])
	slices.Reverse(latest)
	return latest
}

func parseState(output string) (playlistState, error) {
	parts := strings.Split(output, "|")
	if len(parts) != 4 || parts[0] == "" {
		return playlistState{}, fmt.Errorf("unexpected playlist response: %q", output)
	}
	state := playlistState{destinationID: parts[2]}
	if parts[1] != "" {
		state.sourceTracks = strings.Split(parts[1], ",")
	}
	if parts[3] != "" {
		state.tracks = strings.Split(parts[3], ",")
	}
	return state, nil
}

func appleScriptString(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\n", "\\n")
	return "\"" + value + "\""
}
