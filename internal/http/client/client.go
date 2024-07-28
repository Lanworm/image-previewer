package client

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Client структура.
type Client struct {
	client *http.Client
}

// NewHTTPClient функция для создания нового HTTP клиента.
func NewHTTPClient(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoRequest выполняет HTTP запрос с заданным методом, URL, телом запроса и заголовками.
func (c *Client) DoRequest(method string, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	// Создаем новый HTTP запрос с заданным методом, URL и телом запроса
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	// Добавляем переданные заголовки к запросу
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Выполняем HTTP запрос с помощью клиента
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func main() {
//	// Создаем новый HTTP клиент с таймаутом 10 секунд
//	client := NewHTTPClient(10 * time.Second)
//
//	// Создаем контекст с таймаутом 5 секунд
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	// Пример исходного запроса, чтобы получить заголовки
//	url := "https://jsonplaceholder.typicode.com/posts/1"
//	origReq, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		fmt.Println("Ошибка создания исходного запроса:", err)
//		return
//	}
//	origReq.Header.Add("Content-Type", "application/json")
//	origReq.Header.Add("Authorization", "Bearer some_token")
//
//	// Конвертируем заголовки из исходного запроса в map[string]string
//	headers := convertHeaders(origReq.Header)
//
//	// Выполняем новый запрос с заголовками из исходного запроса
//	response, err := client.DoRequest(ctx, "GET", url, nil, headers)
//	if err != nil {
//		fmt.Println("Ошибка выполнения запроса:", err)
//	} else {
//		defer response.Body.Close()
//		body, _ := io.ReadAll(response.Body)
//		fmt.Println("Ответ на GET запрос:", string(body))
//	}
//}
