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
	"github.com/stretchr/testify/assert"
)

// Картинка найдена в кэше.
func TestImageFoundInCache(t *testing.T) {
	imagePath := "image1.jpg"

	// Запрашиваем картинку с сервера
	resp, err := GetImage(imagePath)
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
	resp, err = GetImage(imagePath)
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

// RemoveFile удаляет указанный файл из директории int_test/test_files.
func RemoveFile(fileName string) error {
	filePath := fmt.Sprintf("/int_test/test_files/%s", fileName) // Полный путь к файлу

	// Удаляем файл
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ошибка при удалении файла %s: %w", filePath, err)
	}

	return nil // Возвращаем nil, если файл успешно удален
}

func RestoreFile(fileName string) error {
	srcFilePath := fmt.Sprintf("/int_test/backup_files/%s", fileName) // Путь к резервной копии файла
	dstFilePath := fmt.Sprintf("/int_test/test_files/%s", fileName)   // Путь, куда восстанавливаем файл

	// Открываем файл из резервной копии
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("ошибка при открытии резервного файла %s: %w", srcFilePath, err)
	}
	defer srcFile.Close()

	// Создаем файл в целевой директории
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла %s: %w", dstFilePath, err)
	}
	defer dstFile.Close()

	// Копируем содержимое из резервного файла в целевой файл
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("ошибка при копировании файла из %s в %s: %w", srcFilePath, dstFilePath, err)
	}

	return nil // Возвращаем nil, если файл успешно восстановлен
}

func IsRunningInContainer() bool {
	_, err := os.Stat("/proc/1/cgroup")
	return !os.IsNotExist(err)
}

// GetImage делает HTTP-запрос для получения изображения по указанному пути.
func GetImage(imagePath string) (*http.Response, error) {
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

	// Формируем базовый URL для запроса.
	baseURL := fmt.Sprintf("http://%s:8090/fill/200/300/%s:3080/images/%s", appURL, nginxURL, imagePath)

	// Создаем HTTP-клиент с таймаутом в 10 секунд.
	HTTPClient := client.NewHTTPClient(10 * time.Second)

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

// Дополнительные тесты для других сценариев...
