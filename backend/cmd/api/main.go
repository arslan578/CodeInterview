package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "assets.db"
const defaultLimit = 10 // Show 10 rows per page

// Asset represents an asset with associated IPs and Ports
type Asset struct {
	ID        int
	Host      string
	Comment   string
	Owner     string
	IPs       []IP
	Ports     []Port
	Signature string
}

// IP represents an IP address associated with an asset
type IP struct {
	Address   string
	Signature string
}

// Port represents a port associated with an asset
type Port struct {
	Port      int
	Signature string
}

// generateSignature creates a unique SHA-256 hash for the asset, its IPs, and Ports
func generateSignature(asset Asset) Asset {
	data := asset.Host + asset.Comment + asset.Owner
	hash := sha256.New()
	hash.Write([]byte(data))
	signature := hex.EncodeToString(hash.Sum(nil))

	newAsset := asset
	newAsset.Signature = signature

	// Generate signatures for IPs
	for i, ip := range newAsset.IPs {
		ipData := ip.Address
		hash := sha256.New()
		hash.Write([]byte(ipData))
		ip.Signature = hex.EncodeToString(hash.Sum(nil))
		newAsset.IPs[i] = ip
	}

	// Generate signatures for Ports
	for i, port := range newAsset.Ports {
		portData := fmt.Sprintf("%d", port.Port)
		hash := sha256.New()
		hash.Write([]byte(portData))
		port.Signature = hex.EncodeToString(hash.Sum(nil))
		newAsset.Ports[i] = port
	}

	return newAsset
}

func main() {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize Gin router and set up CORS
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"X-Total-Count"},
		AllowCredentials: true,
	}))

	// Handler to fetch assets with pagination and optional search by host
	router.GET("/assets", func(c *gin.Context) {
		// Pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", fmt.Sprintf("%d", defaultLimit)))
		offset := (page - 1) * limit

		//fmt.Printf("Limit: %d, Offset: %d\n", limit, offset) // Debug logging

		assetID := c.Query("id")
		hostFilter := c.Query("host")
		var rows *sql.Rows

		// Count total assets for pagination
		var totalCount int
		countQuery := "SELECT COUNT(*) FROM assets WHERE 1=1"
		if hostFilter != "" {
			countQuery += " AND host LIKE ?"
			err := db.QueryRow(countQuery, "%"+hostFilter+"%").Scan(&totalCount)
			if err != nil {
				//fmt.Printf("Error in Count Query: %v\n", err) // Log detailed error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in count query"})
				return
			}
		} else {
			err := db.QueryRow(countQuery).Scan(&totalCount)
			if err != nil {
				//fmt.Printf("Error in Count Query: %v\n", err) // Log detailed error
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in count query"})
				return
			}
		}

		//fmt.Printf("Total count of records: %d\n", totalCount) // Debug logging

		// Fetch assets with pagination and IPs/Ports grouped in subqueries
		query := `
    SELECT a.id, a.host, a.comment, a.owner,
           (SELECT GROUP_CONCAT(ip.address) FROM ips ip WHERE ip.asset_id = a.id) AS ip_addresses,
           (SELECT GROUP_CONCAT(p.port) FROM ports p WHERE p.asset_id = a.id) AS port_numbers
    FROM assets a
    WHERE 1=1
`
		if assetID != "" {
			query += " AND a.id = ?"
			rows, err = db.Query(query, assetID)
		} else if hostFilter != "" {
			query += " AND a.host LIKE ? LIMIT ? OFFSET ?"
			rows, err = db.Query(query, "%"+hostFilter+"%", limit, offset)
		} else {
			query += " LIMIT ? OFFSET ?"
			rows, err = db.Query(query, limit, offset)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var assets []Asset
		for rows.Next() {
			var id int
			var host, comment, owner, ipAddresses, portNumbers string
			if err := rows.Scan(&id, &host, &comment, &owner, &ipAddresses, &portNumbers); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Split the IPs and Ports and build the asset structure
			ips := []IP{}
			ports := []Port{}

			for _, ip := range strings.Split(ipAddresses, ",") {
				ips = append(ips, IP{Address: ip})
			}

			for _, port := range strings.Split(portNumbers, ",") {
				portInt, _ := strconv.Atoi(port)
				ports = append(ports, Port{Port: portInt})
			}

			asset := Asset{
				ID:      id,
				Host:    host,
				Comment: comment,
				Owner:   owner,
				IPs:     ips,
				Ports:   ports,
			}

			// Generate the signature for each asset
			processedAsset := generateSignature(asset)
			assets = append(assets, processedAsset)
		}

		// Set the X-Total-Count header for pagination
		c.Header("X-Total-Count", fmt.Sprintf("%d", totalCount))
		// Return the paginated assets
		c.JSON(http.StatusOK, assets)
	})

	// Run the Gin server on port 8080
	router.Run(":8080")
}
