package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

func sendFiletoPeer() {
	// Connect to the IPFS API
	api := shell.NewShell("localhost:5001")
	var filePath string
	var ReqCPU, ReqRAM int
	var loc string
	fmt.Println("Enter full path of the file you want to save:")
	fmt.Scan(&filePath)
	fmt.Println("Enter Required RAM:")
	fmt.Scan(&ReqRAM)
	fmt.Println("Enter Required CPU:")
	fmt.Scan(&ReqCPU)
	fmt.Println("Enter the region e.g. Europe,Asia...:")
	fmt.Scan(&loc)
	// fmt.Printf("File size: %d bytes\n", fileSize)
	data, err := ioutil.ReadFile(string(filePath))
	if err != nil {
		panic(err)
	}
	NodesInfo := getFreeSpace()
	groupedNodes, err := splitNodesByLocation(NodesInfo)
	if err != nil {
		fmt.Println("Error splitting nodes by location:", err)
		return
	}
	// Encode file
	//for i := 0; i < 10; i++ {
	// fmt.Println(As_Nodes)
	// fmt.Println(Af_Nodes)

	encodingTime := time.Now()
	dataShards, parityShards := 2, 4
	encodedData, metadata, err := encodingData(data, dataShards, parityShards)
	fmt.Println("Encoding time",time.Since(encodingTime))
	metadata.fileName = filePath
	if err != nil {
		log.Fatal(err)
	}
	size := len(encodedData[0])
	wList := []Workflow{
		{fileSize: float32(size), CPU: ReqCPU, RAM: ReqRAM, Location: loc},
	}
	
		
	
	
	allResults := executePSO(wList, groupedNodes)
	
	
	dist := time.Now()

	var homeRegionNodes []NodeListWithCost
	var homeRegionIndex int
	var found bool

	// Find the home region
	for i, region := range allResults {
		if len(region) > 0 && region[0].location == loc {
			homeRegionNodes = region
			homeRegionIndex = i
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Home region %s not found", loc)
		return
	}
	totalShards := dataShards + parityShards
	fmt.Println("Required number of nodes:", totalShards)
	fmt.Println("Nodes in home region:", len(homeRegionNodes))

	if len(homeRegionNodes) >= totalShards {
		// If the home region has enough nodes
		distributeData(encodedData, homeRegionNodes, metadata, api)
		fmt.Println("distribution time",time.Since(dist))
		return
	}

	requiredAdditionalNodes := totalShards - len(homeRegionNodes)
	allRegionsCount := len(allResults)
	selectedNodes := homeRegionNodes

	for offset := 1; requiredAdditionalNodes > 0; offset++ {
		nextRegionIndex := (homeRegionIndex + offset) % allRegionsCount
		prevRegionIndex := (homeRegionIndex - offset + allRegionsCount) % allRegionsCount

		if nextRegionIndex != homeRegionIndex && len(allResults[nextRegionIndex]) > 0 {
			nodesToTake := min(requiredAdditionalNodes, len(allResults[nextRegionIndex]))
			selectedNodes = append(selectedNodes, allResults[nextRegionIndex][:nodesToTake]...)
			requiredAdditionalNodes -= nodesToTake
			fmt.Printf("%v Chunks will be stored in %v regions\n", nodesToTake, allResults[nextRegionIndex][0].location)
		}

		if requiredAdditionalNodes > 0 && prevRegionIndex != homeRegionIndex && len(allResults[prevRegionIndex]) > 0 {
			nodesToTake := min(requiredAdditionalNodes, len(allResults[prevRegionIndex]))
			selectedNodes = append(selectedNodes, allResults[prevRegionIndex][:nodesToTake]...)
			requiredAdditionalNodes -= nodesToTake
			fmt.Printf("%v Chunks will be stored in %v regions\n", nodesToTake, allResults[prevRegionIndex][0].location)
		}

		if nextRegionIndex == homeRegionIndex && prevRegionIndex == homeRegionIndex {
			log.Println("Not enough nodes available even in the neighboring regions")
			return
		}
	}

	distributeData(encodedData, selectedNodes, metadata, api)
	fmt.Println("distribution time",time.Since(dist))

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//nodeMat := pso(wList, NodesInfo)
//fmt.Println("ppropriate nodes found in all regions are:", allResults)
// fmt.Println("ppropriate nodes found in this region are:", len(allResults))
// fmt.Println("required number of nodes are:", dataShards+parityShards)
// //fmt.Println("len of All_Nodes_With_Group", groupedNodes)
// fmt.Println("All_Nodes_With_Group", len(groupedNodes))
// //fmt.Println(nodeMat)
// var homeRegion struct {
// 	nodesInHomeRegion []NodeListWithCost
// 	index             int
// }

// for j, region := range allResults {
// 	if region[j].location == loc {
// 		homeRegion.nodesInHomeRegion = region
// 		homeRegion.index = j
// 	}

// }

// fmt.Println("required number of nodes are:", dataShards+parityShards)
// fmt.Println("Nodes in home region:", len(homeRegion.nodesInHomeRegion))
// if len(allResults[homeRegion.index]) >= (dataShards + parityShards) {
// 	distributeData(encodedData, homeRegion.nodesInHomeRegion, metadata, api)
// } else {
// 	reqAdditionalNodes := ((dataShards + parityShards) - len(allResults[homeRegion.index]))
// 	//fmt.Println("required additional nodes:", len(allResults[(homeRegion.index)+1]))
// 	if len(allResults[(homeRegion.index)+1]) >= reqAdditionalNodes {
// 		taking_Nodes := allResults[(homeRegion.index)+1][:reqAdditionalNodes]
// 		homeRegion.nodesInHomeRegion = append(homeRegion.nodesInHomeRegion, taking_Nodes...)
// 		distributeData(encodedData, homeRegion.nodesInHomeRegion, metadata, api)
// 		fmt.Printf("%v Chunks will be stored in %v regions\n", reqAdditionalNodes, allResults[(homeRegion.index)+1][0].location)

// 	} else if len(allResults[(homeRegion.index)-1]) >= reqAdditionalNodes {
// 		taking_Node := allResults[(homeRegion.index)-1][:reqAdditionalNodes]
// 		homeRegion.nodesInHomeRegion = append(homeRegion.nodesInHomeRegion, taking_Node...)
// 		distributeData(encodedData, homeRegion.nodesInHomeRegion, metadata, api)
// 		fmt.Printf("%v Chunks will be stored in %v regions\n", reqAdditionalNodes, allResults[(homeRegion.index)-1][0].location)

// 	}

// }

//fmt.Printf("Nodes in %v and their costs to store the %v file\n", loc, filePath)
// fmt.Println("Nodes in your region and ranking based on cost function")
// for i, dt := range nodeMat[0] {

// 	fmt.Println(i, dt)
// }

//fmt.Println("top peer ", nodeMat[0][0].peer)

// fmt.Printf("List of posible solutions for task 1:%v\nList of posible solutions for task 2:%v\nList of posible solutions for task 3:%v\nList of posible solutions for task 4:%v\n", nodeMat[0], nodeMat[1], nodeMat[2], nodeMat[3])
// var wg sync.WaitGroup
// for i, data := range encodedData {
// 	wg.Add(1)
// 	go func(data []byte, i int) {
// 		defer wg.Done()
// 		buffer := bytes.NewBuffer(data)
// 		//start := time.Now()
// 		fileHash, err := api.Add(buffer)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		topPeer := allResults[1][i].peer
// 		response, err := apiHandle(fileHash, topPeer)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer response.Body.Close()

// 		metadata.IpfsHashes[i] = fileHash

// 		// fmt.Println("########Sending########")
// 		// fmt.Printf("Storage Node:%v\nNode Region: %v\nChunk hash:%v\nChunk size:%v\nData Transfer time(Elapsed time for sending file to suitable Node):%v\n", topPeer, loc, fileHash, size, time.Since(start))
// 		// fmt.Println("########Sent########")
// 	}(data, i)
// }
// wg.Wait()

// SaveMetadata(metadata)

// fmt.Println(time.Since(encodingTime))
//time.Sleep(2 * time.Second)
//}

// func apiHandle(fileHash string, topPeer string) (*http.Response, error) {
// 	addresses := []string{
// 		"http://127.0.0.1:9094",
// 		"http://127.0.0.1:9098",
// 		"http://127.0.0.1:9099",
// 	}

// 	var lastError error
// 	for _, address := range addresses {
// 		url := fmt.Sprintf("%s/pins/ipfs/%s?mode=recursive&name=&replication-max=1&replication-min=1&shard-size=0&user-allocations=%s", address, fileHash, topPeer)
// 		response, err := http.Post(url, "application/json", nil)
// 		if err != nil {
// 			lastError = fmt.Errorf("API %s is not responding: %w", address, err)
// 			continue
// 		}
// 		return response, nil
// 	}
// 	return nil, lastError
// }

// func splitNodesByLocation(nodes []listOfPeer) ([][]listOfPeer, error) {
// 	if len(nodes) == 0 {
// 		return nil, fmt.Errorf("no nodes available")
// 	}

// 	var As_Nodes, Af_Nodes, Eu_Nodes []listOfPeer

// 	for _, node := range nodes {
// 		switch node.Location {
// 		case "Asia":
// 			As_Nodes = append(As_Nodes, node)
// 		case "Africa":
// 			Af_Nodes = append(Af_Nodes, node)
// 		case "Europe": // Added "Europe" assuming Eu_Nodes means European nodes
// 			Eu_Nodes = append(Eu_Nodes, node)
// 		default:
// 			return nil, fmt.Errorf("unknown location: %s", node.Location)
// 		}
// 	}

// 	return [][]listOfPeer{As_Nodes, Af_Nodes, Eu_Nodes}, nil
// }

// func executePSO(wList []Workflow, nodes [][]listOfPeer) [][]NodeListWithCost {
// 	var wg sync.WaitGroup
// 	results := make(chan [][]NodeListWithCost, len(nodes))

// 	for _, nodeGroup := range nodes {
// 		wg.Add(1)

// 		go func(workflows []Workflow, peers []listOfPeer) {
// 			defer wg.Done()
// 			nodeMat := pso(workflows, peers)
// 			for p := range nodeMat {
// 				sort.Slice(nodeMat[p], func(i, j int) bool {
// 					return nodeMat[p][i].cost < nodeMat[p][j].cost
// 				})
// 			}
// 			results <- nodeMat
// 		}(wList, nodeGroup)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	var allResults [][]NodeListWithCost
// 	for result := range results {
// 		allResults = append(allResults, result...)
// 	}

// 	return allResults
// }
