package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
	//"fmt"
	//"os"
	//"path/filepath"
)

func getPostID() {
	//var response *http.Request
	response,
		err := http.Get("http://127.0.0.1:9094/id")
	if err != nil {
		fmt.Println("API ttp://127.0.0.1:9094 is not responding\nSwitching to ttp://127.0.0.1:9098")
		// os.Exit(1)
		response,
			err := http.Get("http://127.0.0.1:9098/id")
		if err != nil {
			fmt.Println("API ttp://127.0.0.1:9098 is not responding\nSwitching to ttp://127.0.0.1:9099")
			response,
				err := http.Get("http://127.0.0.1:9099/id")
			if err != nil {
				fmt.Println("API ttp://127.0.0.1:9099 is not responding")
				os.Exit(1)
			}
			doThis(response)
			return
		}
		doThis(response)
		return
	}
	doThis(response)
}
func doThis(response *http.Response) {
	responseDat,
		err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(string(responseDat))
}

// func statCheck() statusCheck{
// 	response,
// 		err := http.Get("http://127.0.0.1:9094/pins/QmX6vwLvrRft4MAmc9K2aqRdpFQ1QF3C4iqg9PFGtZkL3n?local=false")

// 		//err := http.Get("http://127.0.0.1:9094/monitor/metrics/freespace")
// 		//err := http.Get("http://127.0.0.1:9094/pins/QmX6vwLvrRft4MAmc9K2aqRdpFQ1QF3C4iqg9PFGtZkL3n?local=false")
// 	if err != nil {

// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}
// 	responseDat,
// 		err := io.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}

// 	//fmt.Println(string(responseDat))
// 	status := statusCheck{}
// 	json.Unmarshal(responseDat, &status)
// 	//fmt.Println(status)
//  return status
// }

func SaveMetadata(metadata *Metadata) error {
	// Connect to PostgreSQL database
	dbURL := "postgres://zeeshan:postgress@localhost:5432/udimyproject"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS metadata (
		fileName TEXT,
        filesize INT, 
        dataShards INT, 
        parityShards INT,
        shardHashes TEXT,
        shardOrder TEXT,
		ipfsHashes TEXT,
        CONSTRAINT metadata_pkey PRIMARY KEY (fileName)
    )`)
	if err != nil {
		fmt.Print(err)
		return err
	}

	// Insert metadata
	shardHashes := strings.Join(metadata.ShardHashes, ",")
	ipfsHashes := strings.Join(metadata.IpfsHashes, ",")
	var shardOrderStrings []string
	for _, order := range metadata.ShardOrder {
		shardOrderStrings = append(shardOrderStrings, strconv.Itoa(order))
	}
	shardOrder := strings.Join(shardOrderStrings, ",")
	_, err = db.Exec(`INSERT INTO metadata (fileName, filesize, dataShards, parityShards, shardHashes, shardOrder, ipfsHashes) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (fileName) DO UPDATE
        SET shardHashes = $5, shardOrder = $6, ipfsHashes = $7,fileSize = $2, dataShards = $3, parityShards = $4`,
		metadata.fileName, metadata.FileSize, metadata.DataShards, metadata.ParityShards, shardHashes, shardOrder, ipfsHashes)
	if err != nil {
		return err
	}

	return nil
}

func LoadMetadata(file string) (*Metadata, error) {
	// Connect to PostgreSQL database
	dbURL := "postgres://zeeshan:postgress@localhost:5432/udimyproject"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	// Query for metadata
	row := db.QueryRow("SELECT * FROM metadata WHERE filename = $1", file)
	var (
		fileName     string
		filesize     int
		dataShards   int
		parityShards int
		shardHashes  string
		ipfsHashes   string
		shardOrder   string
	)
	if err := row.Scan(&fileName, &filesize, &dataShards, &parityShards, &shardHashes, &shardOrder, &ipfsHashes); err != nil {
		return nil, err
	}

	// Parse strings into slices
	shardHashesSlice := strings.Split(shardHashes, ",")
	shardOrderSlice := strings.Split(shardOrder, ",")
	ipfsHashe := strings.Split(ipfsHashes, ",")
	shardOrderInts := make([]int, len(shardOrderSlice))
	for i, s := range shardOrderSlice {
		shardOrderInts[i], _ = strconv.Atoi(s)
	}

	// Construct Metadata object
	metadata := &Metadata{
		fileName:     fileName,
		FileSize:     int64(filesize),
		DataShards:   dataShards,
		ParityShards: parityShards,
		IpfsHashes:   ipfsHashe,
		ShardHashes:  shardHashesSlice,
		ShardOrder:   shardOrderInts,
	}

	return metadata, nil
}

func apiHandle(fileHash string, topPeer string) (*http.Response, error) {
	addresses := []string{
		"http://127.0.0.1:9094",
		"http://127.0.0.1:9098",
		"http://127.0.0.1:9099",
	}

	var lastError error
	for _, address := range addresses {
		url := fmt.Sprintf("%s/pins/ipfs/%s?mode=recursive&name=&replication-max=1&replication-min=1&shard-size=0&user-allocations=%s", address, fileHash, topPeer)
		response, err := http.Post(url, "application/json", nil)
		if err != nil {
			lastError = fmt.Errorf("API %s is not responding: %w", address, err)
			continue
		}
		return response, nil
	}
	return nil, lastError
}

// func addByteStreamToIPFS(byteStream []byte) (string, error) {
// 	// Run the ipfs add command with the byte stream as input
// 	cmd := exec.Command("ipfs", "add", "--quiet")
// 	cmd.Stdin = bytes.NewReader(byteStream)
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return "", err
// 	}

// 	// Extract the CID from the command output
// 	cid := string(bytes.TrimSpace(output))

// 	return cid, nil
// }

// func getFileSize(path string) (int64, error) {
// 	absPath, err := filepath.Abs(path)
// 	if err != nil {
// 		return 0, err
// 	}

// 	fileInfo, err := os.Stat(absPath)
// 	if err != nil {
// 		return 0, err
// 	}

//		return fileInfo.Size(), nil
//	}
func distributeData(encodedData [][]byte, allResults []NodeListWithCost, metadata *Metadata, api *shell.Shell) {
	var wg sync.WaitGroup
	for i, data := range encodedData {
		wg.Add(1)
		go func(data []byte, i int) {
			defer wg.Done()

			buffer := bytes.NewBuffer(data)
			//start := time.Now()
			fileHash, err := api.Add(buffer)
			if err != nil {
				log.Fatal(err)
			}

			topPeer := allResults[i].peer
			response, err := apiHandle(fileHash, topPeer)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			metadata.IpfsHashes[i] = fileHash

			// fmt.Println("########Sending########")
			// fmt.Printf("Storage Node:%v\nNode Region: %v\nChunk hash:%v\nChunk size:%v\nData Transfer time(Elapsed time for sending file to suitable Node):%v\n", topPeer, loc, fileHash, size, time.Since(start))
			// fmt.Println("########Sent########")
		}(data, i)
	}
	wg.Wait()

	SaveMetadata(metadata)

	//fmt.Println(time.Since(encodingTime)) // encodingTime needs to be defined or passed to this function as needed
}
