package lru

import "container/list"

/**
* LRU(Least Recently Used)
* 最近最少使用策略，LRU认为，如果数据最近被访问过，那么将来被访问的改了也会更高
* 维护一个队列，将实际的值存进这个队列中，而队列的指针则保存在一个map中,当某条记录被访问了，则移到队尾，队首就肯定是最近访问最少的了，要淘汰的话就淘汰这条数据就好了
 */

type Cache struct {
	maxBytes  int64                         // 内存允许使用的最大内存
	nbytes    int64                         // 当前已经使用的内存
	ll        *list.List                    // 淘汰队列
	cache     map[string]*list.Element      // key是缓存的键，value是队列的位置
	OnEvicted func(key string, value Value) // 是某条记录被移除时的回调函数
}

// entry 是实际存进队列中的key和value
type entry struct {
	key   string
	value Value
}

// 为了通用性，不是所有的值都可以存进去，而是实现了这个接口的值才可以存进去
type Value interface {
	Len() int
}

// 新建一个cache实例
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{maxBytes: maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 有两个步骤，1. 是将内容取出来 2. 将当前节点放到队尾，表示最近访问过
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		v := element.Value.(*entry) //当前节点中值
		return v.value, true
	}
	return
}

// RemoveOldest 淘汰缓存，移除最近最少访问的节点也就是队首
func (c *Cache) RemoveOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		v := element.Value.(*entry) // 这里除了删除队列中的值之外，还要删除map中的值，将key也存在队列中，这样就方便了,也表明了为什么要设计一个entry结构体
		delete(c.cache, v.key)
		c.nbytes -= int64(len(v.key)) + int64(v.value.Len()) // 除了删除值的大小之外，key也是有内存的
		if c.OnEvicted != nil {
			c.OnEvicted(v.key, v.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		v := element.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(v.value.Len()) // 对比新旧值的大小，存进去差值
		v.value = value
	} else {
		element := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = element
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes { // 如果超出了最大值的限制的话，就差将最少最近访问的节点淘汰
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
