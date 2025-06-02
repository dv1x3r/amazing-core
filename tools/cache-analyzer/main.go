package main

import (
	"bytes"
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
	fmt.Printf("cdnid=%s type=%s version=%s \n", cdnid, detectFileType(blob), "")
	return nil
}

func hasPrefix(blob, prefix []byte) bool {
	return len(blob) >= len(prefix) && bytes.Equal(blob[:len(prefix)], prefix)
}

func detectFileType(blob []byte) string {
	switch {
	case hasPrefix(blob, []byte{0x7b}):
		return "json"
	case hasPrefix(blob, []byte{0x4F, 0x67, 0x67, 0x53, 0x00}):
		return "audio/ogg"
	case hasPrefix(blob, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00}):
		return "image/png"
	case hasPrefix(blob, []byte{0x55, 0x6E, 0x69, 0x74, 0x79, 0x57, 0x65, 0x62, 0x00}):
		return "UnityWeb"
	case hasPrefix(blob, []byte{0x55, 0x6E, 0x69, 0x74, 0x79, 0x46, 0x53, 0x00}):
		return "UnityFS"
	case hasPrefix(blob, []byte{0x51, 0x75, 0x65, 0x73, 0x74, 0x0D, 0x0A, 0x7B}):
		return "TreeNode/Quest"
	default:
		return ""
	}
}
