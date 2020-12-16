/*
	Package dictionary provides Dictionary data structure for holding participation data.
*/
package dictionary

import (
	"sync"

	"classboard/helper"
)

// Provides Dictionary structure

type NameKey string // user's name represent Dictionary Key

type ResultMap struct {
	Item map[int]int
	lock sync.RWMutex
}

type Dictionary struct {
	Item map[NameKey]*ResultMap
	lock sync.RWMutex
}

/*
	Dictionary methods
*/
// SetResultMap inserts ResultMap type with NameKey in Dictionary.
func (d *Dictionary) SetResultMap(name NameKey, result *ResultMap) {
	defer helper.CheckPanic()

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.Item == nil {
		d.Item = make(map[NameKey]*ResultMap)
	}
	d.Item[name] = result
}

// GetResultMapBasedName retrieves ResultMap type based on name given.
func (d *Dictionary) GetResultMapBasedName(name NameKey) *ResultMap {
	d.lock.RLock()
	defer d.lock.RUnlock()

	result, exist := d.Item[name]
	if !exist {
		return nil
	}
	return result
}

// DeleteResultMap removes ResultMap type based on name given.
func (d *Dictionary) DeleteResultMap(name NameKey) bool {
	defer helper.CheckPanic()

	d.lock.Lock()
	defer d.lock.Unlock()

	_, exist := d.Item[name]
	if exist {
		delete(d.Item, name)
	}
	return exist
}

// GetSize returns total number of name in Dictionary.
func (d *Dictionary) GetSize() int {
	d.lock.RLock()
	defer d.lock.RUnlock()

	size := len(d.Item)
	return size
}

// GetKeys returns array of names from Dictionary.
func (d *Dictionary) GetKeys() []NameKey {
	d.lock.RLock()
	defer d.lock.RUnlock()

	var dictKeys []NameKey
	dictKeys = []NameKey{}
	for key := range d.Item {
		dictKeys = append(dictKeys, key)
	}
	return dictKeys
}

/*
	ResultMap methods
*/
// SetValue inserts int type with question_id integer in ResultMap.
func (rm *ResultMap) SetValue(question_id int, result int) {
	defer helper.CheckPanic()

	rm.lock.Lock()
	defer rm.lock.Unlock()

	if rm.Item == nil {
		rm.Item = make(map[int]int)
	}
	rm.Item[question_id] = result
}
