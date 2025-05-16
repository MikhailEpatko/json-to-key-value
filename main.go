package main

import (
	iterjson "ezpkg.io/iter.json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	driverName   = "postgres"
	dsn          = "host=127.0.0.1 port=5432 user=user password=password dbname=postgres sslmode=disable"
	tables       = map[string]string{"table1": "column", "table2": "column", "table3": "column"}
	urlPattern   = "://"
	colorPattern = "#"
)

type UiEntity struct {
	Id   int64  `db:"id"`
	Json []byte `db:"jsn"`
}

func main() {
	log.Printf("sarted")
	var db = sqlx.MustConnect(driverName, dsn)
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("close db connetcion: %v", err)
		} else {
			log.Printf("db connection closed")
		}
	}()
	var wg = &sync.WaitGroup{}
	for table, column := range tables {
		wg.Add(1)
		go downloadAndParseUi(wg, db, table, column)
	}
	wg.Wait()
	log.Printf("finisyed")
}

func downloadAndParseUi(
	wg *sync.WaitGroup,
	db *sqlx.DB,
	table string,
	column string,
) {
	defer wg.Done()
	var source []UiEntity
	var query = "select id as id, " + column + " as jsn from " + table
	var err = db.Select(&source, query)
	if err != nil {
		log.Printf("querying %s error: %v\n", table, err)
		return
	}
	err = parseJson(table, source)
	if err != nil {
		log.Printf("parsing %s JSON error: %v\n", table, err)
	}
	log.Printf("parsing %s finished\n", table)
}

func parseJson(
	table string,
	source []UiEntity,
) error {
	f, err := os.Create(table + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	i := 0
	b := iterjson.NewBuilder("", "  ")
	b.Add("", iterjson.TokenObjectOpen)
	b.Add("pairs", iterjson.TokenArrayOpen)
	for _, data := range source {
		for item, err := range iterjson.Parse(data.Json) {
			if err != nil {
				return err
			}
			var path = item.GetPathString()
			var key = item.Key
			var token = item.Token
			if skipByPath(path) || key.IsZero() || token.Type() != iterjson.TokenString {
				continue
			}
			if vs, err := token.GetString(); err != nil {
				return err
			} else if skipByValue(vs) {
				continue
			} else {
				i++
				b.Add("", iterjson.TokenObjectOpen)
				var k = fmt.Sprintf("%d.%s", data.Id, path)
				b.Add(k, vs)
				b.Add("", iterjson.TokenObjectClose)

			}
		}
	}
	b.Add("", iterjson.TokenArrayClose)
	b.Add("", iterjson.TokenObjectClose)
	out, err := b.Bytes()
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	fmt.Printf("%s: %d\n", table, i)
	return err
}

func skipByPath(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, "type") ||
		strings.HasSuffix(path, "coloroption") ||
		strings.HasSuffix(path, "color") ||
		!strings.Contains(path, "text") &&
			!strings.Contains(path, "title") &&
			!strings.Contains(path, "description") &&
			!strings.Contains(path, "caption")
}

func skipByValue(value string) bool {
	return strings.HasPrefix(value, colorPattern) ||
		strings.Contains(value, urlPattern) ||
		strings.TrimSpace(value) == ""
}
