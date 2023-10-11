package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

func getFreeSpace() []listOfPeer {

	response,
		err := http.Get("http://127.0.0.1:9094/monitor/metrics/freespace")
	if err != nil {
		fmt.Println("API ttp://127.0.0.1:9094 is not responding\nSwitching to ttp://127.0.0.1:9098")
		// os.Exit(1)
		response,
			err := http.Get("http://127.0.0.1:9098/monitor/metrics/freespace")
		if err != nil {
			fmt.Println("API ttp://127.0.0.1:9098 is not responding\nSwitching to ttp://127.0.0.1:9099")
			response,
				err := http.Get("http://127.0.0.1:9099/monitor/metrics/freespace")
			if err != nil {
				fmt.Println("API ttp://127.0.0.1:9099 is not responding")
				os.Exit(1)
			}
			data := getInfo(response)
			return data
		}
		dtat := getInfo(response)
		return dtat
	}
	data := getInfo(response)
	return data
}

func getInfo(response *http.Response) []listOfPeer {
	var data = make([]listOfPeer, 0)
	bytes,
		err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	prs := []peers{}

	json.Unmarshal(bytes, &prs)

	//fmt.Println("It is working")
	for _, info := range prs {
		var peerlist = listOfPeer{

			peer:  info.Peer,
			space: float32(info.Weight),

			//CPU:   rand.Intn(10) + 1,
		}

		//fmt.Println(peerlist)

		// Arbitarary assign CPU to nodes

		if peerlist.peer == "12D3KooWBLipb6inEBvJxUo3gAaDrS1trVZHzwgetN9ZvYc4Cmvv" {
			peerlist.CPU = 9
			peerlist.RAM = 128
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWBTqj5VBqJj7Vq1oKULigDieWwkm4V9BzFMgpx5eomWPs" {
			peerlist.CPU = 8

			peerlist.RAM = 16
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWC8Mz4D8JZospKscbkdS7SrzH28p9TLL1s78tg4Py4jdo" {
			peerlist.CPU = 8

			peerlist.RAM = 32
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWDc2VfXXwq37ozZZJUkiQSw6cN31xoFkKp8H8PPqS2VDm" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWFUHPdKo4F9upQMN1vawfPP3iFxZvv8wFvAUBVfwGuncw" {
			peerlist.CPU = 7
			peerlist.RAM = 128
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWHAuHTH2iNrPEYnAS3HE6fY1H22D99nfE5Vg37XbWoatw" {
			peerlist.CPU = 9
			peerlist.RAM = 16
			peerlist.Location = "Asia"

		}
		if peerlist.peer == "12D3KooWJvdiBS1c5s66hrM7Uhm2bcUKGyXvtNjmC6hWcrj2cQeL" {
			peerlist.CPU = 9
			peerlist.RAM = 256
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWL93L4yUu5gQJJgVadx9ZjyWFXeN3qPkuTMCc51fqs6GR" {
			peerlist.CPU = 8
			peerlist.RAM = 128
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWLabiNXQ3K2CCTjo4GWUsQ6zsYCfAo3HD22zGWKQNmMc4" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWMBdBfbt1WofCvan7W4DzqWF4FhLNYLWsci2tCXaZQ2xQ" {
			peerlist.CPU = 11
			peerlist.RAM = 256
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWMJByo59nsjioH7wFZYczQq4QRdSdgVhbAjJtvcQ8hLUz" {
			peerlist.CPU = 11
			peerlist.RAM = 120
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWMZd75CY76mrv3dwyj3QTDa7MnrR1CmwdwHh4paTY9gV7" {
			peerlist.CPU = 11
			peerlist.RAM = 64
			peerlist.Location = "Europe"

		}
		if peerlist.peer == "12D3KooWMfbVmNdZuNubD2wsczxD9Dn7vL6bqTTWu9yJ5ZZmrtb6" {
			peerlist.CPU = 11
			peerlist.RAM = 80
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWNVcwuNdV49dBVYos8L41eBdRkesxGtuXeRLc1dPoeYXs" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWNqR83enQHokNikRq1ynUgFgbUcYWTnhWaQ19jAmHw5TH" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWPTVcLBjZH2RS8TYm393E1kxR39DJYuLWxebtpTFyV6Z1" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWPmM7PHSobUrBCV9wXbn4UPUL4Bj1HuyMf9LFRmNnCYUH" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}
		if peerlist.peer == "12D3KooWSy6eZ1v7YDUbgWbfwhMhdaLKa6szKUUrS8RrpATFH2hT" {
			peerlist.CPU = 11
			peerlist.RAM = 110
			peerlist.Location = "Africa"

		}

		data = append(data, peerlist)
		//fmt.Println("All the Nodes with availabale free disk space")

	}

	// fmt.Println("All the Nodes with availabale free disk space")
	// for i, dt := range data {
	// 	fmt.Println(i, dt)
	// }

	sort.Slice(data, func(i, j int) bool {
		return data[i].space > data[j].space // use ">" if you want descending order
	})
	// for _, fs := range PeerData {

	// 	fsp = append(fsp, fs.space)
	// }

	//fmt.Println("lenght of string ",string(responseDat))
	//ioutil.WriteFile()
	//print(responseDat)

	return data
}
