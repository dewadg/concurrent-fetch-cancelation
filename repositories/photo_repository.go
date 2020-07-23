package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dewadg/concurrent-fetch-cancelation/models"
)

const totalPhotos = 5000

type PhotoRepository interface {
	Get(ctx context.Context, timeout time.Duration) ([]models.Photo, int, int, int, int, error)
}

type photoRepository struct {
	client   *http.Client
	photoIDs [totalPhotos]uint
}

func NewPhotoRepository() PhotoRepository {
	var photoIDs [totalPhotos]uint
	for i := 0; i < totalPhotos; i++ {
		photoIDs[i] = uint(i + 1)
	}

	return &photoRepository{
		client:   &http.Client{},
		photoIDs: photoIDs,
	}
}

func (repository *photoRepository) Get(ctx context.Context, timeout time.Duration) ([]models.Photo, int, int, int, int, error) {
	ctx, markProcessAsDone := context.WithTimeout(ctx, timeout)
	photoChan := make(chan models.Photo)
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	for _, photoID := range repository.photoIDs {
		wg.Add(1)

		go func(wg *sync.WaitGroup, photoID uint) {
			defer wg.Done()

			photo, err := repository.fetchPhotoByID(ctx, photoID)
			if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
				errChan <- err
			} else {
				photoChan <- photo
			}
		}(&wg, photoID)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()

		markProcessAsDone()
	}(&wg)

	var errCount int
	var fetchedPhotos []models.Photo

	for {
		select {
		case photo := <-photoChan:
			fetchedPhotos = append(fetchedPhotos, photo)
		case err := <-errChan:
			if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
				log.Println(err.Error())
				errCount++
			}
		case <-ctx.Done():
			return fetchedPhotos, totalPhotos, len(fetchedPhotos), errCount, totalPhotos - len(fetchedPhotos) - errCount, nil
		}
	}
}

func (repository *photoRepository) fetchPhotoByID(ctx context.Context, id uint) (models.Photo, error) {
	if err := ctx.Err(); err != nil {
		return models.Photo{}, err
	}

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", id)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.Photo{}, err
	}

	response, err := repository.client.Do(request)
	if err != nil {
		return models.Photo{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return models.Photo{}, errors.New("Request failed")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.Photo{}, err
	}

	var payload models.Photo
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return models.Photo{}, err
	}

	if err := ctx.Err(); err != nil {
		return models.Photo{}, err
	}
	return payload, nil
}
