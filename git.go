package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const configGit = "git"

func useGit(path string) (bool, error) {
	c, err := getConfig(path)
	if err != nil {
		return false, err
	}

	if c[configGit] == "" {
		// not set at all
		return false, nil
	}

	return strconv.ParseBool(c[configGit])
}

func syncPull(path string) error {
	g, err := useGit(path)
	if err != nil {
		return err
	}
	if !g {
		return nil
	}

	cmd := exec.Command("git", "pull")
	cmd.Env = os.Environ()
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sync pull: %w: %v", err, string(output))
	}
	return nil
}

func syncPush(path string) error {
	g, err := useGit(path)
	if err != nil {
		return err
	}
	if !g {
		return nil
	}

	cmd := exec.Command("git", "add", "-A")
	cmd.Env = os.Environ()
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sync add: %w: %v", err, string(output))
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("tagebuch %v", time.Now()))
	cmd.Env = os.Environ()
	cmd.Dir = path
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sync commit: %w: %v", err, string(output))
	}

	cmd = exec.Command("git", "push")
	cmd.Env = os.Environ()
	cmd.Dir = path
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sync push: %w: %v", err, string(output))
	}

	return nil
}
