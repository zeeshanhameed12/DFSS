package main

import (
	"fmt"
	//"math/rand"
	//"os/exec"
)

// type ReconFile struct {
// 	index     int
// 	chunkData []byte
// 	flg       int
// }

type Metadata struct {
	fileName     string
	FileSize     int64
	DataShards   int
	ParityShards int
	ShardHashes  []string
	ShardOrder   []int
	IpfsHashes   []string
}
type Workflow struct {
	//SecFlag bool    // Security flag
	//hashOfFile string // Security level
	fileSize float32
	CPU      int // CPU utilization
	RAM      int // Memory utilization
	Location string

	//Best int // Best position
}

// structure for unmarshelling the the data for getting free available space frome nodes
type peers struct {
	Name          string
	Peer          string
	Value         string
	Expire        int64
	Valid         bool
	Weight        int64
	Partitionable bool
	Received_at   int64
}

type listOfPeer struct {
	peer     string
	space    float32
	CPU      int
	RAM      int
	Location string
}

// structure for unmarshelling the data for statcheck function

// type statusCheck struct {
// 	CID string `json:"cid"`
// 	Al []string `json:"allocations"`
// }

var PeerData = make([]listOfPeer, 0)
// var activeroutines = 0
// var validChunks = 0
// var breakLoop = false

//var fsp = make([]int64, 0)

func switchMenu() {

	var input int
	fmt.Print("Please select the available options below.\n 1: For the peer IDs \n 2: Nodes and their ranking based on available free space \n 3: For adding files to a top ranked node \n 4: For retreiving files\n")
	fmt.Scan(&input)
	//fmt.Println(input)
	switch input {
	case 1:
		getPostID()
	case 2:
		getFreeSpace()
	case 3:
		sendFiletoPeer()
	case 4:

		getFile()
	default:
		fmt.Println("Invalid command")
	}
}

func main() {

	for {
		switchMenu()
	}
	//sendFiletoPeer()
	// // define the command to run
	// cmd := exec.Command("./cluster-ctl","id")

	// // run the command and capture the output
	// output, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// // print the output
	// fmt.Println(string(output))

}

// func saveFile(fileName string, data string) error {
// 	return ioutil.WriteFile(fileName, []byte(data), 0666)
// }
// func print(d []byte) { // we can call the function using object created by the type dc
// 	for i, card := range d {
// 		fmt.Println(i, card)
// 	}
// }

// func sendFiletoPeer() {
// 	cmd2 := exec.Command("ls")

// 	// Run the command and wait for it to complete
// 	output2, err2 := cmd2.Output()
// 	if err2 != nil {
// 		fmt.Println("Error running command:", err2)
// 		return
// 	}

// 	// Print the command output
// 	fmt.Println(string(output2))

// 	cmd := exec.Command("cluster-ctl --enc=json --debug  pin add --replication 1 --allocations 12D3KooWFphNrPi8EBPRVgLnQUEAYHwC53RcR2LkcqNnh5brYsZ3 QmWG6fFWjddYnQbLaSsAYaAr8i5xFc4QDiixTM8QXuwi4j")

// 	// Run the command and wait for it to complete
// 	output, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println("Error running command:", err)
// 		return
// 	}

// 	// Print the command output
// 	fmt.Println(string(output))

// }
