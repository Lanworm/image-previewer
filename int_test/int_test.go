package int_test

import (
	"flag"
	"fmt"
	"net/http"
	"testing"
)

var Port = flag.String("port", "8080", "Container port")

// Функция для проверки наличия изображения по URL.
func checkImagePresence(imagePath string) bool {
	fmt.Println(Port)
	baseURL := "http://localhost:" + *Port + "/images/"
	response, err := http.Get(baseURL + imagePath)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK
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
