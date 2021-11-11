package data

import (
	"errors"
	"frame/app/book/internal/biz"
	"sync"
	"time"
)

var errNotFound = errors.New("error: Not found")
var errExist = errors.New("error: record existed")
var errNotExist = errors.New("error: record not existed")
var errDialFailed = errors.New("error: dail failed")

type row struct {
	id    int
	name  string
	saled bool
}

type fakeDB map[int]*row

type BookRepo struct {
	db fakeDB
	sync.RWMutex
}

// simulate connect to db
func (d *fakeDB) Dial(host, port, user, password string) error {
	time.Sleep(time.Millisecond)
	if host == "localhost" && port == "3306" && user == "root" && password == "root" {
		return nil
	}
	return errDialFailed
}

var _ biz.BookRepo = new(BookRepo)

func (r *BookRepo) FindBookByID(id int) (*biz.Book, error) {
	r.RLock()
	defer r.RUnlock()
	book, ok := r.db[id]
	if !ok {
		return nil, errNotFound
	}
	return &biz.Book{
		ID:    book.id,
		Name:  book.name,
		Saled: book.saled,
	}, nil
}

func (r *BookRepo) SaveBook(book *biz.Book) (*biz.Book, error) {
	r.Lock()
	defer r.Unlock()
	if book.ID == 0 {
		id := len(r.db) + 1
		book.ID = id
		return r.createBook(book)
	}
	return r.updateBook(book)
}

func (r *BookRepo) createBook(book *biz.Book) (*biz.Book, error) {
	if _, ok := r.db[book.ID]; ok {
		return nil, errExist
	}
	r.db[book.ID] = &row{
		id:    book.ID,
		name:  book.Name,
		saled: book.Saled,
	}
	return book, nil
}

func (r *BookRepo) updateBook(book *biz.Book) (*biz.Book, error) {
	if _, ok := r.db[book.ID]; !ok {
		return nil, errNotExist
	}
	r.db[book.ID] = &row{
		id:    book.ID,
		name:  book.Name,
		saled: book.Saled,
	}
	return book, nil
}

func (r *BookRepo) DeleteBook(id int) error {
	_, ok := r.db[id]
	if !ok {
		return errNotFound
	}
	delete(r.db, id)
	return nil
}
