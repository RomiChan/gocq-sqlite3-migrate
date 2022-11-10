package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/RomiChan/gocq-sqlite3-migrate/db"
	"github.com/RomiChan/gocq-sqlite3-migrate/leveldb"
	"github.com/RomiChan/gocq-sqlite3-migrate/sqlite3"
)

func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	v3 := flag.String("from", "data/leveldb-v3", "leveldb v3 path")
	sql3 := flag.String("to", "data/sqlite3/msg.db", "sqlite3 db path")
	flag.Parse()
	if flag.Arg(0) == "help" {
		flag.Usage()
		return
	}
	v3db, err := leveldb.Open(*v3)
	handle(err)
	defer v3db.Close()
	sql3db, err := sqlite3.Open(*sql3)
	handle(err)
	i := 0
	errs := v3db.ForEach(func(x any) error {
		fmt.Printf("\rprocess: %d", i)
		i++
		switch msg := x.(type) {
		case *db.StoredGroupMessage:
			return sql3db.InsertGroupMessage(msg)
		case *db.StoredPrivateMessage:
			return sql3db.InsertPrivateMessage(msg)
		case *db.StoredGuildChannelMessage:
			return sql3db.InsertGuildChannelMessage(msg)
		}
		return nil
	})
	fmt.Println(" succeeded.")
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}
	handle(sql3db.Close())
}
