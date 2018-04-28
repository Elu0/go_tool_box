package util

import (
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"unsafe"
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

type LRUCacheAddress *LRUCache

func Get(cache *LRUCacheAddress, key interface{}) (interface{}, error) {
	newCache := LRUCache{}
	err := DeepCopy(*cache, &newCache)
	if err != nil {
		return nil, err
	}

	retValue := newCache.Get(key)

	for {
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(cache)),
			unsafe.Pointer(*cache), unsafe.Pointer(&newCache)) {
			break
		}
	}
	return retValue, nil
}

func Set(cache *LRUCacheAddress, key, value interface{}) error {
	newCache := LRUCache{}
	err := DeepCopy(*cache, &newCache)
	if err != nil {
		return err
	}

	newCache.Set(key, value)

	for {
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(cache)),
			unsafe.Pointer(*cache), unsafe.Pointer(&newCache)) {
			break
		}
	}
	return nil
}

func DeepCopy(src, dst *LRUCache) error {
	if src == nil {
		err := errors.New("src cannot be empty")
		return err
	}
	if dst == nil {
		err := errors.New("dst cannot be empty")
		return err
	}
	dst.Init(src.Cap)
	pd := dst.Head
	for p := src.Head.Next; p != nil; p = p.Next {
		if p == nil {
			err := errors.New("node ptr cannot be empty")
			return err
		}
		dNode := Node{Key: p.Key, Value: p.Value}
		keyStr := formatKeyToStr(p.Key)
		dst.MapContent[keyStr] = &dNode
		pd.Next = &dNode
		dNode.Prev = pd
		pd = pd.Next
	}
	dst.Tail = pd
	return nil
}

func (listMap *LRUCache) Init(capacity int) {
	listMap.Cap = capacity
	listMap.Head = &Node{Key: nil, Value: nil}
	listMap.Tail = listMap.Head
	listMap.MapContent = make(map[string]*Node)
}

func (listMap *LRUCache) Get(key interface{}) interface{} {
	keyStr := formatKeyToStr(key)
	node, ok := listMap.MapContent[keyStr]
	if !ok {
		return nil
	}
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	listMap.appendTail(node)
	return node.Value
}

func (listMap *LRUCache) Set(key interface{}, value interface{}) {
	keyStr := formatKeyToStr(key)
	existsNode, ok := listMap.MapContent[keyStr]
	if ok {
		existsNode.Value = value
		existsNode.Prev.Next = existsNode.Next
		existsNode.Next.Prev = existsNode.Prev
		listMap.appendTail(existsNode)
		return
	}
	if len(listMap.MapContent) == listMap.Cap {
		delNode := listMap.Head.Next
		listMap.Head.Next = listMap.Head.Next.Next
		listMap.Head.Next.Prev = listMap.Head
		delKeyStr := formatKeyToStr(delNode.Key)
		delete(listMap.MapContent, delKeyStr)
	}
	newNode := &Node{Key: key, Value: value}
	listMap.appendTail(newNode)
	listMap.MapContent[keyStr] = newNode
}

func (listMap *LRUCache) appendTail(node *Node) {
	node.Next = nil
	node.Prev = listMap.Tail
	listMap.Tail.Next = node
	listMap.Tail = node
}

func formatKeyToStr(key interface{}) string {
	switch k := key.(type) {
	case int64:
		return strconv.FormatInt(k, 10)
	case int:
		return strconv.Itoa(k)
	case float64:
		return strconv.FormatFloat(k, 'f', -1, 64)
	case string:
		return k
	default:
		return fmt.Sprint(k)
	}
}
