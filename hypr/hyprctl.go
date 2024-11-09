package hypr

import (
	"os/exec"
	"strings"
)

func Windows() ([]Window, error) {
	cmd := exec.Command("hyprctl", "clients")
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseHyprCtl(b), nil
}

// Parses output from hyprctl clients.
func parseHyprCtl(b []byte) []Window {
	windows := []Window{}

	chunks := splitChunks(string(b), 22)

	for _, chunk := range chunks {
		window := Window{}

		chunk = trim(chunk)

		for _, line := range strings.Split(chunk, "\n") {
			switch prefix(line) {
			case "Window":
				words := strings.Split(line, " ")
				if len(words) < 2 {
					break
				}

				window.Id = words[1]

			case "class:":
				words := strings.Split(line, " ")

				if len(words) < 2 {
					break
				}

				window.Class = strings.Join(words[1:], " ")

			case "title:":
				words := strings.Split(line, " ")

				if len(words) < 2 {
					break
				}

				window.Title = strings.Join(words[1:], " ")

			}
		}

		if window.Class != "" {
			windows = append(windows, window)
		}
	}

	return windows
}

func splitChunks(input string, size uint8) []string {
	lines := strings.Split(input, "\n")
	numLines := len(lines)
	chunkSize := int(size)
	chunks := []string{}

	for i := 0; i < numLines; i += chunkSize {
		end := i + int(size)

		if end > numLines {
			end = numLines
		}

		chunks = append(chunks, strings.Join(lines[i:end], "\n"))
	}

	return chunks
}

func trim(s string) string {
	lines := strings.Split(s, "\n")
	newLines := []string{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		newLines = append(newLines, strings.TrimSpace(line))
	}

	return strings.Join(newLines, "\n")
}

func prefix(s string) string {
	words := strings.Split(s, " ")

	if len(words) > 0 {
		return words[0]
	}

	return ""
}
