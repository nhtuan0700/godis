package core

import (
	"sync"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

var dictStore *data_structure.Dict
var setStore map[string]*data_structure.SimpleSet
var once sync.Once

func init() {
	once.Do(func() {
		dictStore = data_structure.NewDict()
		setStore = make(map[string]*data_structure.SimpleSet)
	})
}
