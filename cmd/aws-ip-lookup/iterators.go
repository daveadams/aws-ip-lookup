package main

import (
	"bufio"
	"os"
	"strings"
)

type StringIterator interface {
	Next() (string, bool)
}

type SliceIterator struct {
	list  []string
	index int
}

func NewSliceIterator(list []string) *SliceIterator {
	sliceCopy := make([]string, 0, len(list))
	for _, item := range list {
		sliceCopy = append(sliceCopy, item)
	}
	return &SliceIterator{list: sliceCopy}
}

func (it *SliceIterator) Next() (string, bool) {
	if it.index >= len(it.list) {
		return "", false
	}
	rv := it.list[it.index]
	it.index++
	return rv, true
}

type StdinIterator struct {
	reader *bufio.Reader
}

func NewStdinIterator() *StdinIterator {
	return &StdinIterator{reader: bufio.NewReader(os.Stdin)}
}

func (it *StdinIterator) Next() (string, bool) {
	line, err := it.reader.ReadSlice('\n')
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(line)), true
}
