package project_for_practice

import (
	"fmt"
	"geecache/singleflight"
	"log"
	"sync"
)

// Group 为最核心的数据结构，负责与用户交互，并且控制缓存的存储和获取
// 一个Group是一个缓存命名空间，每个Group都是一个name
// getter 是缓存没命中时，获取源数据的回调
// mainCache 并发缓存
type Group struct {
	name string
	getter Getter
	mainCache cache

	peers PeerPicker

	loader *singleflight.Group
}

// 具体不清楚啥意思
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}



var (
	mu	sync.RWMutex
	groups = make(map[string]*Group)
)
// 实例化
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}
// GetGroup 使用读锁
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get 方法
func (g *Group) Get(key string) (ByteView, error) {
	// 如果没有key，返回错误
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 从mainCachet查找缓存 ，有！
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	// 没有，调用load方法
	return g.load(key)
}
// load方法调用getLocally
func (g *Group) load(key string) (value ByteView, err error) {

	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)

	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return

}

// getLocally 调用用户回调函数获取源数据，并添加到缓存中
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		b: cloneBytes(bytes),
	}
	g.populateCache(key, value)
	return value, nil

}

func (g *Group) populateCache(key string, value ByteView)  {
	g.mainCache.add(key, value)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}



func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}