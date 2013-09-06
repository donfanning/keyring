package keyring

import (
	"fmt"
	"os/exec"
	"regexp"
)

type osxProvider struct {
}

var pwRe = regexp.MustCompile("password: \"(.+)\"")

func (p osxProvider) Get(Service, Username string) (string, error) {
	args := []string{"find-generic-password",
		"-s", Service,
		"-a", Username,
		"-g"}
	c := exec.Command("security", args...)
	o, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("security: %s", err)
	}
	matches := pwRe.FindStringSubmatch(string(o))
	if len(matches) != 2 {
		return "", fmt.Errorf("expected two submatches, got %d in: '%s'",
			len(matches), string(o))
	}
	return matches[1], nil
}

func (p osxProvider) Set(Service, Username, Password string) error {
	args := []string{"add-generic-password",
		"-s", Service,
		"-a", Username,
		"-w", Password,
		"-U"}
	c := exec.Command("security", args...)
	err := c.Run()
	if err != nil {
		o, _ := c.CombinedOutput()
		return fmt.Errorf(string(o))
	}
	return nil
}

func init() {
	registerProvider("osxkeychain", osxProvider{}, true)
}
