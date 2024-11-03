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

var (
	timeFormats = []string{
		// RFC3339 and variants
		time.RFC3339,
		"2006-01-02T15:04:05Z",      // Basic UTC format
		"2006-01-02T15:04:05-0700",  // With numeric timezone
		"2006-01-02T15:04:05+0700",  // With positive timezone
		"2006-01-02T15:04:05Z07:00", // With colon in timezone
		"2006-01-02T15:04:05+0000",  // UTC with +0000 timezone

		// Common date formats
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"02 Jan 2006 15:04:05",
		"01/02/2006 03:04:05 PM",
		"Mon 02 Jan 2006 15.04 MST",
		"Mon 2 Jan 2006 15.04 MST",
		"January 2, 2006",
		"02 January 2006",

		// News website formats
		"3:04 PM MST Mon January 2 2006",
		"3:04 PM MST Mon January 02 2006",
		"3:04 PM MST January 2 2006",
		"3:04 PM MST January 02 2006",
		"15:04 MST January 2 2006",
		"15:04 MST January 02 2006",
		"January 2, 2006, 3:04 PM",
		"January 02, 2006, 3:04 PM",
		"January 2 2006 3:04 PM",
		"January 02 2006 3:04 PM",
		"2006-01-02",

		// Specific format for your CNN/news case
		"3:04 PM MST Mon October 2 2006",
		"3:04 PM MST Mon October 02 2006",
		"3:04 PM MST Thu October 31 2006", // Exact match for your case
	}
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
	// Preprocess the time string
	timeString = preprocessTimeString(timeString)

	if isVerboseTime(timeString) {
		return parseVerboseTime(timeString)
	}

	// Try each format directly
	for _, format := range timeFormats {
		if t, err := time.Parse(format, timeString); err == nil {
			return t.Unix(), nil
		}
	}

	// Try additional format variations
	variations := []struct {
		process func(string) string
		formats []string
	}{
		{
			// Try with dots converted to colons
			process: func(s string) string { return strings.ReplaceAll(s, ".", ":") },
			formats: timeFormats,
		},
		{
			// Try without commas
			process: func(s string) string { return strings.ReplaceAll(s, ",", "") },
			formats: timeFormats,
		},
		{
			// Try with both transformations
			process: func(s string) string {
				s = strings.ReplaceAll(s, ".", ":")
				return strings.ReplaceAll(s, ",", "")
			},
			formats: timeFormats,
		},
	}

	for _, variation := range variations {
		processedTime := variation.process(timeString)
		for _, format := range variation.formats {
			if t, err := time.Parse(format, processedTime); err == nil {
				return t.Unix(), nil
			}
		}
	}

	return 0, fmt.Errorf("the time %s is not in a valid format", timeString)
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

func preprocessTimeString(input string) string {
	// Remove common prefixes
	cleaned := strings.TrimPrefix(input, "Updated ")
	cleaned = strings.TrimPrefix(cleaned, "Last updated ")
	cleaned = strings.TrimPrefix(cleaned, "Published ")

	// Handle cases where time zone is separated with comma
	cleaned = strings.ReplaceAll(cleaned, ", ", " ")

	// Convert +0000 format to Z format if it matches UTC time
	if strings.HasSuffix(cleaned, "+0000") {
		cleaned = strings.TrimSuffix(cleaned, "+0000") + "Z"
	}

	return strings.TrimSpace(cleaned)
}
