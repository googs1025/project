package lru

/*

 */

import (
	"container/list"
)



// 建立一个有map和双向链表的结构体
// 使用标准库中，list.list双向链表
// map[string]*list.Element  key:字符串 value:双向链表指针
// maxBytes:最大内存容量 nbytes:当前使用的内存
//
type Cache struct {
	maxBytes	int64
	nbytes 		int64
	ll 			*list.List
	cache		map[string]*list.Element
	OnEvicted	func(key string, value Value)
}

// entry 是双向链表节点的数据类型
type entry struct {
	key		string
	value	Value
	
}
// 实现Value接口
type Value interface {
	Len()	int // 用来返回占用内存大小
}

// 实例化Cache New()
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}

}

// 新增/修改方法
func (c *Cache) Add(key string, value Value)  {
	// 如果key存在，更新对应节点值，并且移到队尾
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 新增
		// 先在队尾新增entry 节点
		ele := c.ll.PushFront(&entry{key, value})
		// 在map中增加key的映射
		c.cache[key] = ele
		// 更新内存容量
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超过容量，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

	
}


// 查找方法
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 查找，找到移到队尾
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
	
}

// 删除方法
func (c *Cache) RemoveOldest()  {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}

	}

}

// 用来获取添加多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}

