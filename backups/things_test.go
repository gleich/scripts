package main

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestThingsUpdateURLPreservesSpacesAndPlusSigns(t *testing.T) {
	data := []byte(`{"title":"Edited photos + café & 100%"}`)
	link, err := url.Parse(thingsUpdateURL("test+token", data))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(link.RawQuery, "+") {
		t.Fatal("Things requires percent-encoded spaces")
	}
	want := map[string]string{"auth-token": "test+token", "data": string(data)}
	for _, parameter := range strings.Split(link.RawQuery, "&") {
		key, value, _ := strings.Cut(parameter, "=")
		decoded, err := url.PathUnescape(value)
		if err != nil {
			t.Fatal(err)
		}
		if decoded != want[key] {
			t.Fatalf("%s decoded to %q, want %q", key, decoded, want[key])
		}
		delete(want, key)
	}
	if len(want) != 0 {
		t.Fatalf("missing parameters: %v", want)
	}
}

func TestCompleteChecklist(t *testing.T) {
	items := []checklistItem{
		{Title: "Strava", Status: 3},
		{Title: "GitHub", Status: 0},
		{Title: "Edited photos", Status: 0},
		{Title: "Skipped", Status: 2},
	}
	original := append([]checklistItem(nil), items...)
	checklist, changed, err := completeChecklist(items, "GitHub")
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected a checklist update")
	}
	want := []thingsChecklistItem{
		{Type: "checklist-item", Attributes: checklistAttributes{Title: "Strava", Completed: true}},
		{Type: "checklist-item", Attributes: checklistAttributes{Title: "GitHub", Completed: true}},
		{Type: "checklist-item", Attributes: checklistAttributes{Title: "Edited photos"}},
		{Type: "checklist-item", Attributes: checklistAttributes{Title: "Skipped", Canceled: true}},
	}
	if !reflect.DeepEqual(checklist, want) {
		t.Fatalf("got %#v, want %#v", checklist, want)
	}
	if !reflect.DeepEqual(items, original) {
		t.Fatal("changed the input checklist")
	}

	items[1].Status = 3
	_, changed, err = completeChecklist(items, "GitHub")
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("already completed item should not require an update")
	}
}

func TestCompleteChecklistRejectsInvalidItems(t *testing.T) {
	cases := map[string][]checklistItem{
		"missing":   {{Title: "Strava"}},
		"duplicate": {{Title: "GitHub"}, {Title: "GitHub"}},
		"canceled":  {{Title: "GitHub", Status: 2}},
		"unknown":   {{Title: "GitHub", Status: 1}},
		"too many":  make([]checklistItem, 101),
	}
	for name, items := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := completeChecklist(items, "GitHub")
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
