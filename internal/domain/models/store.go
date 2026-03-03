package models

const SchemaVersion = 1

type Store struct {
	SchemaVersion int     `json:"schema_version"`
	Aliases       []Alias `json:"aliases"`
}

func DefaultStore() Store {
	return Store{
		SchemaVersion: SchemaVersion,
		Aliases:       []Alias{},
	}

}
