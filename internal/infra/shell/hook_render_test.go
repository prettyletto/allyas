package shell

import (
	"strings"
	"testing"
)

func TestRenderHookReloadsAfterSync(t *testing.T) {
	out := RenderHook("/tmp/allyas/aliases.sh", "allyas")
	if !strings.Contains(out, "create|edit|remove|init|import|sync)") {
		t.Fatalf("hook does not reload after sync:\n%s", out)
	}
}
