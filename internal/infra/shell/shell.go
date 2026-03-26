package shell

type Type string

const (
	Bash Type = "bash"
	Zsh  Type = "zsh"
	Sh   Type = "sh"
)

func (t Type) Valid() bool {
	switch t {
	case Bash, Zsh, Sh:
		return true
	default:
		return false
	}
}
