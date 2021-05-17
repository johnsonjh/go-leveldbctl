package jleveldbctl

import (
	"fmt"
	"os"
	"path"

	jleveldb "github.com/johnsonjh/jleveldb/leveldb"
	util "github.com/johnsonjh/jleveldb/leveldb/util"
)

func dbexists(dbpath string) bool {
	_, err := os.Stat(path.Join(dbpath, "LOG"))
	return err == nil
}

func Init(dbpath string) error {
	if dbexists(dbpath) {
		return fmt.Errorf("%s was initialized as jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	return nil
}

func Put(dbpath string, key, value []byte) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	err = db.Put(key, value, nil)
	return err
}

func Get(dbpath string, key []byte) (string, bool, error) {
	if !dbexists(dbpath) {
		return "", false, fmt.Errorf("%s is not jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return "", false, fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	has, err := db.Has(key, nil)
	if err != nil {
		return "", false, fmt.Errorf("cannot open jleveldb")
	}
	if !has {
		return "", false, nil
	}

	value, err := db.Get(key, nil)
	if err != nil {
		return "", true, fmt.Errorf("cannot get value")
	}
	return string(value), true, nil
}

func Delete(dbpath string, key []byte) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	err = db.Delete(key, nil)
	return err
}

func Walk(dbpath string, f func(string, string)) error {
	if !dbexists(dbpath) {
		return fmt.Errorf("%s is not jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	s, err := db.GetSnapshot()
	if err != nil {
		return fmt.Errorf("cannot make snapshot jleveldb")
	}
	defer s.Release()

	i := s.NewIterator(nil, nil)
	for i.Next() {
		key := string(i.Key())
		value := string(i.Value())
		f(key, value)
	}

	return nil
}

func Search(dbpath string, key []byte) (string, bool, error) {
	if !dbexists(dbpath) {
		return "", false, fmt.Errorf("%s is not jleveldb", dbpath)
	}

	db, err := jleveldb.OpenFile(dbpath, nil)
	if err != nil {
		return "", false, fmt.Errorf("cannot open jleveldb")
	}
	defer db.Close()

	out := ""
	iter := db.NewIterator(util.BytesPrefix(key), nil)
	for iter.Next() {
		out = out + string(iter.Key()) + ": " + string(iter.Value()) + "\n"
	}
	iter.Release()

	return string(out), true, iter.Error()
}
