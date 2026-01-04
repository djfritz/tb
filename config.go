package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getConfig(path string) (map[string]string, error) {
	configPath := filepath.Join(path, tagebuchMagic)
	c := make(map[string]string)

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		kv := strings.Split(text, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid config item: %v", text)
		}
		c[kv[0]] = kv[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return c, nil
}
