package lrucache

import (
	"image"
	"testing"

	"github.com/Lanworm/image-previewer/internal/storage/filestorage"
	"github.com/stretchr/testify/require"
)

func TestLRUCache(t *testing.T) {
	storage := filestorage.NewFileStorage("../../test_images")
	cache := NewCache(2, storage)

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
	storage := filestorage.NewFileStorage("../../test_images")
	testCache := NewCache(capacity, storage)

	// Создание временного файла с изображением для теста
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	err := storage.Set(img, "temp_image.jpg")
	if err != nil {
		return
	}
	err = testCache.InitCache("../../test_images")
	require.NoError(t, err, "Ошибка при инициализации кеша изображений")

	// Проверка добавления изображения в кеш
	retrievedImg, found := testCache.Get(Key("temp_image.jpg"))
	require.True(t, found, "Изображение 'temp_image.jpg' не найдено в кеше")
	require.NotNil(t, retrievedImg, "Изображение 'temp_image.jpg' не было добавлено в кеш")

	// Удаление временного файла после теста
	err = storage.Delete("temp_image.jpg")
	if err != nil {
		return
	}
}
