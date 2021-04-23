package main

type Set map[string]struct{}

func Newset() (Set) {
	return make(map[string]struct{})
}

func (self Set) delete(key string) {
	delete(self, key)
}

func (self Set) Add(key string) Set {
	self[key] = struct{}{}
	return self
}

func (self Set) Exists(key string) bool {
	_, ok := self[key]
	return ok
}
func (self Set) Clear() {
	self = make(map[string]struct{})
}
func (self Set) Keys() ([]string) {  //keys 去重的切片
	count := len(self)
	if count == 0 {
		return []string{}
	}
	sl := []string{}
	for k := range self {
		sl = append(sl, k)
	}
	return sl
}
