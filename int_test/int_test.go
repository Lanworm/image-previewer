package int_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/Lanworm/image-previewer/internal/http/client"
	"github.com/stretchr/testify/assert"
)

func GetImage(imagePath string) (*http.Response, error) {
	baseURL := "http://localhost:8090/fill/200/300/previewer-web-nginx:3080/images/"

	HTTPClient := client.NewHTTPClient(5 * time.Second)

	resp, err := HTTPClient.DoRequest("GET", baseURL+imagePath, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Получение изображения.
func TestImageExistence(t *testing.T) {
	imagePath := "image1.jpg"
	resp, err := GetImage(imagePath)

	assert.NoError(t, err, "Unexpected error occurred")
	assert.NotNil(t, resp, "Response should not be nil")
	resp.Body.Close()
}

// Дополнительные тесты для других сценариев...
