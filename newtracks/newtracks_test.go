package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestLatestTracks(t *testing.T) {
	for _, count := range []int{0, 1, 49, 50, 51, 100} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			tracks := make([]string, count)
			for i := range tracks {
				tracks[i] = fmt.Sprint(i)
			}
			original := slices.Clone(tracks)
			got := latestTracks(tracks)
			if len(got) != min(count, 50) {
				t.Fatalf("got %d tracks", len(got))
			}
			for i, id := range got {
				if id != fmt.Sprint(count-1-i) {
					t.Fatalf("unexpected track at %d: %s", i, id)
				}
			}
			if !slices.Equal(tracks, original) {
				t.Fatal("modified source tracks")
			}
		})
	}
	if !slices.Equal(latestTracks([]string{"A", "B", "A"}), []string{"A", "B", "A"}) {
		t.Fatal("duplicate entries were lost")
	}
}

func TestSyncPlaylist(t *testing.T) {
	cases := []struct {
		name     string
		snapshot string
		result   string
		changed  bool
	}{
		{"unchanged", "S|A,B|D|B,A", "", false},
		{"create", "S|A,B||", "B,A", true},
		{"replace", "S|A,B|D|C", "B,A", true},
		{"reorder", "S|A,B|D|A,B", "B,A", true},
		{"empty source", "S||D|A", "", true},
		{"empty destination", "S||D|", "", false},
		{"create empty", "S|||", "", true},
		{"duplicates", "S|A,B,A|D|B,A", "A,B,A", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			err := syncPlaylist("chill", func(script string) (string, error) {
				calls++
				if calls == 1 {
					return tc.snapshot, nil
				}
				if !strings.Contains(script, "set expectedSnapshot to "+appleScriptString(tc.snapshot)) {
					t.Fatal("missing snapshot check")
				}
				return tc.result, nil
			})
			if err != nil {
				t.Fatal(err)
			}
			wantCalls := 1
			if tc.changed {
				wantCalls = 2
			}
			if calls != wantCalls {
				t.Fatalf("got %d calls, want %d", calls, wantCalls)
			}
		})
	}
}

func TestSyncFailures(t *testing.T) {
	for _, failure := range []string{"read", "malformed response", "replace", "verification"} {
		t.Run(failure, func(t *testing.T) {
			calls := 0
			err := syncPlaylists([]string{"chill", "focus"}, func(script string) (string, error) {
				calls++
				if strings.Contains(script, `set sourceName to "focus"`) {
					t.Fatal("continued after failure")
				}
				if calls == 1 {
					if failure == "read" {
						return "", errors.New("Music unavailable")
					}
					if failure == "malformed response" {
						return "invalid", nil
					}
					return "S|A|D|B", nil
				}
				if failure == "replace" {
					return "", errors.New("Music rejected track")
				}
				return "B", nil
			})
			if err == nil || !strings.Contains(err.Error(), "new chill") {
				t.Fatalf("expected contextual error, got %v", err)
			}
		})
	}
}

func TestMultiplePlaylists(t *testing.T) {
	var scripts []string
	err := syncPlaylists([]string{"chill", "focus"}, func(script string) (string, error) {
		scripts = append(scripts, script)
		return "S||D|", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(scripts) != 2 || !strings.Contains(scripts[1], `set destinationName to "new focus"`) {
		t.Fatal("did not process both playlists")
	}
}

func TestAppleScriptString(t *testing.T) {
	got := appleScriptString("chill \"mix\"\\live\nnext\rline")
	want := `"chill \"mix\"\\live\nnext\rline"`
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestMusicIntegration(t *testing.T) {
	if os.Getenv("NEWTRACKS_MUSIC_TEST") != "1" {
		t.Skip("set NEWTRACKS_MUSIC_TEST=1 to test with temporary Music playlists")
	}
	name := fmt.Sprintf("newtracks test %d", time.Now().UnixNano())
	source := appleScriptString(name)
	destination := appleScriptString("new " + name)
	folder := appleScriptString(name + " folder")
	t.Cleanup(func() {
		_, err := runAppleScript("delete (every user playlist whose name is " + destination + ")\n" +
			"delete (every user playlist whose name is " + source + ")\n" +
			"delete (every folder playlist whose name is " + folder + ")")
		if err != nil {
			t.Error(err)
		}
	})
	_, err := runAppleScript("set fixture to make new user playlist with properties {name:" + source + "}\n" +
		"duplicate track 1 of library playlist 1 to fixture\n" +
		"duplicate track 1 of library playlist 1 to fixture\n" +
		"duplicate track 2 of library playlist 1 to fixture")
	if err != nil {
		t.Fatal(err)
	}
	err = syncPlaylist(name, runAppleScript)
	if err != nil {
		t.Fatal(err)
	}
	destinationID, err := runAppleScript("return persistent ID of user playlist " + destination)
	if err != nil {
		t.Fatal(err)
	}
	checkDestination := func(folderName string) {
		t.Helper()
		output, err := runAppleScript("set targetPlaylist to user playlist " + destination + `
set folderName to ""
try
	set currentParent to parent of targetPlaylist
	if currentParent is not missing value then set folderName to name of currentParent
on error messageText number errorNumber
	if errorNumber is not -1728 then error messageText number errorNumber
end try
return (persistent ID of targetPlaylist) & "|" & folderName`)
		if err != nil {
			t.Fatal(err)
		}
		want := destinationID + "|" + folderName
		if output != want {
			t.Fatalf("destination identity or folder changed: got %q, want %q", output, want)
		}
	}
	checkDestination("")
	for i, folderName := range []string{"", name + " folder"} {
		if folderName != "" {
			_, err = runAppleScript("set fixtureFolder to make new folder playlist with properties {name:" + folder + "}\n" +
				"move user playlist " + destination + " to fixtureFolder")
			if err != nil {
				t.Fatal(err)
			}
		}
		calls := 0
		err = syncPlaylist(name, func(script string) (string, error) {
			calls++
			return runAppleScript(script)
		})
		if err != nil {
			t.Fatal(err)
		}
		if calls != 1 {
			t.Fatal("unchanged playlist was not a no-op")
		}
		checkDestination(folderName)
		_, err = runAppleScript("duplicate track 2 of library playlist 1 to user playlist " + source)
		if err != nil {
			t.Fatal(err)
		}
		err = syncPlaylist(name, runAppleScript)
		if err != nil {
			t.Fatal(err)
		}
		checkDestination(folderName)
		output, err := runAppleScript("return count of tracks of user playlist " + source)
		if err != nil {
			t.Fatal(err)
		}
		if output != fmt.Sprint(4+i) {
			t.Fatalf("source changed: %s tracks", output)
		}
	}
}

func TestLogging(t *testing.T) {
	mode := os.Getenv("NEWTRACKS_LOG_TEST")
	if mode != "" {
		calls := 0
		err := syncPlaylist("chill", func(script string) (string, error) {
			calls++
			if mode == "unchanged" {
				return "S|A,B|D|B,A", nil
			}
			if calls == 1 {
				return "S|A,B|D|C", nil
			}
			if mode == "failure" {
				return "", errors.New("Music failed")
			}
			return "B,A", nil
		})
		if err != nil && mode != "failure" {
			t.Fatal(err)
		}
		return
	}
	cases := map[string][]string{
		"unchanged": {"No updates to new chill needed"},
		"updated":   {"Found latest tracks from chill", "Removed 1 of old tracks", "Add 2 of new tracks"},
		"failure":   {"Found latest tracks from chill"},
	}
	for mode, want := range cases {
		t.Run(mode, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestLogging$")
			cmd.Env = append(os.Environ(), "NEWTRACKS_LOG_TEST="+mode)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%v: %s", err, output)
			}
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			if len(lines) != len(want)+1 || lines[len(lines)-1] != "PASS" {
				t.Fatalf("unexpected output: %s", output)
			}
			for i, message := range want {
				if !strings.HasSuffix(lines[i], message) {
					t.Fatalf("unexpected log: %s", lines[i])
				}
			}
		})
	}
}
