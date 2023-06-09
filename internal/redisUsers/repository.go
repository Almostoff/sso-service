package redisUsers

type Repository interface {
	SetUser(accessToken string, refreshToken string, fingerprint string, uid string) (*RedisUser, error)
	CheckAccessTokenByFingerprint(fingerprint string, uid string) bool
	GetUIDByFingerprint(fingerprint string) (string, error)
	CheckAccessToken(accessToken string, fingerprint string, uid string) (bool, error)
	CheckRefreshToken(refreshToken string, fingerprint string, uid string) (bool, error)
	GetUser(fingerprint string, uid string) (*RedisUser, error)
	SetNewFingerprint(accessToken string, refreshToken string, fingerprint string, uid string) (*RedisUser, error)
	GetFingerprintsCount(uid string) int64
	RemoveUser(params *RemoveUser) (int64, error)
}
