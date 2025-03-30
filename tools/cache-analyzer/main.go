package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	blobDB := flag.String("db", "", "Path to the blob.db database.")

	flag.Parse()

	if *blobDB == "" {
		flag.PrintDefaults()
		return
	}

	if err := analyzeBlobDB(*blobDB); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func analyzeBlobDB(blobDB string) error {
	db, err := sql.Open("sqlite3", blobDB)
	if err != nil {
		return err
	}
	defer db.Close()

	var cdnid string
	var blob []byte

	rows, err := db.Query("select cdnid, blob from asset_file order by id;")
	for rows.Next() {
		if err := rows.Scan(&cdnid, &blob); err != nil {
			return err
		}

		if err := processRow(cdnid, blob); err != nil {
			return err
		}
	}

	return nil
}

func processRow(cdnid string, blob []byte) error {
	fmt.Println(cdnid)
	return nil
}
