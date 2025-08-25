package memoize

type store struct {
	records map[string]any
}

type IStore interface {
	memoize(key string, fn callback) (any, error, bool)
}

func Store() IStore {
	return &store{records: make(map[string]any)}
}

type callback func() (any, error)

func (s *store) memoize(key string, fn callback) (any, error, bool) {
	if value, exists := s.records[key]; exists {
		return value, nil, true
	}
	result, err := fn()
	if err != nil {
		return nil, err, false
	}
	s.records[key] = result
	return result, nil, false
}

func Memoize[Output any](key string, m IStore, f func() Output) Output {
	result, err, _ := m.memoize(key, func() (any, error) {
		data := f()
		return data, nil
	})
	if err != nil {
		panic(err)
	}
	return result.(Output)
}

func Memoize2[Output any](key string, m IStore, f func() (Output, error)) (Output, error) {
	result, err, _ := m.memoize(key, func() (any, error) {
		data, err := f()
		return data, err
	})
	return result.(Output), err
}
