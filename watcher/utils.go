package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isNixOS() bool {
	cmd := exec.Command("nixos-version")
	err := cmd.Run()
	return err == nil
}

func readJSON[T any](filePath string) (T, error) {
	var data T
	file, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// try to find on the system where chromium is by using which command
func getChromiumPath() string {
	cmd := exec.Command("which", "chromium")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	// remove the newline character from the output
	out = out[:len(out)-1]
	return string(out)
}

func parseTimeAndConvertToUnix(timeString string) (int64, error) {
	if isVerboseTime(timeString) {
		return parseVerboseTime(timeString)
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"02 Jan 2006 15:04:05",
		"01/02/2006 03:04:05 PM",
		"Mon 02 Jan 2006 15.04 MST",
		"January 2, 2006",
		"02 January 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeString); err == nil {
			return t.Unix(), nil
		}
	}
	return 0, fmt.Errorf("unsupported time format")
}

// isVerboseTime checks if the time string is in a verbose format.
func isVerboseTime(timeString string) bool {
	verboseTimeStrings := []string{"day", "hr", "min"}

	// go through each verbose strings and check if it is in the time string
	for _, v := range verboseTimeStrings {
		if strings.Contains(timeString, v) {
			return true
		}
	}

	return false
}

// parseVerboseTime parses a verbose time string and returns the unix timestamp.
func parseVerboseTime(timeString string) (int64, error) {
	fmt.Printf("Parsing verbose time: %s\n", timeString)

	// split the time string into parts
	timeParts := strings.Split(timeString, " ")

	// get the time value and the time unit
	timeValue := timeParts[0]
	timeUnit := timeParts[1]

	// get the current time
	now := time.Now()

	//convert the time value to an integer
	intTimeValue, err := strconv.Atoi(timeValue)
	if err != nil {
		return 0, err
	}

	// if timeUnit is day , do time now - timeValue days
	if strings.Contains(timeUnit, "day") {
		newTime := now.AddDate(0, 0, -1*intTimeValue)
		return newTime.Unix(), nil
	}

	// if timeUnit is hr, do time now - timeValue hours
	if strings.Contains(timeUnit, "hr") {
		newTime := now.Add(-1 * time.Duration(intTimeValue) * time.Hour)
		return newTime.Unix(), nil
	}

	// if timeUnit is min, do time now - timeValue minutes
	if strings.Contains(timeUnit, "min") {
		newTime := now.Add(-1 * time.Duration(intTimeValue) * time.Minute)
		return newTime.Unix(), nil
	}

	return 0, fmt.Errorf("unsupported time unit")
}

func isExistant(link string) bool {
	for _, a := range Articles {
		if a.Link == link {
			return true
		}
	}
	return false
}
