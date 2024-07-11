package lrucache

import (
	"fmt"
	"image"
	"sync"

	"github.com/Lanworm/image-previewer/internal/storage"
)

type Key string

type Cache interface {
	Set(key Key, value image.Image) bool
	Get(key Key) (image.Image, bool)
	Clear()
	InitCache(path string) error
}

type CacheListItem struct {
	value image.Image
	key   Key
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
	storage  storage.Storage
}

func NewCache(capacity int, storage storage.Storage) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		storage:  storage,
	}
}

func (c *lruCache) Set(key Key, value image.Image) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.Lock()
	defer c.mu.Unlock()
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

func (c *lruCache) InitCache(folderPath string) error {
	fileNames, err := c.storage.GetFileList(folderPath)
	if err != nil {
		return err
	}

	for _, fileName := range fileNames {
		imgFile, err := c.storage.Get(fileName)
		if err != nil {
			return err
		}

		c.Set(Key(fileName), imgFile)
		fmt.Printf("added to the cache: %s\n", fileName)
	}

	return nil
}
