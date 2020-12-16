package dictionary

import (
	"reflect"
	"testing"
)

func TestSetResultMap(t *testing.T) {
	var dict *Dictionary = &Dictionary{}
	var rm *ResultMap = &ResultMap{}
	rm.SetValue(1, 1)
	rm.SetValue(2, 0)
	rm.SetValue(3, -1)

	var name NameKey = "test1"
	dict.SetResultMap(name, rm)
	result := dict.GetResultMapBasedName(name)
	if result == nil {
		t.Errorf("SetResultMap failed to set %s into dictionary\n", name)
	}
}

func TestGetResultMapBasedName(t *testing.T) {
	var dict *Dictionary = &Dictionary{}
	var rm *ResultMap = &ResultMap{}
	rm.SetValue(1, 1)
	rm.SetValue(2, 0)
	rm.SetValue(3, -1)

	var name NameKey = "test1"
	dict.SetResultMap(name, rm)

	// retrieve uninitialize NameKey
	result := dict.GetResultMapBasedName(NameKey("test2"))
	if result != nil {
		t.Errorf("GetResultMapBasedName should not retrieve ResultMap from NameKey : %s\n", name)
	}

	// retrieve initialize NameKey
	result = dict.GetResultMapBasedName(name)
	if result == nil {
		t.Errorf("GetResultMapBasedName failed to retrieve ResultMap from NameKey : %s\n", name)
	}
}

func TestDeleteResultMap(t *testing.T) {
	var dict *Dictionary = &Dictionary{}
	var rm *ResultMap = &ResultMap{}
	rm.SetValue(1, 1)
	rm.SetValue(2, 0)
	rm.SetValue(3, -1)

	// delete uninitialize NameKey
	result := dict.DeleteResultMap(NameKey("test2"))
	if result == true {
		t.Error("DeleteResultMap failed. Expecting false when NameKey is not initialize yet")
	}

	var name NameKey = "test1"
	dict.SetResultMap(name, rm)

	// delete initialize NameKey
	result = dict.DeleteResultMap(name)
	if result == false {
		t.Errorf("DeleteResultMap failed. Expecting true when NameKey is already initialize")
	}
}

func TestGetSize(t *testing.T) {
	var dict *Dictionary = &Dictionary{}
	var rm *ResultMap = &ResultMap{}
	rm.SetValue(1, 1)
	rm.SetValue(2, 0)
	rm.SetValue(3, -1)

	var nameKeys []NameKey = []NameKey{"test1", "test2", "test3"}
	// initialize 3 NameKey
	dict.SetResultMap(nameKeys[0], rm)
	dict.SetResultMap(nameKeys[1], rm)
	dict.SetResultMap(nameKeys[2], rm)

	var expectedResult int = 3
	result := dict.GetSize()
	if result != expectedResult {
		t.Errorf("GetSize failed. Expecting %d but returned %d\n", expectedResult, result)
	}
}

func TestGetKeys(t *testing.T) {
	var dict *Dictionary = &Dictionary{}
	var rm *ResultMap = &ResultMap{}
	rm.SetValue(1, 1)
	rm.SetValue(2, 0)
	rm.SetValue(3, -1)

	var nameKeys []NameKey = []NameKey{"test1", "test2", "test3"}
	// initialize 3 NameKey
	dict.SetResultMap(nameKeys[0], rm)
	dict.SetResultMap(nameKeys[1], rm)
	dict.SetResultMap(nameKeys[2], rm)

	result := dict.GetKeys()
	if !reflect.DeepEqual(result, nameKeys) {
		t.Errorf("GetKeys failed. Expecting %+v but returned %+v\n", nameKeys, result)
	}
}

func TestSetValue(t *testing.T) {
	var rm *ResultMap = &ResultMap{}

	expectedMap := make(map[int]int)
	expectedMap[1] = 1
	expectedMap[2] = 0
	expectedMap[3] = -1
	rm.SetValue(1, expectedMap[1])
	rm.SetValue(2, expectedMap[2])
	rm.SetValue(3, expectedMap[3])

	for i, v := range rm.Item {
		switch i {
		case 1:
			if v != expectedMap[1] {
				t.Errorf("SetValue failed in case 1. %d value doesn't match with expected %d in\n", v, expectedMap[1])
			}
		case 2:
			if v != expectedMap[2] {
				t.Errorf("SetValue failed in case 2. %d value doesn't match with expected %d\n", v, expectedMap[2])
			}
		case 3:
			if v != expectedMap[3] {
				t.Errorf("SetValue failed in case 3. %d value doesn't match with expected %d\n", v, expectedMap[3])
			}
		default:
			t.Errorf("SetValue failed. %d key not found\n", i)
		}
	}
}
