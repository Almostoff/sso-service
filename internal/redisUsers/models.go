package redisUsers

type RedisUser struct {
	AccessToken  string
	RefreshToken string
}

type RemoveUser struct {
	Fingerprint string
	Uuid        string
}
