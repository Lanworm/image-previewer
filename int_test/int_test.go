package int_test

import (
	"flag"
	"net/http"
	"testing"
	"time"

	"github.com/Lanworm/image-previewer/internal/http/client"
)

var Port = flag.String("port", "8080", "Container port")

func checkImagePresence(imagePath string) bool {
	baseURL := "http://localhost:" + *Port + "/temp/"

	HTTPClient := client.NewHTTPClient(100 * time.Second)

	resp, err := HTTPClient.DoRequest("GET", baseURL+imagePath, nil, nil)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Тестирование различных сценариев.
func TestImageExistence(t *testing.T) {
	if imagePath := "image1.jpg"; !checkImagePresence(imagePath) {
		t.Errorf("Expected image %s to exist, but it was not found", imagePath)
	}
}

func TestImageNotFound(t *testing.T) {
	imagePath := "non_existent_image.jpg"
	if checkImagePresence(imagePath) {
		t.Errorf("Expected image %s to not exist, but it was found", imagePath)
	}
}

func TestNonImageFile(t *testing.T) {
	imagePath := "not_an_image.exe"
	if checkImagePresence(imagePath) {
		t.Errorf("Expected %s to not be an image, but it was found", imagePath)
	}
}

// Дополнительные тесты для других сценариев...
