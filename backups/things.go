package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const backupsTaskScript = `
tell application "Things3"
	set matches to to dos of list "Today" whose name is "Backups" and status is open
	if (count of matches) is not 1 then
		error "Expected exactly one open Backups task in Today"
	end if
	return id of item 1 of matches
end tell
`

type checklistItem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
}

type thingsChecklistItem struct {
	Type       string              `json:"type"`
	Attributes checklistAttributes `json:"attributes"`
}

type checklistAttributes struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Canceled  bool   `json:"canceled"`
}

func completeBackupChecklist(home, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, "osascript", "-e", backupsTaskScript).CombinedOutput()
	if err != nil {
		return fmt.Errorf("find Backups task: %w: %s", err, strings.TrimSpace(string(output)))
	}
	taskID := strings.TrimSpace(string(output))

	databases, err := filepath.Glob(filepath.Join(home,
		"Library/Group Containers/JLMPQHK86H.com.culturedcode.ThingsMac",
		"ThingsData-*", "Things Database.thingsdatabase", "main.sqlite"))
	if err != nil {
		return err
	}
	if len(databases) != 1 {
		return fmt.Errorf("expected one Things database, found %d", len(databases))
	}
	database := databases[0]
	items, err := readChecklist(ctx, database, taskID)
	if err != nil {
		return err
	}
	checklist, changed, err := completeChecklist(items, name)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	token, err := thingsAuthToken(ctx)
	if err != nil {
		return err
	}
	data, err := json.Marshal([]any{map[string]any{
		"type":      "to-do",
		"operation": "update",
		"id":        taskID,
		"attributes": map[string]any{
			"checklist-items": checklist,
		},
	}})
	if err != nil {
		return err
	}
	err = exec.CommandContext(ctx, "open", "-g", thingsUpdateURL(token, data)).Run()
	if err != nil {
		return fmt.Errorf("send checklist update to Things: %w", err)
	}

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("Things did not confirm checklist completion; check that Things URLs are enabled and the authorization token is valid")
		case <-ticker.C:
			items, err = readChecklist(ctx, database, taskID)
			if err != nil {
				return err
			}
			_, changed, err = completeChecklist(items, name)
			if err != nil {
				return err
			}
			if !changed {
				return nil
			}
		}
	}
}

func thingsUpdateURL(token string, data []byte) string {
	query := url.Values{
		"auth-token": {token},
		"data":       {string(data)},
	}
	return "things:///json?" + strings.ReplaceAll(query.Encode(), "+", "%20")
}

func thingsAuthToken(ctx context.Context) (string, error) {
	token := strings.TrimSpace(os.Getenv("THINGS_AUTH_TOKEN"))
	if token != "" {
		return token, nil
	}
	output, err := exec.CommandContext(ctx, "/usr/bin/security", "find-generic-password",
		"-s", "go.mattglei.ch.scripts.backups.things", "-a", "things-auth-token", "-w").Output()
	if err != nil {
		return "", fmt.Errorf("read Things authorization token from Keychain (see backups/README.md for setup): %w", err)
	}
	token = strings.TrimSpace(string(output))
	if token == "" {
		return "", errors.New("Things authorization token in Keychain is empty")
	}
	return token, nil
}

func readChecklist(ctx context.Context, database, taskID string) ([]checklistItem, error) {
	query := `SELECT title, status FROM TMChecklistItem WHERE task = '` +
		strings.ReplaceAll(taskID, "'", "''") + `' ORDER BY "index"`
	output, err := exec.CommandContext(ctx, "sqlite3", "-readonly", "-json", database, query).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("read Things checklist: %w: %s", err, strings.TrimSpace(string(output)))
	}
	if len(output) == 0 {
		return nil, errors.New("Backups task has no checklist")
	}
	var items []checklistItem
	err = json.Unmarshal(output, &items)
	if err != nil {
		return nil, fmt.Errorf("decode Things checklist: %w", err)
	}
	return items, nil
}

func completeChecklist(items []checklistItem, name string) ([]thingsChecklistItem, bool, error) {
	if len(items) > 100 {
		return nil, false, errors.New("Things URL updates support at most 100 checklist items")
	}
	checklist := make([]thingsChecklistItem, 0, len(items))
	matches := 0
	changed := false
	for _, item := range items {
		if item.Status != 0 && item.Status != 2 && item.Status != 3 {
			return nil, false, fmt.Errorf("unknown checklist status %d for %q", item.Status, item.Title)
		}
		attributes := checklistAttributes{
			Title:     item.Title,
			Completed: item.Status == 3,
			Canceled:  item.Status == 2,
		}
		if item.Title == name {
			matches++
			if attributes.Canceled {
				return nil, false, fmt.Errorf("checklist item %q is canceled", name)
			}
			changed = !attributes.Completed
			attributes.Completed = true
		}
		checklist = append(checklist, thingsChecklistItem{
			Type:       "checklist-item",
			Attributes: attributes,
		})
	}
	if matches != 1 {
		return nil, false, fmt.Errorf("expected one checklist item named %q, found %d", name, matches)
	}
	return checklist, changed, nil
}
