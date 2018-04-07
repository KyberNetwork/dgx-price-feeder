package common

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStringToActivityID(t *testing.T) {
	input := "1512189195897392628|1872552297_OMGETH"
	expectedOutput := ActivityID{1512189195897392628, "1872552297_OMGETH"}
	activityID, err := StringToActivityID(input)
	if err != nil {
		t.Fatalf("Expected convert successfully but got error: %v", err)
	} else {
		if activityID != expectedOutput {
			t.Fatalf("Expected %v, got %v", expectedOutput, activityID)
		}
	}
}

func TestStringToActivityIDError(t *testing.T) {
	input := "1512189195897392628_1872552297_OMGETH"
	_, err := StringToActivityID(input)
	if err == nil {
		t.Fatalf("Expected to return error")
	}
}

func TestActivityIDStringable(t *testing.T) {
	id := ActivityID{1512189195897392628, "1872552297_OMGETH"}
	expectedOutput := "1512189195897392628|1872552297_OMGETH"
	output := fmt.Sprintf("%s", id)
	if output != expectedOutput {
		t.Fatalf("Expected %s, got %s", expectedOutput, output)
	}
}

func TestActivityIDToJSON(t *testing.T) {
	id := ActivityID{1512189195897392628, "1872552297_OMGETH"}
	expectedOutput := `"1512189195897392628|1872552297_OMGETH"`
	b, err := json.Marshal(id)
	output := string(b)
	if err != nil {
		t.Fatalf("Expected convert successfully but got error: %v", err)
	} else {
		if output != expectedOutput {
			t.Fatalf("Expected %v, got %v", expectedOutput, output)
		}
	}
}

func TestJSONToActivityID(t *testing.T) {
	expectedOutput := ActivityID{1512189195897392628, "1872552297_OMGETH"}
	input := `"1512189195897392628|1872552297_OMGETH"`
	output := ActivityID{}
	err := json.Unmarshal([]byte(input), &output)
	if err != nil {
		t.Fatalf("Expected convert successfully but got error: %v", err)
	} else {
		if output != expectedOutput {
			t.Fatalf("Expected %v, got %v", expectedOutput, output)
		}
	}
}
