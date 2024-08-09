package dip

type Logic interface {
	SayHello(userId string) (string, error)
}

func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{
		Logger:    l,
		DataStore: ds,
	}
}
