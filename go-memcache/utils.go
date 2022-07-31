package go_memcache

import "time"

func nowTimeStamp() int64 {
	return time.Now().Unix()
}

func isExpired(createdAt int64, expiration int32) bool {
	if expiration <= 0 {
		return false
	}
	return createdAt+int64(expiration) < nowTimeStamp()
}
