package dip

type DataStore interface {
	UserNameForId(userID string) (string, bool)
}

type SimpleDataStore struct {
	userData map[string]string
}

func (sds SimpleDataStore) UserNameForId(userID string) (string, bool) {
	name, ok := sds.userData[userID]
	return name, ok
}

func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[string]string{
			"1": "Freed",
			"2": "Mary",
			"3": "Pat",
		},
	}
}
