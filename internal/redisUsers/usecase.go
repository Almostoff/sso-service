package redisUsers

type UseCase interface {
	CheckAccessTokenByFingerprint(fingerprint string, uid string) bool
	SetNewFingerprint(accessToken string, refreshToken string, fingerprint string, uid string) (*RedisUser, error)
	GetUIDByFingerprint(fingerprint string) (string, error)
	UpdateToken(accessToken string, refreshToken string, fingerprint string) (*RedisUser, error)
	Validate(accessToken string, fingerprint string, uid string) (bool, error)
	ValidateRefresh(refreshToken string, fingerprint string, uid string) (bool, error)
	RemoveUser(params *RemoveUser) (int64, error)
	UpdateAccess(refreshToken string, accessToken string, email string, uid string) (string, error)
}
