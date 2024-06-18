package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpictures"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "52439f5e49ba33b4b9b2373dcd6bfc9e97cafb7f"
	expectedPathname := "52439/f5e49/ba33b/4b9b2/373dc/d6bfc/9e97c/afb7f"
	if pathKey.Pathname != expectedPathname {
		t.Errorf("%s want %s\n", pathKey.Pathname, expectedPathname)
	}
	if pathKey.Original != expectedOriginalKey {
		t.Errorf("%s want %s\n", pathKey.Original, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("Some jpg bytes"))
	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}

}
