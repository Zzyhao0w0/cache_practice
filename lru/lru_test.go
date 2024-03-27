package lru

import "testing"

type str string

func (d str) Len() int {
	return len(d)
}

// 测试获取缓存
func TestGet(t *testing.T) {
	var (
		aStr str = "key_test_get_value"
		aKey     = "key_test_get_key"
	)
	lru := New(0, nil)
	lru.Add(aKey, aStr)
	if v, ok := lru.Get(aKey); !ok || v.(str) != aStr {
		t.Fatalf("cache hit %s failed", aKey)
	}

	if _, ok := lru.Get("other_key"); ok {
		t.Fatalf("failed cache hit a miss key")
	}
}

// 测试当内存超过了限制时会不会触发移除无用节点
func TestRemoveoldest(t *testing.T) {
	var (
		key1, key2, key3           = "test_key_1", "test_key_2", "test_key_3"
		value1, value2, value3 str = "test_value_1", "test_value_2", "test_value_3"
		cap                        = len(key1 + key2 + string(value1) + string(value2))
	)

	lru := New(int64(cap), nil)
	lru.Add(key1, value1)
	lru.Add(key2, value2)
	lru.Add(key3, value3)

	if _, ok := lru.Get(key1); ok || lru.Len() != 2 {
		t.Fatalf("cache removeoldest failed ")
	}

}

// 测试回调函数能否被调用
func TestOnEvicted(t *testing.T) {

}
