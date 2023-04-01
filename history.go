package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS content (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tags TEXT,
		text TEXT,
		created_at DATETIME
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func storeHistoryItem(db *sql.DB, tags, text string) error {
	fmt.Printf("Storing history item with tags: %s", tags)
	fmt.Printf("Storing history item with text: %s", text)
	stmt, err := db.Prepare("INSERT INTO content(tags, text, created_at) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tags, text, time.Now())
	return err
}

func retrieveHistoryItems(db *sql.DB, tags string) ([]string, []string, error) {
	last10Entries := []string{}
	taggedEntries := []string{}

	tagList := strings.Split(tags, ",")

	rows, err := db.Query("SELECT created_at, text FROM content ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var createdAt string
		var text string
		err = rows.Scan(&createdAt, &text)
		if err != nil {
			return nil, nil, err
		}
		last10Entries = append(last10Entries, fmt.Sprintf("%s: %s", createdAt, text))
	}

	if tags == "" {
		return last10Entries, taggedEntries, nil
	}

	sqlQuery := "SELECT created_at, text FROM content WHERE "
	for i, tag := range tagList {
		if i > 0 {
			sqlQuery += " OR "
		}
		sqlQuery += fmt.Sprintf("tags LIKE '%%%s%%'", tag)
	}
	sqlQuery += " ORDER BY created_at DESC LIMIT 10"

	rows, err = db.Query(sqlQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var createdAt string
		var text string
		err = rows.Scan(&createdAt, &text)
		if err != nil {
			return nil, nil, err
		}
		taggedEntries = append(taggedEntries, fmt.Sprintf("%s: %s", createdAt, text))
	}

	return last10Entries, taggedEntries, nil
}

//func main() {
//	db, err := initDB()
//	if err != nil {
//		fmt.Println("Error initializing database:", err)
//		return
//	}
//	defer db.Close()
//
//	err = storeContent(db, "tag1,tag2", "Sample text")
//	if err != nil {
//		fmt.Println("Error storing content:", err)
//		return
//	}
//
//	last10, tagged, err := retrieveContent(db, "tag1,tag2")
//	if err != nil {
//		fmt.Println("Error retrieving content:", err)
//		return
//	}
//
//	fmt.Println("Last 10 entries:")
//	for _, entry := range last10 {
//		fmt.Println(entry)
//	}
//
//	fmt.Println("Tagged entries:")
//	for _, entry := range tagged {
//		fmt.Println(entry)
//	}
//}
//
