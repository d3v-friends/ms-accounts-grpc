package fnFakeit

import (
	"github.com/brianvoe/gofakeit"
	"time"
)

func Init() {
	gofakeit.Seed(time.Now().Unix())
}
