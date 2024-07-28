package int_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Lanworm/image-previewer/internal/http/client"
	"github.com/Lanworm/image-previewer/internal/http/server/dto"
	"github.com/Lanworm/image-previewer/internal/service"
	"github.com/stretchr/testify/assert"
)

// Картинка найдена в кэше.
func TestImageFoundInCache(t *testing.T) {
	imagePath := "image1.jpg"

	// Запрашиваем картинку с сервера
	resp, err := GetImage(imagePath, "100", "200", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()

	// Проверяем, что картинка успешно получена
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Читаем содержимое полученного изображения
	expectedImageContents, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Удаляем файл из контейнера
	err = RemoveFile(imagePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Запрашиваем картинку с сервера повторно
	resp, err = GetImage(imagePath, "100", "200", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()

	// Проверяем, что картинка успешно получена (из кэша)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Читаем содержимое полученного изображения из кэша
	actualImageContents, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Сравниваем содержимое полученного изображения из кэша с эталонным
	assert.True(t, bytes.Equal(expectedImageContents, actualImageContents), "Полученное изображение из кэша не совпадает с эталонным.")
}

// Удаленный сервер не существует;.
func TestRemoteServerDoesNotExist(t *testing.T) {
	imagePath := "image1.jpg"

	// Запрашиваем картинку с сервера
	resp, err := GetImage(imagePath, "300", "400", "http://NotExistHost")

	// Проверяем, что ошибка не равна nil
	assert.NotNil(t, err, "expected an error but got nil")

	if err != nil {
		// Проверяем, что текст ошибки содержит "remote server does not exist"
		assert.Contains(t, err.Error(), service.ErrServerDoesNotExist.Error(), "unexpected error message")
		return
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	t.Fatal("expected error, but got none")
}

// Удаленный сервер существует, но изображение не найдено (404 Not Found);.
func TestImageNotFound(t *testing.T) {
	imagePath := "nonexistent_image.jpg"

	// Запрашиваем картинку с сервера, который существует
	resp, err := GetImage(imagePath, "400", "300", "")

	// Проверяем, что ошибка не равна nil
	assert.NotNil(t, err, "expected an error, but got none")

	if err != nil {
		// Проверяем, что текст ошибки содержит "image not found on remote server"
		assert.Contains(t, err.Error(), service.ErrImageNotFound.Error(), "unexpected error message")
		return
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	t.Fatal("expected error, but got none")
}

// Удаленный сервер существует, но изображение не изображение.
func TestFileIsNotAnImage(t *testing.T) {
	imagePath := "image.exe"

	// Запрашиваем файл с сервера, который существует
	resp, err := GetImage(imagePath, "400", "300", "")

	// Проверяем, что ошибка не равна nil
	assert.NotNil(t, err, "expected an error, but got none")

	if err != nil {
		// Проверяем, что текст ошибки содержит "target file is not an image"
		assert.Contains(t, err.Error(), service.ErrTargetNotImage.Error(), "unexpected error message")
		return
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	t.Fatal("expected error, but got none")
}

// Дополнительные тесты для других сценариев...

// RemoveFile удаляет указанный файл из директории int_test/test_files.
func RemoveFile(fileName string) error {
	filePath := fmt.Sprintf("./test_files/%s", fileName) // Полный путь к файлу

	// Удаляем файл
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ошибка при удалении файла %s: %w", filePath, err)
	}

	return nil // Возвращаем nil, если файл успешно удален
}

// IsRunningInContainer проверяет, выполняется ли приложение внутри контейнера.
func IsRunningInContainer() bool {
	// Проверяет наличие файла /proc/1/cgroup
	_, err := os.Stat("/proc/1/cgroup")
	return !os.IsNotExist(err)
}

// GetImage делает HTTP-запрос для получения изображения по указанному пути.
func GetImage(imgPath string, imgH string, imgW string, hostPath string) (*http.Response, error) {
	var nginxURL string
	var appURL string

	// Определяем URL сервиса в зависимости от того, запущен ли код в контейнере.
	if IsRunningInContainer() {
		nginxURL = "previewer-web-nginx" // Используем имя сервиса в Docker.
		appURL = "app"
	} else {
		nginxURL = "localhost" // Используем localhost для локальной разработки.
		appURL = "localhost"
	}

	// Подменяем URL сервиса в если это необходимо для теста.
	if hostPath != "" {
		nginxURL = hostPath
	}

	// Формируем базовый URL для запроса.
	baseURL := fmt.Sprintf("http://%s:8090/fill/%s/%s/%s:3080/images/%s", appURL, imgH, imgW, nginxURL, imgPath)

	// Создаем HTTP-клиент с таймаутом в 10 секунд.
	HTTPClient := client.NewHTTPClient(50 * time.Second)

	// Выполняем GET-запрос по сформированному URL.
	resp, err := HTTPClient.DoRequest("GET", baseURL, nil, nil)
	if err != nil {
		return nil, err // Возвращаем ошибку, если запрос не удался.
	}

	// Проверяем, не вернулся ли статус 500 (Internal Server Error).
	if resp.StatusCode == http.StatusInternalServerError {
		var res dto.Result
		decoder := json.NewDecoder(resp.Body) // Создаем декодер для чтения JSON из тела ответа.
		defer resp.Body.Close()               // Закрываем тело запроса после использования.

		// Декодируем ответ в структуру Result.
		err := decoder.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest) // Возвращаем ошибку, если декодирование не удалось.
		}
		return nil, errors.New(res.Message) // Возвращаем сообщение об ошибке из ответа.
	}

	// Возвращаем полученный ответ, если все прошло успешно.
	return resp, nil
}
