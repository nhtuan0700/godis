package data_structure

import "time"

type Obj struct {
	Value any
}

type Dict struct {
	dictStore        map[string]*Obj
	expiredDictStore map[string]uint64
}

func NewDict() *Dict {
	return &Dict{
		dictStore:        make(map[string]*Obj),
		expiredDictStore: make(map[string]uint64),
	}
}

func (d *Dict) GetExpiredDictStore() map[string]uint64 {
	return d.expiredDictStore
}

func (d *Dict) NewObj(k string, v any, ttlMs uint64) *Obj {
	obj := &Obj{
		Value: v,
	}

	if ttlMs > 0 {
		d.SetExpiry(k, ttlMs)
	}

	return obj
}

func (d *Dict) SetExpiry(k string, ttlMs uint64) {
	_, ok := d.expiredDictStore[k]
	if !ok {
		HashKeySpaceStat.Expire++
	}
	d.expiredDictStore[k] = uint64(time.Now().UnixMilli()) + ttlMs
}

func (d *Dict) GetExpiry(k string) (uint64, bool) {
	exp, ok := d.expiredDictStore[k]
	return exp, ok
}

func (d *Dict) HasExpired(k string) bool {
	if expired, ok := d.expiredDictStore[k]; ok {
		return expired <= uint64(time.Now().UnixMilli())
	}

	return false
}

func (d *Dict) Get(k string) *Obj {
	if obj, ok := d.dictStore[k]; ok {
		// delete epxired key in passive mode
		if d.HasExpired(k) {
			d.Del(k)
			return nil
		}
		return obj
	}

	return nil
}

func (d *Dict) Set(k string, obj *Obj) {
	_, ok := d.dictStore[k]
	if !ok {
		HashKeySpaceStat.Key++
	}
	d.dictStore[k] = obj
}

func (d *Dict) Del(k string) bool {
	if _, ok := d.expiredDictStore[k]; ok {
		delete(d.expiredDictStore, k)
		HashKeySpaceStat.Expire--
	}
	if _, ok := d.dictStore[k]; ok {
		delete(d.dictStore, k)
		HashKeySpaceStat.Key--
		return true
	}

	return false
}
