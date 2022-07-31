package go_memcache

type Item struct {
	Key        string
	Value      string
	Expiration int32
	CreatedAt  int64
}

func (i *Item) IsExpired() bool {
	return isExpired(i.CreatedAt, i.Expiration)
}
