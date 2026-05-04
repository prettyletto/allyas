package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct{}

func NewClient() Client {
	return Client{}
}

func (c Client) Clone(remote, dir string) error {
	_, err := run("", "clone", remote, dir)
	return err
}

func (c Client) Init (dir string) error {
	_, err := run(dir, "init")
	return err
}

func (c Client)RemoteAddOrigin(dir, remote string) error {
	_, err := run(dir, "remote", "add", "origin", remote)
	return err
}

func (c Client)PullRebase(dir string, files ...string) error {
	args := append([]string{"add"},files...)
	_, err := run(dir, args...)
	return err
}

func (c Client)Commit (dir, message string) error {
	_, err := run(dir, "commit", "-m", message)
	return err
}

func (c Client)Push (dir  string) error {
	_, err := run(dir, "push", "-u","origin" , "main" )
	return err
}

  func (c Client) HasStagedChanges(dir string) (bool, error) {
        _, err := run(dir, "diff", "--cached", "--quiet")
        if err == nil {
                return false, nil
        }

        if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
                return true, nil
        }

        return false, err
  }

func run(dir string, args ...string) (string, error) {
	cmd:= exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir 
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer


	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}
		if msg == "" {
			msg = err.Error()
		}

		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}
