package lrucache

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLRUCache(t *testing.T) {
	cache := NewCache(2)

	// Проверка добавления и получения изображения из кеша
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	cache.Set(Key("image1"), img1)

	retrievedImg1, found1 := cache.Get(Key("image1"))
	require.True(t, found1, "Изображение 'image1' не найдено в кеше")
	require.Equal(t, img1, retrievedImg1, "Изображение 'image1' не соответствует ожидаемому")

	// Проверка замещения изображения в кеше
	img2 := image.NewRGBA(image.Rect(0, 0, 200, 200))
	cache.Set(Key("image2"), img2)

	img3 := image.NewRGBA(image.Rect(0, 0, 300, 300))
	cache.Set(Key("image3"), img3)

	_, found2 := cache.Get(Key("image1"))
	require.False(t, found2, "Изображение 'image1' должно быть замещено")

	// Проверка очистки кеша
	cache.Clear()
	require.Empty(t, cache.(*lruCache).items, "Элементы кеша не были очищены")
	require.Zero(t, cache.(*lruCache).queue.Len(), "Длина очереди кеша должна быть нулевой")
}

func TestInitCache(t *testing.T) {
	capacity := 2
	testCache := NewCache(capacity)

	// Создание временного файла с изображением для теста
	tempImageFilename := filepath.Join("../../test_images", "temp_image.jpg")
	createTempImageFile(tempImageFilename)

	err := InitCache("../../test_images", testCache)
	require.NoError(t, err, "Ошибка при инициализации кеша изображений")

	// Проверка добавления изображения в кеш
	retrievedImg, found := testCache.Get(Key("temp_image.jpg"))
	require.True(t, found, "Изображение 'temp_image.jpg' не найдено в кеше")
	require.NotNil(t, retrievedImg, "Изображение 'temp_image.jpg' не было добавлено в кеш")

	// Удаление временного файла после теста
	err = os.Remove(tempImageFilename)
	if err != nil {
		panic(err)
	}
}

func createTempImageFile(filename string) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	err := os.MkdirAll(filepath.Dir(filename), 0o777)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := jpeg.Encode(file, img, nil); err != nil {
		panic(err)
	}
}
