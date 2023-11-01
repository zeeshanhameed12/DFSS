type Workflow struct {
    //SecFlag bool    // Security flag
    //hashOfFile string // Security level
    fileSize float32
    CPU      int // CPU utilization
    RAM      int // Memory utilization
    Location string
}

type listOfPeer struct {
    peer     string
    space    float32
    CPU      int
    RAM      int
    Location string
}

size := len(encodedData[0])
	wList := []Workflow{
		{fileSize: float32(size), CPU: ReqCPU, RAM: ReqRAM, Location: loc},
	}

	allResults := executePSO(wList, listOfPeer)