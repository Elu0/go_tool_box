package util

import (
	"fmt"
	"strconv"
)

type Node struct {
	Key   interface{}
	Value interface{}
	Next  *Node
	Prev  *Node
}

type LRUCache struct {
	Head       *Node
	Tail       *Node
	Cap        int
	MapContent map[string]*Node
}

func (cache *LRUCache) Init(capacity int) {
	cache.Cap = capacity
	cache.Head = &Node{Key: nil, Value: nil}
	cache.Tail = &Node{Key: nil, Value: nil}
	cache.Head.Next = cache.Tail
	cache.Tail.Prev = cache.Head
	cache.MapContent = make(map[string]*Node)
}

func (cache *LRUCache) Get(key interface{}) (value interface{}) {
	keyStr := cache.formatKeyToStr(key)
	node, ok := cache.MapContent[keyStr]
	if !ok {
		value = nil
		return
	}
	value = node.Value
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	cache.appendTail(node)
	return
}

func (cache *LRUCache) Set(key interface{}, value interface{}) {
	keyStr := cache.formatKeyToStr(key)
	existsNode, ok := cache.MapContent[keyStr]
	if ok {
		cache.MapContent[keyStr].Value = value
		existsNode.Prev.Next = existsNode.Next
		existsNode.Next.Prev = existsNode.Prev
		cache.appendTail(existsNode)
		return
	}
	if len(cache.MapContent) == cache.Cap {
		delNode := cache.Head.Next
		cache.Head.Next = cache.Head.Next.Next
		cache.Head.Next.Prev = cache.Head
		delKeyStr := cache.formatKeyToStr(delNode.Key)
		delete(cache.MapContent, delKeyStr)
	}
	newNode := &Node{Key: key, Value: value}
	cache.appendTail(newNode)
	cache.MapContent[keyStr] = newNode
}

func (cache *LRUCache) appendTail(node *Node) {
	node.Next = cache.Tail
	node.Prev = cache.Tail.Prev
	cache.Tail.Prev.Next = node
	cache.Tail.Prev = node
}

func (cache *LRUCache) formatKeyToStr(key interface{}) (keyStr string) {
	keyStr = ""
	switch k := key.(type) {
	case int64:
		keyStr = strconv.FormatInt(k, 10)
	case int:
		keyStr = strconv.Itoa(k)
	case float64:
		keyStr = strconv.FormatFloat(k, 'f', -1, 64)
	case string:
		keyStr = k
	default:
		keyStr = fmt.Sprint(k)
	}
	return
}
