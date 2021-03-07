package api

import (
	"math/rand"
	"time"
)

type (
	Rnder interface {
		Rnd() int64
	}

	rnd struct {
	}
)

func NewRnder() *rnd {
	return &rnd{}
}

func (r *rnd) Rnd() int64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Int63()
}
