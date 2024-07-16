package client

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

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
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// Создаем контекст с таймаутом, основываясь на таймауте клиента
	ctx, cancel := context.WithTimeout(context.Background(), c.client.Timeout)
	defer cancel()

	// Присваиваем созданный контекст к запросу
	req = req.WithContext(ctx)

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
