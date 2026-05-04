package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prettyletto/allyas/internal/domain/models"
	"github.com/prettyletto/allyas/internal/infra/git"
	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

const (
	ProviderGit   = "git"
	DefaultBranch = "main"
	RepoDirName   = "sync"
)

type Mode string

const (
	ModeDefault   Mode = "default"
	ModeSetup     Mode = "setup"
	ModePreview   Mode = "preview"
	ModePull      Mode = "pull"
	ModePush      Mode = "push"
	ModeForcePull Mode = "force-pull"
	ModeForcePush Mode = "force-push"
)

type Input struct {
	ConfigPath string
	StorePath  string
	SourcePath string
	Mode       Mode
	Remote     string
	Branch     string
}

type Output struct {
	Message       string
	RepoDir       string
	Remote        string
	Branch        string
	LocalChanges  string
	RemoteChanges string
}

func Run(in Input) (Output, error) {
	if in.ConfigPath == "" || in.StorePath == "" || in.SourcePath == "" {
		return Output{}, fmt.Errorf("sync requires config, store, and source paths")
	}
	if in.Mode == "" {
		in.Mode = ModeDefault
	}

	cfg, err := storage.LoadConfig(in.ConfigPath)
	if err != nil {
		return Output{}, fmt.Errorf("load config: %w", err)
	}

	repoDir := SyncRepoDir(in.ConfigPath)
	branch := resolveBranch(in.Branch, cfg.Sync.Branch)
	remote := strings.TrimSpace(in.Remote)
	if remote == "" {
		remote = strings.TrimSpace(cfg.Sync.Remote)
	}

	switch in.Mode {
	case ModeSetup:
		if remote == "" {
			return Output{}, fmt.Errorf("sync setup requires a git remote")
		}
		return setup(in, cfg, repoDir, remote, branch)
	case ModePreview:
		return preview(in, cfg, repoDir, remote, branch)
	case ModeDefault:
		if !syncConfigured(cfg) {
			return Output{}, fmt.Errorf("sync is not configured; run allyas sync setup <git-remote>")
		}
		return syncDefault(in, repoDir, cfg.Sync.Remote, resolveBranch("", cfg.Sync.Branch))
	case ModePull:
		if !syncConfigured(cfg) {
			return Output{}, fmt.Errorf("sync is not configured; run allyas sync setup <git-remote>")
		}
		return pull(in, repoDir, cfg.Sync.Remote, resolveBranch("", cfg.Sync.Branch), false)
	case ModePush:
		if !syncConfigured(cfg) {
			return Output{}, fmt.Errorf("sync is not configured; run allyas sync setup <git-remote>")
		}
		return push(in, repoDir, cfg.Sync.Remote, resolveBranch("", cfg.Sync.Branch), false)
	case ModeForcePull:
		if !syncConfigured(cfg) {
			return Output{}, fmt.Errorf("sync is not configured; run allyas sync setup <git-remote>")
		}
		return pull(in, repoDir, cfg.Sync.Remote, resolveBranch("", cfg.Sync.Branch), true)
	case ModeForcePush:
		if !syncConfigured(cfg) {
			return Output{}, fmt.Errorf("sync is not configured; run allyas sync setup <git-remote>")
		}
		return push(in, repoDir, cfg.Sync.Remote, resolveBranch("", cfg.Sync.Branch), true)
	default:
		return Output{}, fmt.Errorf("unknown sync mode %q", in.Mode)
	}
}

func SyncRepoDir(configPath string) string {
	return filepath.Join(filepath.Dir(configPath), RepoDirName)
}

func setup(in Input, cfg models.Config, repoDir, remote, branch string) (Output, error) {
	if err := ensureRepo(repoDir, remote, branch); err != nil {
		return Output{}, err
	}
	client := git.NewClient()
	if err := client.Fetch(repoDir, "origin", branch); err == nil {
		_ = client.CheckoutRemoteBranch(repoDir, branch)
	}
	if repoHasCanonical(repoDir) {
		if err := applyRepoToCanonical(in.ConfigPath, in.StorePath, in.SourcePath, repoDir, true); err != nil {
			return Output{}, err
		}
		cfg, err := storage.LoadConfig(in.ConfigPath)
		if err != nil {
			return Output{}, fmt.Errorf("load synced config: %w", err)
		}
		cfg.Sync = models.SyncConfig{
			Enabled:  true,
			Provider: ProviderGit,
			Remote:   remote,
			Branch:   branch,
		}
		if err := storage.SaveConfig(in.ConfigPath, cfg); err != nil {
			return Output{}, fmt.Errorf("save config: %w", err)
		}
		if err := regenerateSource(in.ConfigPath, in.StorePath, in.SourcePath); err != nil {
			return Output{}, err
		}
		if err := copyCanonicalToRepo(in.ConfigPath, in.StorePath, repoDir); err != nil {
			return Output{}, err
		}
		if err := commitIfChanged(repoDir, "sync allyas config"); err != nil {
			return Output{}, err
		}
		if err := client.Push(repoDir, "origin", branch); err != nil {
			return Output{}, fmt.Errorf("push sync repo: %w", err)
		}
		return Output{Message: "Sync configured from remote", RepoDir: repoDir, Remote: remote, Branch: branch}, nil
	}

	cfg.Sync = models.SyncConfig{
		Enabled:  true,
		Provider: ProviderGit,
		Remote:   remote,
		Branch:   branch,
	}
	if err := storage.SaveConfig(in.ConfigPath, cfg); err != nil {
		return Output{}, fmt.Errorf("save config: %w", err)
	}
	if err := copyCanonicalToRepo(in.ConfigPath, in.StorePath, repoDir); err != nil {
		return Output{}, err
	}
	if err := commitIfChanged(repoDir, "sync allyas config"); err != nil {
		return Output{}, err
	}
	if err := client.Push(repoDir, "origin", branch); err != nil {
		return Output{}, fmt.Errorf("push sync repo: %w", err)
	}
	return Output{Message: "Sync configured", RepoDir: repoDir, Remote: remote, Branch: branch}, nil
}

func preview(in Input, cfg models.Config, repoDir, remote, branch string) (Output, error) {
	if remote == "" {
		remote = cfg.Sync.Remote
	}
	if branch == "" {
		branch = resolveBranch("", cfg.Sync.Branch)
	}
	status := "not configured"
	if syncConfigured(cfg) {
		status = "configured"
	}
	if !git.NewClient().IsRepo(repoDir) {
		return Output{Message: fmt.Sprintf("Sync %s; repo missing at %s", status, repoDir), RepoDir: repoDir, Remote: remote, Branch: branch}, nil
	}
	return Output{
		Message:       fmt.Sprintf("Sync %s; repo ready at %s", status, repoDir),
		RepoDir:       repoDir,
		Remote:        remote,
		Branch:        branch,
		LocalChanges:  previewLocalChanges(in, repoDir),
		RemoteChanges: previewRemoteChanges(repoDir, remote, branch),
	}, nil
}

func syncDefault(in Input, repoDir, remote, branch string) (Output, error) {
	if err := ensureRepo(repoDir, remote, branch); err != nil {
		return Output{}, err
	}
	if err := copyCanonicalToRepo(in.ConfigPath, in.StorePath, repoDir); err != nil {
		return Output{}, err
	}
	if err := commitIfChanged(repoDir, "sync allyas config"); err != nil {
		return Output{}, err
	}
	if err := git.NewClient().PullRebase(repoDir, "origin", branch); err != nil {
		return Output{}, fmt.Errorf("pull sync repo: %w", err)
	}
	if err := applyRepoToCanonical(in.ConfigPath, in.StorePath, in.SourcePath, repoDir, false); err != nil {
		return Output{}, err
	}
	if err := copyCanonicalToRepo(in.ConfigPath, in.StorePath, repoDir); err != nil {
		return Output{}, err
	}
	if err := commitIfChanged(repoDir, "sync allyas config"); err != nil {
		return Output{}, err
	}
	if err := git.NewClient().Push(repoDir, "origin", branch); err != nil {
		return Output{}, fmt.Errorf("push sync repo: %w", err)
	}
	return Output{Message: "Sync complete", RepoDir: repoDir, Remote: remote, Branch: branch}, nil
}

func pull(in Input, repoDir, remote, branch string, force bool) (Output, error) {
	if err := ensureRepo(repoDir, remote, branch); err != nil {
		return Output{}, err
	}
	client := git.NewClient()
	if force {
		if err := client.Fetch(repoDir, "origin", branch); err != nil {
			return Output{}, fmt.Errorf("fetch sync repo: %w", err)
		}
		if err := client.CheckoutRemoteBranch(repoDir, branch); err != nil {
			return Output{}, fmt.Errorf("checkout remote sync branch: %w", err)
		}
	} else if err := client.PullRebase(repoDir, "origin", branch); err != nil {
		return Output{}, fmt.Errorf("pull sync repo: %w", err)
	}

	if err := applyRepoToCanonical(in.ConfigPath, in.StorePath, in.SourcePath, repoDir, force); err != nil {
		return Output{}, err
	}
	return Output{Message: "Sync pulled", RepoDir: repoDir, Remote: remote, Branch: branch}, nil
}

func push(in Input, repoDir, remote, branch string, force bool) (Output, error) {
	if err := ensureRepo(repoDir, remote, branch); err != nil {
		return Output{}, err
	}
	if err := copyCanonicalToRepo(in.ConfigPath, in.StorePath, repoDir); err != nil {
		return Output{}, err
	}
	if err := commitIfChanged(repoDir, "sync allyas config"); err != nil {
		return Output{}, err
	}
	client := git.NewClient()
	if force {
		if err := client.ForcePush(repoDir, "origin", branch); err != nil {
			return Output{}, fmt.Errorf("force push sync repo: %w", err)
		}
	} else if err := client.Push(repoDir, "origin", branch); err != nil {
		return Output{}, fmt.Errorf("push sync repo: %w", err)
	}
	return Output{Message: "Sync pushed", RepoDir: repoDir, Remote: remote, Branch: branch}, nil
}

func ensureRepo(repoDir, remote, branch string) error {
	client := git.NewClient()
	if client.IsRepo(repoDir) {
		if err := client.RemoteSetOrigin(repoDir, remote); err != nil {
			return fmt.Errorf("set sync remote: %w", err)
		}
		if err := client.CheckoutBranch(repoDir, branch); err != nil {
			return fmt.Errorf("checkout sync branch: %w", err)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(repoDir), storage.DirPerm); err != nil {
		return err
	}
	if err := client.Clone(remote, repoDir); err != nil {
		return fmt.Errorf("clone sync repo: %w", err)
	}
	if err := client.CheckoutBranch(repoDir, branch); err != nil {
		return fmt.Errorf("checkout sync branch: %w", err)
	}
	return nil
}

func copyCanonicalToRepo(configPath, storePath, repoDir string) error {
	if err := copyFile(configPath, filepath.Join(repoDir, storage.ConfigFileName)); err != nil {
		return fmt.Errorf("copy config to sync repo: %w", err)
	}
	if err := copyFile(storePath, filepath.Join(repoDir, storage.StoreFileName)); err != nil {
		return fmt.Errorf("copy store to sync repo: %w", err)
	}
	return nil
}

func previewLocalChanges(in Input, repoDir string) string {
	configSame := sameFile(in.ConfigPath, filepath.Join(repoDir, storage.ConfigFileName))
	storeSame := sameFile(in.StorePath, filepath.Join(repoDir, storage.StoreFileName))
	if configSame && storeSame {
		return "no"
	}
	return "yes"
}

func previewRemoteChanges(repoDir, remote, branch string) string {
	if remote == "" {
		return "unknown"
	}
	client := git.NewClient()
	head, err := client.Head(repoDir)
	if err != nil {
		return "unknown"
	}
	remoteHead, err := client.LsRemote(remote, branch)
	if err != nil || remoteHead == "" {
		return "unknown"
	}
	if head == remoteHead {
		return "no"
	}
	return "yes"
}

func sameFile(left, right string) bool {
	leftData, err := os.ReadFile(left)
	if err != nil {
		return false
	}
	rightData, err := os.ReadFile(right)
	if err != nil {
		return false
	}
	return string(leftData) == string(rightData)
}

func applyRepoToCanonical(configPath, storePath, sourcePath, repoDir string, backup bool) error {
	if backup {
		if err := backupFile(configPath); err != nil {
			return fmt.Errorf("backup config: %w", err)
		}
		if err := backupFile(storePath); err != nil {
			return fmt.Errorf("backup store: %w", err)
		}
	}
	if err := copyFile(filepath.Join(repoDir, storage.ConfigFileName), configPath); err != nil {
		return fmt.Errorf("apply config from sync repo: %w", err)
	}
	if err := copyFile(filepath.Join(repoDir, storage.StoreFileName), storePath); err != nil {
		return fmt.Errorf("apply store from sync repo: %w", err)
	}

	cfg, err := storage.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load synced config: %w", err)
	}
	store, err := storage.LoadStore(storePath)
	if err != nil {
		return fmt.Errorf("load synced store: %w", err)
	}
	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(sourcePath, source); err != nil {
		return fmt.Errorf("save source: %w", err)
	}
	return nil
}

func regenerateSource(configPath, storePath, sourcePath string) error {
	cfg, err := storage.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	store, err := storage.LoadStore(storePath)
	if err != nil {
		return fmt.Errorf("load store: %w", err)
	}
	source := shell.RenderSource(store, cfg.DefaultGroup, cfg.Shell, cfg.AliasMode)
	if err := storage.SaveSource(sourcePath, source); err != nil {
		return fmt.Errorf("save source: %w", err)
	}
	return nil
}

func commitIfChanged(repoDir, message string) error {
	client := git.NewClient()
	if err := client.Add(repoDir, storage.ConfigFileName, storage.StoreFileName); err != nil {
		return fmt.Errorf("stage sync files: %w", err)
	}
	changed, err := client.HasStagedChanges(repoDir)
	if err != nil {
		return fmt.Errorf("check sync changes: %w", err)
	}
	if !changed {
		return nil
	}
	if err := client.Commit(repoDir, message); err != nil {
		return fmt.Errorf("commit sync changes: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), storage.DirPerm); err != nil {
		return err
	}
	return os.WriteFile(dst, data, storage.WritePerm)
}

func repoHasCanonical(repoDir string) bool {
	if exists, err := storage.FileExists(filepath.Join(repoDir, storage.ConfigFileName)); err != nil || !exists {
		return false
	}
	if exists, err := storage.FileExists(filepath.Join(repoDir, storage.StoreFileName)); err != nil || !exists {
		return false
	}
	return true
}

func backupFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	backupPath := fmt.Sprintf("%s.%d.bak", path, time.Now().UnixNano())
	return os.WriteFile(backupPath, data, storage.WritePerm)
}

func syncConfigured(cfg models.Config) bool {
	return cfg.Sync.Enabled && strings.TrimSpace(cfg.Sync.Provider) == ProviderGit && strings.TrimSpace(cfg.Sync.Remote) != ""
}

func resolveBranch(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return DefaultBranch
}
