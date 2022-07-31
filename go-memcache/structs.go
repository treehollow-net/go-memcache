package go_memcache

type Item struct {
	Key        string
	Value      string
	Expiration int32
	CreatedAt  int64
}
