package int_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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
	err = RemoveFileFromDocker(imagePath)
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

func RemoveFileFromDocker(imagePath string) error {
	// Команда для удаления файла из контейнера
	cmd := exec.Command("docker-compose", "exec", "nginx", "rm", "/usr/share/nginx/html/images/"+imagePath)

	// Выполняем команду
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ошибка при удалении файла: %w", err)
	}

	// Проверка, что файл был успешно удален
	checkCmd := exec.Command("docker-compose", "exec", "nginx", "ls", "/usr/share/nginx/html/images")
	output, err := checkCmd.Output()
	if err != nil {
		return fmt.Errorf("ошибка при проверке наличия файла: %w", err)
	}
	// Преобразуем вывод в строку и проверяем, что файла нет в выводе
	if string(output) == imagePath {
		return fmt.Errorf("файл %s все еще существует после удаления", imagePath)
	}
	return nil
}

func IsRunningInContainer() bool {
	_, err := os.Stat("/proc/1/cgroup")
	return !os.IsNotExist(err)
}

// GetImage делает HTTP-запрос для получения изображения по указанному пути.
func GetImage(imagePath string) (*http.Response, error) {
	var serviceURL string

	// Определяем URL сервиса в зависимости от того, запущен ли код в контейнере.
	if IsRunningInContainer() {
		serviceURL = "previewer-web-nginx" // Используем имя сервиса в Docker.
	} else {
		serviceURL = "localhost" // Используем localhost для локальной разработки.
	}

	// Формируем базовый URL для запроса.
	baseURL := fmt.Sprintf("http://localhost:8090/fill/200/300/%s:3080/images/%s", serviceURL, imagePath)

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
