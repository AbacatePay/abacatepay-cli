package auth

type MemoryStore struct {
	tokens        map[string]string
	defaultKey    string
	activeProfile string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tokens:     make(map[string]string),
		defaultKey: "default",
	}
}

func (m *MemoryStore) Save(token string) error {
	return m.SaveNamed(m.defaultKey, token)
}

func (m *MemoryStore) SaveNamed(name, token string) error {
	m.tokens[name] = token
	return nil
}

func (m *MemoryStore) Get() (string, error) {
	return m.GetNamed(m.defaultKey)
}

func (m *MemoryStore) GetNamed(name string) (string, error) {
	return m.tokens[name], nil
}

func (m *MemoryStore) Delete() error {
	return m.DeleteNamed(m.defaultKey)
}

func (m *MemoryStore) DeleteNamed(name string) error {
	delete(m.tokens, name)
	return nil
}

func (m *MemoryStore) SetActiveProfile(name string) error {
	m.activeProfile = name
	return nil
}

func (m *MemoryStore) GetActiveProfile() (string, error) {
	return m.activeProfile, nil
}
