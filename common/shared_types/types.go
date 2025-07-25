package shared_types

type RedisUserInfo struct {
	Email string `json:"email"`
	Token string `json:"token"`
	TTL   int64  `json:"ttl"`
}
