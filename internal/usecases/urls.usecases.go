package usecases

import (
	"encoding/json"
	"shortner-url/internal"
	"shortner-url/internal/domain"
	"shortner-url/internal/helpers"
	"shortner-url/internal/repository/mysql"
	"time"
)

type UrlUsecase struct {
	repository *mysql.UrlRepository
	cache      *RedisUsecase
}

func NewUrlUseCase(repository *mysql.UrlRepository, redis *RedisUsecase) *UrlUsecase {
	return &UrlUsecase{repository, redis}
}

func (u *UrlUsecase) FindUrlByHashedId(hashedId string) (*domain.Urls, error) {
	cached, err := u.cache.Get(hashedId)

	if err == nil {
		var url domain.Urls

		if err := json.Unmarshal([]byte(cached), &url); err != nil {
			return nil, err
		}

		return &url, nil
	}

	url, err := u.repository.FindUrlByHashedId(hashedId)

	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, internal.NewAPIError("Url not found", 404, 100)
	}

	bytes, err := json.Marshal(url)

	if err != nil {
		return nil, err
	}

	// Seta o cache por 24 horas
	if err := u.cache.Set(url.HashedDomain, string(bytes), 1440*time.Minute); err != nil {
		return nil, err
	}

	return url, nil
}

func (u *UrlUsecase) CreateUrl(url string, expiresAt *time.Time) (bool, error) {
	hashedDomain := helpers.GenerateHash(0)

	create, err := u.repository.CreateUrl(url, hashedDomain, expiresAt)

	if err != nil {
		return false, err
	}

	bytes, err := json.Marshal(create)

	if err != nil {
		return false, err
	}

	if err := u.cache.Set(create.HashedDomain, string(bytes), 1440*time.Minute); err != nil {
		return false, err
	}

	return true, nil
}
