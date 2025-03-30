package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cacheDir := flag.String("cache-dir", "", "Path to the directory containing the game client cache.")
	blobDB := flag.String("db", "", "Path to the target blob.db database.")
	everything := flag.Bool("everything", false, "Do not skip any files. Optional, default: false.")

	flag.Parse()

	if *cacheDir == "" || *blobDB == "" {
		flag.PrintDefaults()
		return
	}

	if err := importCacheDir(*cacheDir, *blobDB, *everything); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func importCacheDir(cacheDir string, blobDB string, everything bool) error {
	db, err := sql.Open("sqlite3", blobDB)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	dirEntries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	var inserted int64

	for _, dirEntry := range dirEntries {
		cdnid, ok := getEntryCDNID(dirEntry, everything)
		if !ok {
			continue
		}

		blobPath := filepath.Join(cacheDir, dirEntry.Name())
		blob, err := os.ReadFile(blobPath)
		if err != nil {
			return err
		}

		hashSum := sha1.Sum(blob)
		blobHash := hex.EncodeToString(hashSum[:])

		affected, err := insertBlobToDB(tx, cdnid, blob, blobHash)
		if err != nil {
			return fmt.Errorf("cdnid: %s, err: %w", cdnid, err)
		}

		inserted += affected
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Done:")
	fmt.Printf("  total files: %d \n", len(dirEntries))
	fmt.Printf("  new files: %d \n", inserted)
	fmt.Printf("  ignored files: %d \n", int64(len(dirEntries))-inserted)

	return nil
}

func getEntryCDNID(dirEntry os.DirEntry, everything bool) (string, bool) {
	if dirEntry.IsDir() {
		fmt.Printf("skipping directory %s \n", dirEntry.Name())
		return "", false
	}

	if dirEntry.Name() == ".DS_Store" {
		fmt.Printf("skipping file %s (you do not need this) \n", dirEntry.Name())
		return "", false
	}

	if everything {
		return dirEntry.Name(), true
	}

	// trim extension
	fileName := dirEntry.Name()
	fileExt := filepath.Ext(fileName)
	cdnid := fileName[:len(fileName)-len(fileExt)]

	if len(cdnid) != 18 {
		// skip index.txt and other non-asset files
		fmt.Printf("skipping file %s (not an asset) \n", dirEntry.Name())
		return "", false
	}

	// if cdnid != fileName {
	// 	fmt.Println(fileName, "->", cdnid)
	// }

	return cdnid, true
}

func insertBlobToDB(tx *sql.Tx, cdnid string, blob []byte, blobHash string) (int64, error) {
	var dbHash string
	row := tx.QueryRow("select hash from asset_file where cdnid = ?;", cdnid)
	if err := row.Scan(&dbHash); err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if dbHash == blobHash {
		return 0, nil
	} else if dbHash != "" {
		return 0, fmt.Errorf("hashes do not match")
	}

	res, err := tx.Exec("insert into asset_file (cdnid, blob, hash) values (?, ?, ?);", cdnid, blob, blobHash)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
