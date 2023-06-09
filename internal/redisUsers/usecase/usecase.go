package usecase

import (
	"AuthService/internal/redisUsers"
	"AuthService/pkg/secure"
	"fmt"
)

type RedisUsersUseCase struct {
	repo   redisUsers.Repository
	shield *secure.Shield
}

func RedisUseCase(repo redisUsers.Repository, shield *secure.Shield) redisUsers.UseCase {
	return &RedisUsersUseCase{
		repo:   repo,
		shield: shield,
	}
}

func (r *RedisUsersUseCase) UpdateAccess(refreshToken string, accessToken string, ua string, uid string) (string, error) {
	ok, err := r.ValidateRefresh(refreshToken, ua, uid)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	//RU, err := r.SetNewFingerprint(accessToken, refreshToken, fmt.Sprintf("%s%d", ua, id), id)
	//if err != nil {
	//	return "", err
	//}
	//fmt.Println(RU.AccessToken)
	use, err := r.repo.SetUser(accessToken, refreshToken, ua, uid)
	fmt.Println(accessToken == use.AccessToken)
	return accessToken, nil
}

func (r *RedisUsersUseCase) CheckAccessTokenByFingerprint(fingerprint string, uid string) bool {
	return r.repo.CheckAccessTokenByFingerprint(fingerprint, uid)
}

func (r *RedisUsersUseCase) SetNewFingerprint(accessToken string, refreshToken string, fingerprint string,
	uid string) (*redisUsers.RedisUser, error) {
	return r.repo.SetNewFingerprint(accessToken, refreshToken, fingerprint, uid)
}

func (r *RedisUsersUseCase) GetUIDByFingerprint(fingerprint string) (string, error) {
	return r.repo.GetUIDByFingerprint(fingerprint)
}

func (r *RedisUsersUseCase) UpdateToken(accessToken string, refreshToken string, fingerprint string) (*redisUsers.RedisUser, error) {
	uid, err := r.GetUIDByFingerprint(fingerprint)
	if err != nil {
		return nil, err
	}
	return r.repo.SetUser(accessToken, refreshToken, fingerprint, uid)
}

func (r *RedisUsersUseCase) Validate(accessToken string, fingerprint string, uid string) (bool, error) {
	ok, err := r.repo.CheckAccessToken(accessToken, fingerprint, uid)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (r *RedisUsersUseCase) ValidateRefresh(refreshToken string, fingerprint string, uid string) (bool, error) {
	ok, err := r.repo.CheckRefreshToken(refreshToken, fingerprint, uid)
	fmt.Println(ok, err)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (r *RedisUsersUseCase) RemoveUser(params *redisUsers.RemoveUser) (int64, error) {
	return r.repo.RemoveUser(params)
}
