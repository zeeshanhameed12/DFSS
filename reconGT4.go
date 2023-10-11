package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

type ChunkResult struct {
	Index int
	Data  []byte
	Valid bool
}

func getFile() {
	fmt.Println("Provide the name of the file you want to retrieve:")
	var fileName string
	fmt.Scan(&fileName)

	metadata, err := LoadMetadata(fileName)
	if err != nil {
		panic(err)
	}

	retreiveddata := make([][]byte, len(metadata.IpfsHashes))
	fmt.Println("########Retreiving########")
	api := shell.NewShell("localhost:5001")

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan bool)
		startt := time.Now()
		go func() {
			results, validChunks := downloadChunks(ctx, cancel, metadata, api, done)
			if validChunks < metadata.DataShards {
				fmt.Println("Not enough valid chunks to reconstruct the file")
				return
			}
			sort.Slice(results, func(i, j int) bool {
				return results[i].Index < results[j].Index
			})

			for i, chunk := range results {
				retreiveddata[i] = chunk.Data
			}

			_, err := decodingData(retreiveddata, metadata)
			if err != nil {
				panic(err)
			}

			// Save the decoded file
			// err = ioutil.WriteFile("RetrievedFile.txt", decodedData, 0644)
			// if err != nil {
			// 	panic(err)
			// }
			// fmt.Println("File retrieved successfully!")
			// close(done)
		}()

		<-done // This will block until the file is retrieved or an error occurs
		elapsed := time.Since(startt)
		fmt.Println(elapsed)
		time.Sleep(2 * time.Second)
	}
}

func downloadChunks(ctx context.Context, cancel context.CancelFunc, metadata *Metadata, api *shell.Shell, done chan bool) ([]ChunkResult, int) {
	results := make([]ChunkResult, len(metadata.IpfsHashes))

	validChunks := 0
	mux := &sync.Mutex{}

	for i := range metadata.IpfsHashes {
		go func(i int) {
			defer func() {
				if validChunks == metadata.DataShards {
					done <- true
				}
			}()

			select {
			case <-ctx.Done():
				return
			default:
			}

			fileContents, err := api.Cat(metadata.IpfsHashes[i])
			if err != nil {
				results[i] = ChunkResult{Index: i}
				return
			}

			data, err := io.ReadAll(fileContents)
			if err != nil {
				results[i] = ChunkResult{Index: i}
				return
			}

			hash := md5.Sum(data)
			if fmt.Sprintf("%x", hash) == metadata.ShardHashes[i] {
				results[i] = ChunkResult{Index: i, Data: data, Valid: true}
				//fmt.Printf("Chunk %v retrieved successfully!\n", i)

				mux.Lock()
				validChunks++
				if validChunks == metadata.DataShards {
					cancel()
				}
				mux.Unlock()
			} else {
				results[i] = ChunkResult{Index: i}
			}
		}(i)
	}

	<-done // Wait for the done signal
	return results, validChunks
}
