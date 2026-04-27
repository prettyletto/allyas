package commands

// type modeFlags struct {
// 	Force     bool
// 	Shell     shell.Type
// 	AliasMode models.AliasMode
// }
//
// func parseAliasMode(s string) (models.AliasMode, error) {
// 	mode := models.AliasMode(strings.ToLower(strings.TrimSpace(s)))
//
// 	if !mode.Valid() {
// 		return "", fmt.Errorf("invalid alias mode %q, expected one of: plain, tracked", s)
// 	}
//
// 	return mode, nil
// }
//
// func parseShell(s string) (shell.Type, error) {
// 	t := shell.Type(s)
// 	if !t.Valid() {
// 		return "", fmt.Errorf("unsupported shell: %q", s)
// 	}
// 	return t, nil
// }
//
// func parseInitargs(args []string) (initFlags, error) {
// 	var f initFlags
//
// 	for i := 0; i < len(args); i++ {
// 		a := args[i]
//
// 		switch a {
// 		case "--shell":
// 			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
// 				return f, fmt.Errorf("%s requires value", a)
// 			}
// 			sh, err := parseShell(args[i+1])
// 			if err != nil {
// 				return f, err
// 			}
// 			f.Shell = sh
// 			i++
// 		case "--force", "-f":
// 			f.Force = true
// 		case "--alias-mode", "--mode":
// 			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
// 				return f, fmt.Errorf("%s requires value", a)
// 			}
// 			mode, err := parseAliasMode(args[i+1])
// 			if err != nil {
// 				return f, err
// 			}
// 			f.AliasMode = mode
// 			i++
//
// 		default:
// 			return f, fmt.Errorf("unkown arg: %s", a)
// 		}
//
// 	}
// 	return f, nil
// }

type ModeCommand struct{}

func NewModeCommand() *ModeCommand {
	return &ModeCommand{}
}

func (c *ModeCommand) Names() []string {
	return []string{"mode", "-m"}
}

func (c *ModeCommand) Usage() string {
	return "allyas mode"
}

func (c *ModeCommand) Description() string {
	return `Configure the files needed to the app run locally 
	and prepare the file to be injected in the shell`
}

func (c *ModeCommand) Execute(ctx CommandContext, args []string) error {

	return nil
}
