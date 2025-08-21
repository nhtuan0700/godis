package core

import (
	"sync"

	"github.com/nhtuan0700/godis/internal/core/data_structure"
)

var dictStore *data_structure.Dict
var once sync.Once

func init() {
	once.Do(func ()  {
		dictStore = data_structure.NewDict()
	})
}
