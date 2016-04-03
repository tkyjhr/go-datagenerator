package main

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"sync"
	"time"
)

// ZeroReader implements io.Reader's Read(p []byte). It fills data with 0x00.
var ZeroReader = &zeroReader{}

// FfReader implements io.Reader's Read(p []byte). It fills data with 0xFF.
var FfReader = &ffReader{}

// RandReader is rand.New(rand.NewSource(time.Now().UnixNano()))
var RandReader = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomAlphabetReader implements io.Reader's Read(p []byte). It fills data with random alphabet characters.
var RandomAlphabetReader = &randAlphabetReader{}

// RandomAlphabetNumericCharacterReader implements io.Reader's Read(p []byte). It fills data with random alphabet and numeric characters.
var RandomAlphabetNumericCharacterReader = &randAlphaNumericReader{}

type zeroReader struct{}

func (r *zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

type ffReader struct{}

func (r *ffReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0xff
	}
	return len(p), nil
}

type randAlphabetReader struct{}

func (r *randAlphabetReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte('A' + rand.Intn('Z' - 'A') + ('a' - 'A') * rand.Intn(2))
	}
	return len(p), nil
}

type randAlphaNumericReader struct{}

const alphanumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (r *randAlphaNumericReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return len(p), nil
}

// CounterReader is the interface that wraps the Reset() method in addition to io.Reader's Read(p []byte) method.
type CounterReader interface {
	Read(p []byte) (int, error)
	Reset()
}

// ByteCounterReader implements io.Reader's Read(p []byte).
// Data starts from 0x00 and goes up to 0xFF, then returns to 0x00 again.
type ByteCounterReader struct {
	ctr byte
	lk  sync.Mutex
}

func (r *ByteCounterReader) Read(p []byte) (n int, err error) {
	r.lk.Lock()
	defer r.lk.Unlock()
	for i := range p {
		p[i] = r.ctr
		r.ctr++
	}
	return len(p), nil
}

// Reset set the counter value to 0.
func (r *ByteCounterReader) Reset() {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.ctr = 0
}

// Uint16CounterReader implements io.Reader's Read(p []byte).
// Data starts from 0x0000 and goes up to 0xFFFF, then returns to 0x0000 again (BigEndian).
type Uint16CounterReader struct {
	ctr uint16
	lk  sync.Mutex
}

func (r *Uint16CounterReader) Read(p []byte) (n int, err error) {
	l := len(p)
	if l % 2 != 0 {
		return 0, errors.New("length must be multiple of 2")
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	for i := 0; i < l; i += 2 {
		binary.BigEndian.PutUint16(p[i:i + 2], r.ctr)
		r.ctr++
	}
	return len(p), nil
}

// Reset set the counter value to 0.
func (r *Uint16CounterReader) Reset() {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.ctr = 0
}

// Uint32CounterReader implements io.Reader's Read(p []byte).
// Data starts from 0x00000000 and goes up to 0xFFFFFFFF, then returns to 0x00000000 again (BigEndian).
type Uint32CounterReader struct {
	ctr uint32
	lk  sync.Mutex
}

func (r *Uint32CounterReader) Read(p []byte) (n int, err error) {
	l := len(p)
	if l % 4 != 0 {
		return 0, errors.New("length must be multiple of 4")
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	for i := 0; i < l; i += 4 {
		binary.BigEndian.PutUint32(p[i:i + 4], r.ctr)
		r.ctr++
	}
	return len(p), nil
}

// Reset set the counter value to 0.
func (r *Uint32CounterReader) Reset() {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.ctr = 0
}

// Uint64CounterReader implements io.Reader's Read(p []byte).
// Data starts from 0x0000000000000000 and goes up to 0xFFFFFFFFFFFFFFFF, then returns to 0x0000000000000000 again (BigEndian).
type Uint64CounterReader struct {
	ctr uint64
	lk  sync.Mutex
}

func (r *Uint64CounterReader) Read(p []byte) (n int, err error) {
	l := len(p)
	if l % 8 != 0 {
		return 0, errors.New("length must be multiple of 8")
	}
	r.lk.Lock()
	defer r.lk.Unlock()
	for i := 0; i < l; i += 8 {
		binary.BigEndian.PutUint64(p[i:i + 8], r.ctr)
		r.ctr++
	}
	return len(p), nil
}

// Reset set the counter value to 0.
func (r *Uint64CounterReader) Reset() {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.ctr = 0
}
