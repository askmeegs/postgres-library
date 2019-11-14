// adapted from https://github.com/go-pg/pg

package main

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"math/rand"
	"os"
	"time"
)

type Story struct {
	Id    int64
	Title string
}

func (s Story) String() string {
	return fmt.Sprintf("Story<%d %s>", s.Id, s.Title)
}

var db *pg.DB
var titles = []string{
	"A Tale of Two Cities",
	"The Golden Compass",
	"To the Lighthouse",
	"The Poet X",
	"The Great Gatsby",
	"Infinite Jest",
	"Fun Home",
	"Pride and Prejudice",
	"A Tale for the Time Being",
	"Homecoming"}

func initialize() error {
	pw := os.Getenv("POSTGRES_PASSWORD")
	db = pg.Connect(&pg.Options{
		Addr:     "postgresql-1-vm.default.svc.cluster.local:5432",
		User:     "postgres",
		Password: pw,
	})
	return createSchema(db)
}

func insertAndSelect(t time.Time) {
	if err := initialize(); err != nil {
		panic(err)
	}
	defer db.Close()
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	title := titles[rand.Intn(len(titles))]
	story1 := &Story{
		Title: title,
	}
	err := db.Insert(story1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("✅ inserted %s\n", title)
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*Story)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	doEvery(2*time.Second, insertAndSelect)
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
