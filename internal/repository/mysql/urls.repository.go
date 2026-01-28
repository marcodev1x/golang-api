package mysql

import (
	"shortner-url/internal/domain"
	"shortner-url/internal/usecases"
	"time"

	"gorm.io/gorm"
)

type UrlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) usecases.UrlRepository {
	return &UrlRepository{db}
}

func (r *UrlRepository) UpdateUrl(url *domain.Urls) error {
	return r.db.
		Model(&domain.Urls{}).
		Where("id = ?", url.Id).
		Updates(url).
		Error
}

func (r *UrlRepository) FindUrlByHashedId(hashedId string, ref string) (*domain.Urls, error) {
	var url domain.Urls

	err := r.db.
		Where("hashed_domain = ?", hashedId).
		Where("reference = ?", ref).
		Find(&url).
		Error

	return &url, err
}

func (r *UrlRepository) CreateUrl(url string, hashedDomain string, expiresAt *time.Time, ref string) (*domain.Urls, error) {
	var createdUrl domain.Urls

	err := r.db.
		Transaction(func(tx *gorm.DB) error {
			var obj domain.Urls

			if err := tx.Create(&domain.Urls{
				ShortenedUrl: url,
				ExpiresAt:    expiresAt,
				Reference:    ref,
			}).Scan(&obj).
				Error; err != nil {
				tx.Rollback()
				return err
			}

			if err := tx.Model(&domain.Urls{}).
				Where("id = ?", obj.Id).
				Update("hashed_domain", hashedDomain).
				Scan(&createdUrl).
				Error; err != nil {
				tx.Rollback()
				return err
			}

			return nil
		})

	return &createdUrl, err
}
