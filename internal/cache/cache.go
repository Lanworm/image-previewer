package lrucache

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
)

type Key string

type Cache interface {
	Set(key Key, value image.Image) bool
	Get(key Key) (image.Image, bool)
	Clear()
}

type CacheListItem struct {
	value image.Image
	key   Key
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value image.Image) bool {
	cacheListItem, ok := c.items[key]

	if ok {
		cacheItem := cacheListItem.Value.(CacheListItem)
		cacheItem.value = value
		cacheListItem.Value = cacheItem

		c.queue.MoveToFront(cacheListItem)

		return true
	}

	newCacheItem := CacheListItem{
		value: value,
		key:   key,
	}

	if c.queue.Len() >= c.capacity {
		lastListItem := c.queue.Back()
		cad := lastListItem.Value.(CacheListItem)
		delete(c.items, cad.key)
	}

	c.items[key] = c.queue.PushFront(newCacheItem)

	return false
}

func (c *lruCache) Get(key Key) (image.Image, bool) {
	cacheItem, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(cacheItem)
		cad := cacheItem.Value.(CacheListItem)

		return cad.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func InitCache(folderPath string, cap int, cache Cache) error {
	file, err := os.Open(folderPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for i, fileInfo := range fileInfos {
		if fileInfo.IsDir() && i <= cap {
			continue
		}
		filename := fileInfo.Name()
		filePath := filepath.Join(folderPath, filename)

		imgFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			return err
		}
		cache.Set(Key(filename), img)
		fmt.Printf("added to the cache: %s\n", filename)
	}
	return nil
}
