package bootstrap

import (
	"fmt"

	"github.com/prettyletto/allyas/internal/infra/shell"
	"github.com/prettyletto/allyas/internal/infra/storage"
)

type InstallInput struct {
	Shell    shell.Type
	Manual   bool
	HookPath string
}

type InstallOutput struct {
	RCPath      string
	SourceLine  string
	AlreadyDone bool
}

func RenderRCSourceBlock(hookPath string) string {
	return fmt.Sprintf(
		"%s\n"+
			"if [ -f %q ]; then\n"+
			"       . %q\n"+
			"fi\n"+
			"%s\n",
		storage.AllyasRCStart,
		hookPath,
		hookPath,
		storage.AllyasRCEnd,
	)
}

func PreviewInstall(in InstallInput) (InstallOutput, string, error) {
	rcPath, err := storage.RCPath(in.Shell)
	if err != nil {
		return InstallOutput{}, "", err
	}

	block := RenderRCSourceBlock(in.HookPath)
	content, err := storage.LoadOptionalFile(rcPath)
	if err != nil {
		return InstallOutput{}, "", err
	}

	return InstallOutput{
		RCPath:      rcPath,
		SourceLine:  fmt.Sprintf("source %q", in.HookPath),
		AlreadyDone: storage.HasRCBlock(content),
	}, block, nil
}

func RunInstall(in InstallInput) (InstallOutput, string, error) {
	out, block, err := PreviewInstall(in)
	if err != nil {
		return InstallOutput{}, "", err
	}

	if in.Manual {
		return out, block, nil
	}

	if err := storage.AppendRCBlock(out.RCPath, block); err != nil {
		return InstallOutput{}, "", err
	}
	out.AlreadyDone = false
	return out, block, nil
}
