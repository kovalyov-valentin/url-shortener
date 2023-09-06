package random

import (
	"math/rand"
	"time"
)

var character = []rune("abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789",
)

func NewRandomString(size int64) string {
	// Генерирует псевдослучайные числа
	// Устанавливаем сид для генератора рандомных чисел, чтобы он каждый раз генерировал разные alias
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Буфер для символов
	b := make([]rune, size)
	for i := range b {
		b[i] = character[rnd.Intn(len(character))]
	}
	return string(b)
}
