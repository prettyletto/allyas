package storage

func Load(path string) (Store, error) {
	return Store{SchemaVersion: 1, Aliases: []Alias{}}, nil
}

func Save(path string, store Store) error {
	return nil
}
