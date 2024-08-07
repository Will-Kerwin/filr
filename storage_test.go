package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpictures"
	pathKey := CASPathTransformFunc(key)
	expectedFileName := "52439f5e49ba33b4b9b2373dcd6bfc9e97cafb7f"
	expectedPathname := "52439/f5e49/ba33b/4b9b2/373dc/d6bfc/9e97c/afb7f"
	if pathKey.Pathname != expectedPathname {
		t.Errorf("%s want %s\n", pathKey.Pathname, expectedPathname)
	}
	if pathKey.Filename != expectedFileName {
		t.Errorf("%s want %s\n", pathKey.Filename, expectedFileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	id := generateID()
	defer tearDown(t, s)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("Some jpg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("expected to have key in %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		fmt.Println(string(b))

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); ok {
			t.Errorf("expected to not have key in %s", key)
		}
	}

}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)

}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
