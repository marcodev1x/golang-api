package urls_tests

import (
	"encoding/json"
	"errors"
	"shortner-url/internal"
	"shortner-url/internal/domain"
	"shortner-url/internal/usecases"
	"shortner-url/internal/usecases/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUrlUseCases_UpdateUrl(t *testing.T) {
	t.Run("should update url successfully", func(t *testing.T) {
		controller := gomock.NewController(t) // Setup de controller
		defer controller.Finish()             // Finaliza o controller ao final do teste

		mockRepo := mocks.NewMockUrlRepository(controller)
		useCase := usecases.NewUrlUseCase(mockRepo, nil)

		url := &domain.Urls{
			Id:            1,
			CountedClicks: 10,
		}

		mockRepo.
			EXPECT().
			UpdateUrl(url).
			Return(nil).
			Times(1)

		err := useCase.UpdateUrl(url)

		assert.NoError(t, err)
	})

	t.Run("should return error if repository returns error", func(t *testing.T) {
		// AAA
		// ARRANGE
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := mocks.NewMockUrlRepository(controller)
		useCase := usecases.NewUrlUseCase(mockRepo, nil)

		url := &domain.Urls{
			Id:            1,
			CountedClicks: 10,
		}
		expectedError := internal.NewAPIError("Error updating url", 500, 100)

		mockRepo.
			EXPECT().
			UpdateUrl(url).
			Return(expectedError).
			Times(1)
		// ACT
		err := useCase.UpdateUrl(url)

		// ASSERT
		assert.Equal(t, expectedError, err)
	})
}

func TestUrlUseCases_FindUrlByHashedId(t *testing.T) {

	// common values
	hashedId := "hashedId"
	ref := "ref"
	mockedUrl := &domain.Urls{
		Id:            1,
		HashedDomain:  hashedId,
		ShortenedUrl:  "https://google.com",
		CountedClicks: 10,
	}

	t.Run("should return url successfully in cache (Redis)", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := mocks.NewMockUrlRepository(controller)
		mockCache := mocks.NewMockRedisUseCase(controller)

		useCase := usecases.NewUrlUseCase(mockRepo, mockCache)

		cachedBytes, _ := json.Marshal(mockedUrl)

		// Mock do cache
		mockCache.
			EXPECT().
			Get(hashedId+ref).
			Return(string(cachedBytes), nil).
			Times(1)

		// Mock do repository
		mockRepo.
			EXPECT().
			UpdateUrl(gomock.AssignableToTypeOf(&domain.Urls{})).
			DoAndReturn(func(u *domain.Urls) error {
				assert.Equal(t, mockedUrl.Id, u.Id)
				assert.Equal(t, mockedUrl.CountedClicks+1, u.CountedClicks)

				return nil
			}).
			Times(1)

		// Mock do repository (não deve ser chamado)
		mockRepo.
			EXPECT().
			FindUrlByHashedId(gomock.Any(), gomock.Any()).
			Times(0)

		result, err := useCase.FindUrlByHashedId(hashedId, ref)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockedUrl.Id, result.Id)
	})

	t.Run("should return error if json unmarshal fails in cache (Redis)", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := mocks.NewMockUrlRepository(controller)
		mockCache := mocks.NewMockRedisUseCase(controller)

		useCase := usecases.NewUrlUseCase(mockRepo, mockCache)

		mockedNotOkUrl := "not ok url"

		mockCache.
			EXPECT().
			Get(hashedId+ref).
			Return(mockedNotOkUrl, nil).
			Times(1)

		// Não chama pois deu erro antes
		mockRepo.
			EXPECT().
			UpdateUrl(gomock.AssignableToTypeOf(&domain.Urls{})).
			Times(0)

		result, err := useCase.FindUrlByHashedId(hashedId, ref)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should throw error if cache was not found and repository returns error", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := mocks.NewMockUrlRepository(controller)
		mockCache := mocks.NewMockRedisUseCase(controller)
		useCase := usecases.NewUrlUseCase(mockRepo, mockCache)

		expectedError := internal.NewAPIError("Error trying to find url", 500, 480)

		mockCache.
			EXPECT().
			Get(hashedId+ref).
			Return("", errors.New("not found")).
			Times(1)

		mockRepo.
			EXPECT().
			FindUrlByHashedId(hashedId, ref).
			Return(nil, errors.New("not found")).
			Times(1)

		// Não chama pois não encontrou
		mockRepo.
			EXPECT().
			UpdateUrl(gomock.AssignableToTypeOf(&domain.Urls{})).
			Times(0)

		result, err := useCase.FindUrlByHashedId(hashedId, ref)

		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

	t.Run("should throw error if url was not found in cache and repository", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := mocks.NewMockUrlRepository(controller)
		mockCache := mocks.NewMockRedisUseCase(controller)
		useCase := usecases.NewUrlUseCase(mockRepo, mockCache)

		expectedError := internal.NewAPIError("Url not found", 404, 100)

		mockCache.
			EXPECT().
			Get(hashedId+ref).
			Return("", errors.New("not found")).
			Times(1)

		mockRepo.
			EXPECT().
			FindUrlByHashedId(hashedId, ref).
			Return(nil, nil).
			Times(1)

		// Não chama pois não encontrou
		mockRepo.
			EXPECT().
			UpdateUrl(gomock.AssignableToTypeOf(&domain.Urls{})).
			Times(0)

		result, err := useCase.FindUrlByHashedId(hashedId, ref)

		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})

}
