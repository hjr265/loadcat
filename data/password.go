// Copyright 2015 The Loadcat Authors. All rights reserved.

package data

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
)

type Password struct {
	Hash []byte
	Salt []byte
}

func (p *Password) Set(clear string) error {
	p.Salt = make([]byte, 32)
	_, err := rand.Read(p.Salt)
	if err != nil {
		return err
	}
	sum := sha1.Sum([]byte(string(p.Salt) + clear))
	p.Hash = sum[:]
	return nil
}

func (p *Password) Equal(clear string) bool {
	sum := sha1.Sum([]byte(string(p.Salt) + clear))
	return bytes.Equal(p.Hash, sum[:])
}
