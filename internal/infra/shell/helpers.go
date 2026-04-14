package shell

import "os"

func RunningFromWrapper() bool {
	return os.Getenv("ALLYAS_SHELL_WRAPPER") == "1"
}
