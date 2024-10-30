package main

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "assets.db"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomDomain() string {
	return faker.DomainName()
}

func randomWord() string {
	return faker.Word()
}

func randomComment() string {
	return faker.Sentence()
}

func randomIP() string {
	ip := make(net.IP, 4)
	rand.Read(ip)
	return ip.String()
}

func randomOwner() string {
	return faker.Name()
}

func randomPort() int {
	return rand.Intn(65535-1) + 1
}

func main() {
	if _, err := os.Stat(dbFileName); err == nil {
		os.Remove(dbFileName)
	}

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tables with indexes for faster querying
	createTablesSQL := `
		CREATE TABLE assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host TEXT NOT NULL,
			comment TEXT,
			owner TEXT NOT NULL
		);

		CREATE TABLE ips (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER,
			address TEXT NOT NULL,
			FOREIGN KEY(asset_id) REFERENCES assets(id)
		);

		CREATE TABLE ports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER,
			port INTEGER NOT NULL,
			FOREIGN KEY(asset_id) REFERENCES assets(id)
		);

		CREATE INDEX idx_asset_host ON assets(host);
		CREATE INDEX idx_ip_address ON ips(address);
		CREATE INDEX idx_port_number ON ports(port);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare statements for inserts
	insertAssetStmt, err := db.Prepare("INSERT INTO assets(host, comment, owner) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertAssetStmt.Close()

	insertIPStmt, err := db.Prepare("INSERT INTO ips(asset_id, address) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertIPStmt.Close()

	insertPortStmt, err := db.Prepare("INSERT INTO ports(asset_id, port) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertPortStmt.Close()

	// Insert data: 100.000 rows
	totalEntries := 10000
	for i := 0; i < totalEntries; i++ {
		// Insert asset
		host := fmt.Sprintf("%s.%s", randomWord(), randomDomain())
		comment := randomComment()
		owner := randomOwner()
		result, err := insertAssetStmt.Exec(host, comment, owner)
		if err != nil {
			log.Fatal(err)
		}

		assetID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		// Insert related IPs
		for j := 0; j < rand.Intn(3)+1; j++ { // Each asset has 1-3 IPs
			ip := randomIP()
			_, err := insertIPStmt.Exec(assetID, ip)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Insert related Ports
		for j := 0; j < rand.Intn(5)+1; j++ { // Each asset has 1-5 ports
			port := randomPort()
			_, err := insertPortStmt.Exec(assetID, port)
			if err != nil {
				log.Fatal(err)
			}
		}

		if i%1000 == 0 {
			log.Printf("Inserted %d/%d entries", i, totalEntries)
		}
	}

	log.Println("Database created and populated successfully")
}
