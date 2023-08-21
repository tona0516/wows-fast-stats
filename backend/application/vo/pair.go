package vo

type Pair struct {
	Key   string
	Value string
}

func NewPair(key string, value string) Pair {
	return Pair{
		Key:   key,
		Value: value,
	}
}
