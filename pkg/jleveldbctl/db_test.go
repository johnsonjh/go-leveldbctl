package jleveldbctl_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	db "github.com/johnsonjh/jleveldbctl/pkg/jleveldbctl"
	u "github.com/johnsonjh/leaktestfe"
)

func TestInit(t *testing.T) {
	defer u.Leakplug(t)
	tmpdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpdir)

	dbdir := path.Join(tmpdir, "init")

	// init db
	err = db.Init(dbdir)
	if err != nil {
		t.Error(err)
	}

	// when db is initialized, "LOG" is exist under the dbdir
	if _, err := os.Stat(path.Join(dbdir, "LOG")); err != nil {
		t.Error(err)
	}
}

func TestCRUD(t *testing.T) {
	defer u.Leakplug(t)
	tmpdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpdir)

	dbdir := path.Join(tmpdir, "crud")

	// init db
	err = db.Init(dbdir)
	if err != nil {
		t.Error(err)
	}

	// const
	v := []byte("Foo")
	k := []byte("asdf")
	k2 := []byte("asdf2")

	// get from blank db
	value, ok, err := db.Get(dbdir, k)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Errorf("Key{%s} must not be exist.", k)
	}
	if value != "" {
		t.Error("value should be blank when key is not exist")
	}

	// set a value with key
	err = db.Put(dbdir, k, v)
	if err != nil {
		t.Error(err)
	}

	// get with key
	value, ok, err = db.Get(dbdir, k)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf("Key{%s} must be exist.", k)
	}
	if !bytes.Equal(v, []byte(value)) {
		t.Errorf("value should be %s but that is %s.", v, value)
	}

	// get with unset key
	value, ok, err = db.Get(dbdir, k2)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Errorf("Key{%s} must not be exist.", k2)
	}
	if value != "" {
		t.Error("value should be blank when key is not exist")
	}

	// delete
	err = db.Delete(dbdir, k)
	if err != nil {
		t.Error(err)
	}

	// key is deleted.
	value, ok, err = db.Get(dbdir, k)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Errorf("Key{%s} must not be exist.", k)
	}
	if value != "" {
		t.Error("value should be blank when key is not exist")
	}
}

func TestWalk(t *testing.T) {
	defer u.Leakplug(t)
	tmpdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpdir)

	dbdir := path.Join(tmpdir, "walk")

	// initialize
	err = db.Init(dbdir)
	if err != nil {
		t.Error(err)
	}

	// const
	keyvalue := map[string]string{
		"k1": "asdf",
		"k2": "fsda",
		"k3": "aaaa",
	}

	// set value
	for k, v := range keyvalue {
		err := db.Put(dbdir, []byte(k), []byte(v))
		if err != nil {
			t.Error(err)
		}
	}

	// walk
	actual_keyvalue := map[string]string{}
	err = db.Walk(dbdir, func(k, v string) {
		actual_keyvalue[k] = v
	})
	if err != nil {
		t.Error(err)
	}

	// Test
	if !reflect.DeepEqual(keyvalue, actual_keyvalue) {
		t.Errorf("k,v found with walk are not equal. expected: %v, actual: %v", keyvalue, actual_keyvalue)
	}
}

func TestCheckingExistenceDB(t *testing.T) {
	defer u.Leakplug(t)
	tmpdir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpdir)

	dbdir := path.Join(tmpdir, "existence")

	// Uninitialized Get, Delete, Put, Walk
	err = db.Put(dbdir, []byte("k"), []byte("v"))
	if err == nil {
		t.Error("Put not check whether db is initialized")
	}

	err = db.Delete(dbdir, []byte("k"))
	if err == nil {
		t.Error("Delete not check whether db is initialized")
	}

	value, ok, err := db.Get(dbdir, []byte("k"))
	if err == nil {
		t.Error("Get not check whether db is initialized")
	}

	if ok {
		t.Error("Get is returing wrong value that 'ok' should be false")
	}

	if value != "" {
		t.Error("Get is returing wrong value that 'value' should be blank string")
	}

	err = db.Walk(dbdir, func(k, v string) {
		t.Error("handling functions should be called when db is not initialized")
	})
	if err == nil {
		t.Error("Walk not check whether db is initialized")
	}

	// Initialized DB with init
	err = db.Init(dbdir)
	if err != nil {
		t.Error(err)
	}

	// Check
	err = db.Init(dbdir)
	if err == nil {
		t.Error("Init not check whether db is initialized")
	}
}
