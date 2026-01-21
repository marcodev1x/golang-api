package mysql

import (
	"fmt"
	"shortner-url/internal/domain"
	"shortner-url/internal/helpers"
	"time"

	"gorm.io/gorm"
)

type UrlRepository struct {
	db *gorm.DB
}

func NewUrlRepository(db *gorm.DB) *UrlRepository {
	return &UrlRepository{db}
}

func (r *UrlRepository) FindUrlByHashedId(hashedId string) (*domain.Urls, error) {
	var url domain.Urls

	err := r.db.
		Where("hashed_domain = ?", hashedId).
		Find(&url).
		Error

	return &url, err
}

func (r *UrlRepository) CreateUrl(url string, hashedDomain string, expiresAt *time.Time) (*domain.Urls, error) {
	var createdUrl domain.Urls

	err := r.db.
		Transaction(func(tx *gorm.DB) error {
			var obj domain.Urls

			if err := tx.Create(&domain.Urls{
				ShortenedUrl: url,
				ExpiresAt:    expiresAt,
			}).Scan(&obj).
				Error; err != nil {
				return err
			}

			finalDomain := helpers.GenerateHash(obj.Id) + hashedDomain

			if err := tx.Model(&domain.Urls{}).
				Where("id = ?", obj.Id).
				Update("hashed_domain", finalDomain).
				Scan(&createdUrl).
				Error; err != nil {
				return err
			}

			return nil
		})

	fmt.Println(createdUrl)

	return &createdUrl, err
}
