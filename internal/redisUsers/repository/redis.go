package repository

import (
	"AuthService/internal/redisUsers"
	"AuthService/pkg/secure"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	client *redis.Client
	shield *secure.Shield
}

func NewRedisRepository(client *redis.Client, shield *secure.Shield) redisUsers.Repository {
	return &redisRepository{
		client: client,
		shield: shield,
	}
}

var ctx = context.Background()

func (r redisRepository) CheckAccessTokenByFingerprint(fingerprint string, uid string) bool {
	data, err := r.client.Get(ctx, fmt.Sprintf("%s:%s:access", fingerprint, uid)).Result()
	fmt.Println(data, err)
	if err == redis.Nil {
		return false
	} else {
		return true
	}
}

func (r redisRepository) GetUIDByFingerprint(fingerprint string) (string, error) {
	fingerprintKey, _ := r.client.Keys(ctx, fmt.Sprintf("%s:*:uid", fingerprint)).Result()
	if len(fingerprintKey) == 0 {
		return "", nil
	}
	fingerprintValue, err := r.client.Get(ctx, fingerprintKey[0]).Result()
	if err == redis.Nil {
		return "", errors.New("no uid by fingerprint")
	}
	return fingerprintValue, nil
}

func (r redisRepository) GetFingerprintsCount(uid string) int64 {
	keys, _ := r.client.Keys(ctx, fmt.Sprintf("*:%s:uid", uid)).Result()
	return int64(len(keys))
}

func (r redisRepository) SetUser(accessToken string, refreshToken string, fingerprint string, uid string) (*redisUsers.RedisUser, error) {
	baseFingerprintUidKey := fmt.Sprintf("%s:%s", fingerprint, uid)
	formattedAccessKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "access")
	formattedRefreshKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "refresh")
	formattedUIDKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "uid")

	r.client.Set(ctx, formattedAccessKey, accessToken, time.Hour*3)
	r.client.Set(ctx, formattedRefreshKey, refreshToken, time.Hour*72)
	r.client.Set(ctx, formattedUIDKey, uid, time.Hour*24*360)

	return &redisUsers.RedisUser{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (r redisRepository) CheckAccessToken(accessToken string, fingerprint string, uid string) (bool, error) {
	accessTokenKey := fmt.Sprintf("%s:%s:%s", fingerprint, uid, "access")
	redisAccessToken, err := r.client.Get(ctx, accessTokenKey).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if redisAccessToken == accessToken {
		return true, nil
	} else {
		return false, nil
	}
}

func (r redisRepository) CheckRefreshToken(refreshToken string, fingerprint string, uid string) (bool, error) {
	refreshTokenKey := fmt.Sprintf("%s:%s:%s", fingerprint, uid, "refresh")
	redisRefreshToken, err := r.client.Get(ctx, refreshTokenKey).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if redisRefreshToken == refreshToken {
		return true, nil
	} else {
		return false, nil
	}
}

func (r redisRepository) GetUser(fingerprint string, uid string) (*redisUsers.RedisUser, error) {
	baseFingerprintUidKey := fmt.Sprintf("%s:%s", fingerprint, uid)
	formattedAccessKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "access")
	formattedRefreshKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "refresh")

	accessToken, _ := r.client.Get(ctx, formattedAccessKey).Result()
	refreshToken, _ := r.client.Get(ctx, formattedRefreshKey).Result()
	return &redisUsers.RedisUser{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r redisRepository) SetNewFingerprint(accessToken string, refreshToken string, fingerprint string, uid string) (*redisUsers.RedisUser, error) {
	if r.GetFingerprintsCount(uid) >= 8 {
		_, err := r.RemoveUser(nil)
		if err != nil {
			return nil, err
		}
		//return nil, errors.New("a lot of fingerprints")
	}
	user, err := r.SetUser(accessToken, refreshToken, fingerprint, uid)

	return user, err
}

func (r redisRepository) RemoveUser(params *redisUsers.RemoveUser) (int64, error) {
	baseFingerprintUidKey := fmt.Sprintf("%s:%s", params.Fingerprint, params.Uuid)
	formattedAccessKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "access")
	formattedRefreshKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "refresh")
	formattedUIDKey := fmt.Sprintf("%s:%s", baseFingerprintUidKey, "uid")
	result, err := r.client.Del(ctx, formattedAccessKey, formattedRefreshKey, formattedUIDKey).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
