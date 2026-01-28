package usecases

import (
	"encoding/json"

	"shortner-url/internal"
	"shortner-url/internal/domain"
	"shortner-url/internal/helpers"
	"time"
)

type UrlRepository interface {
	UpdateUrl(url *domain.Urls) error
	FindUrlByHashedId(hashedId, ref string) (*domain.Urls, error)
	CreateUrl(url, hashedDomain string, expiresAt *time.Time, ref string) (*domain.Urls, error)
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type UrlUsecase struct {
	repository UrlRepository
	cache      Cache
}

func NewUrlUseCase(repo UrlRepository, cache Cache) *UrlUsecase {
	return &UrlUsecase{
		repository: repo,
		cache:      cache,
	}
}

func (u *UrlUsecase) UpdateUrl(url *domain.Urls) error {
	err := u.repository.UpdateUrl(url)

	if err != nil {
		return internal.NewAPIError("Error updating url", 500, 100)
	}

	return nil
}

func (u *UrlUsecase) FindUrlByHashedId(hashedId string, ref string) (*domain.Urls, error) {
	cached, err := u.cache.Get(hashedId + ref)

	if err == nil {
		var url domain.Urls

		if err := json.Unmarshal([]byte(cached), &url); err != nil {
			return nil, err
		}

		u.UpdateUrl(&domain.Urls{
			Id:            url.Id,
			CountedClicks: url.CountedClicks + 1,
		})

		return &url, nil
	}

	url, err := u.repository.FindUrlByHashedId(hashedId, ref)

	if err != nil {
		return nil, internal.NewAPIError("Error trying to find url", 500, 480)
	}

	if url == nil {
		return nil, internal.NewAPIError("Url not found", 404, 100)
	}

	u.UpdateUrl(&domain.Urls{
		Id:            url.Id,
		CountedClicks: url.CountedClicks + 1,
	})

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

func (u *UrlUsecase) CreateUrl(url string, expiresAt *time.Time, ref string) (bool, error) {
	hashedDomain := helpers.GenerateHash(0, url)

	create, err := u.repository.CreateUrl(url, hashedDomain, expiresAt, ref)

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
