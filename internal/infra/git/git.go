package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func (c Client) Init(dir string) error {
	_, err := run(dir, "init")
	return err
}

func (c Client) RemoteAddOrigin(dir, remote string) error {
	_, err := run(dir, "remote", "add", "origin", remote)
	return err
}

func (c Client) RemoteSetOrigin(dir, remote string) error {
	if _, err := run(dir, "remote", "set-url", "origin", remote); err == nil {
		return nil
	}
	return c.RemoteAddOrigin(dir, remote)
}

func (c Client) Add(dir string, files ...string) error {
	args := append([]string{"add"}, files...)
	_, err := run(dir, args...)
	return err
}

func (c Client) PullRebase(dir, remote, branch string) error {
	_, err := run(dir, "pull", "--rebase", remote, branch)
	return err
}

func (c Client) Commit(dir, message string) error {
	_, err := run(dir, "commit", "-m", message)
	return err
}

func (c Client) Push(dir, remote, branch string) error {
	_, err := run(dir, "push", "-u", remote, branch)
	return err
}

func (c Client) ForcePush(dir, remote, branch string) error {
	_, err := run(dir, "push", "--force-with-lease", "-u", remote, branch)
	return err
}

func (c Client) Fetch(dir, remote, branch string) error {
	_, err := run(dir, "fetch", remote, branch)
	return err
}

func (c Client) CheckoutBranch(dir, branch string) error {
	_, err := run(dir, "checkout", "-B", branch)
	return err
}

func (c Client) CheckoutRemoteBranch(dir, branch string) error {
	_, err := run(dir, "checkout", "-B", branch, "origin/"+branch)
	return err
}

func (c Client) Head(dir string) (string, error) {
	return run(dir, "rev-parse", "HEAD")
}

func (c Client) LsRemote(remote, branch string) (string, error) {
	out, err := run("", "ls-remote", remote, branch)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", nil
	}
	return fields[0], nil
}

func (c Client) HasStagedChanges(dir string) (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = dir
	err := cmd.Run()
	if err == nil {
		return false, nil
	}

	if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
		return true, nil
	}

	return false, err
}

func (c Client) IsRepo(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && info.IsDir()
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Allyas",
		"GIT_AUTHOR_EMAIL=allyas@example.invalid",
		"GIT_COMMITTER_NAME=Allyas",
		"GIT_COMMITTER_EMAIL=allyas@example.invalid",
	)

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
