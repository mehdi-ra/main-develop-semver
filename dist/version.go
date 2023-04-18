package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
)

const (
	developBranch  = "develop"
	mainBranch     = "main"
	breakingPrefix = "BREAKING CHANGE:"
)

var (
	errGitCommandFailed = errors.New("git command failed")
	errInvalidVersion   = errors.New("invalid version")
)

func cleanVersion(s string) (string, error) {
	// Check if the input string is empty
	if s == "" {
		return "", fmt.Errorf("empty string")
	}

	// Remove any leading or trailing spaces
	s = strings.TrimSpace(s)

	// Remove any prefix that starts with "v"
	s = strings.TrimPrefix(s, "v")

	// Remove any suffix that starts with "-"
	if i := strings.Index(s, "-"); i != -1 {
		s = s[:i]
	}

	// Check if the resulting string is empty
	if s == "" {
		return "", fmt.Errorf("invalid version number")
	}

	return s, nil
}

func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%w: %s", errGitCommandFailed, err.Error())
	}
	return strings.TrimSpace(string(output)), nil
}

func getCurrentVersion() (string, error) {
	version, err := runGitCommand("describe", "--tags", "--abbrev=0")
	if err != nil {
		return "1.0.0", err
	}

	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	if pos := strings.Index(version, "-"); pos != -1 {
		version = version[:pos]
	}

	if _, err := semver.NewVersion(version); err != nil {
		return "", fmt.Errorf("%w: %s", errInvalidVersion, err.Error())
	}

	return version, nil
}

func hasBreakingChange() (bool, error) {
	commitMessage, err := runGitCommand("log", "--format=%B", "-n", "1", "HEAD")
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(commitMessage, breakingPrefix), nil
}

func getNextVersion(currentVersion string, branchName string, hasBreaking string) (string, error) {
	v, err := semver.NewVersion(currentVersion)
	if err != nil {
		return "", err
	}
	switch branchName {
	case developBranch:
		nextVersion := v.IncPatch()
		return nextVersion.String() + "-stage", nil
	case mainBranch:

		if hasBreaking == "true" {
			newVer := v.IncMajor()
			return newVer.String(), nil
		} else {
			newVer := v.IncMinor()
			return newVer.String(), nil
		}
	default:
		return "", fmt.Errorf("unsupported branch: %s", branchName)
	}
}

func createTag(newVersion string) error {
	tagCmd := exec.Command("git", "tag", "-a", newVersion, "-m", fmt.Sprintf("Release version %s", newVersion))
	err := tagCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}

	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: go run version.go <branch-name> <latest-release-version> <has-breaking-change>")
		os.Exit(1)
	}

	branchName := os.Args[1]
	latestReleaseVersion, err := cleanVersion(os.Args[2])
	hasBreakingChange := os.Args[3]

	// Verify that latestReleaseVersion is a valid semantic version
	if _, err := semver.NewVersion(latestReleaseVersion); err != nil {
		fmt.Println("error: latest-release-version is not a valid semantic version")
		os.Exit(1)
	}

	// Verify that hasBreakingChange is a valid boolean value
	if _, err := strconv.ParseBool(hasBreakingChange); err != nil {
		fmt.Println("error: has-breaking-change must be a boolean value (true or false)")
		os.Exit(1)
	}

	nextVersion, err := getNextVersion(latestReleaseVersion, branchName, hasBreakingChange)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println(nextVersion)
}
