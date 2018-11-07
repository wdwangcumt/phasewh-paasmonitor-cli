/**
 * Author: Wang Weidong (Hxdi)
 * Created Date: 2018-01-13
 * Project Description: paasmonitor-cli
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego/httplib"
	//"log"
	"os"
	//"os/exec"
	//"os/user"
	"sort"
	"strconv"
	"time"
)

const (
	LOOP_INTERVAL time.Duration = 1
)

type ClusterMonitorListResponse struct {
	DateTime           string `json:"DateTime"`
	Count              int64
	ClusterMonitorList []ClusterMonitor `json:"ClusterMonitorList"`
}

type ClusterMonitor struct {
	Host               string  `orm:"pk" json:"Host"`
	HostType           string  `json:"HostType"`
	HostDescription    string  `json:"HostDescription"`
	Status             string  `json:"Status"`
	CpuCores           int     `json:"CpuCores"`
	CpuUsagePercent    float64 `json:"CpuUsagePercent"`
	MemoryTotal        string  `json:"MemoryTotal"`
	MemoryUsed         string  `json:"MemoryUsed"`
	MemoryFree         string  `json:"MemoryFree"`
	MemoryUsagePercent float64 `json:"MemoryUsagePercent"`
	DiskSize           string  `json:"DiskSize"`
	DiskUsed           string  `json:"DiskUsed"`
	DiskFree           string  `json:"DiskFree"`
	DiskUsagePercent   float64 `json:"DiskUsagePercent"`
	Processes          string  `json:"Processes"`
}

var (
	paasmonitor = os.Getenv("PAAS_MONITOR_SERVER")
	cleanScreen = []byte{27, 91, 72, 27, 91, 50, 74}
)

func main() {
	var (
		arg string
	)

	if paasmonitor == "" {
		fmt.Println("Please set environment variable 'PAAS_MONITOR_SERVER'!")
		return
	}

	flag.Parse()
	arg = flag.Arg(0)
	if arg != "help" && arg != "cluster" && arg != "" {
		arg = "host"
	}
	//fmt.Println("arg is :", arg)

	for {
		//clear the console
		//	cmd := exec.Command("clear")
		//	clearConsole, err := cmd.Output()
		//	if err != nil {
		//		fmt.Println(err.Error())
		//	}
		//	fmt.Println(clearConsole)
		//	fmt.Print(string(clearConsole))
		//	cmd.Wait()
		fmt.Print(string(cleanScreen))

		switch arg {
		case "cluster":
			getCluster()
		case "host":
			getHostDetails(flag.Arg(0))
		case "help":
			fmt.Println("Usage of paasmonitor-cli:")
			fmt.Println("    cluster : View cluster information in real time.")
			fmt.Println("     [host] : Host IP or name, eg: 10.20.30.40, view the specified host information in real time.")
			return
		default:
			fmt.Println("Try 'paasmonitor-cli help' for more information.")
			return
		}

		//sleep
		time.Sleep(LOOP_INTERVAL * time.Second)
	}

}

/**
 * Get cluster information
 */
func getCluster() {
	var (
		pme      string = paasmonitor + "/pm/cluster/information"
		response ClusterMonitorListResponse
		err      error
	)

	request := httplib.Get(pme)
	httpResponse, err := request.DoRequest()
	if err != nil {
		fmt.Println("Cluster monitor endpoint '"+pme+"' is unreachable!", err.Error())
		os.Exit(1)
	}

	request.ToJSON(&response)
	err = httpResponse.Body.Close()
	if err != nil {
		fmt.Println("Http response error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("***** PaasMonitor-CLI *****\n")
	fmt.Printf("%-18s  %-10s  %-25s  %-7s  %-10s  %-20s  %-10s  %-20s  %-10s\n", "Host", "Type", "Description", "Status", "CPU(%)", "Memory Usage/Limit", "Memory(%)", "Disk Usage/Total", "Disk(%)")
	for i := 0; i < int(response.Count); i++ {
		cm := response.ClusterMonitorList[i]
		fmt.Printf("%-18s  %-10s  %-25s  %-7s  %-10s  %-20s  %-10s  %-20s  %-10s\n",
			cm.Host,
			cm.HostType,
			cm.HostDescription,
			cm.Status,
			strconv.FormatFloat(cm.CpuUsagePercent, 'f', 3, 64)+"%",
			cm.MemoryUsed+"/"+cm.MemoryTotal,
			strconv.FormatFloat(cm.MemoryUsagePercent, 'f', 3, 64)+"%",
			cm.DiskUsed+"/"+cm.DiskSize,
			strconv.FormatFloat(cm.DiskUsagePercent, 'f', 3, 64)+"%")
	}
}

/**
 * Get a host's details
 */
func getHostDetails(host string) {
	var (
		pme          string = paasmonitor + "/pm/cluster/" + host + "/information"
		response     ClusterMonitor
		processesMap map[string]string
		err          error
	)
	processesMap = make(map[string]string)

	request := httplib.Get(pme)
	httpResponse, err := request.DoRequest()
	if err != nil {
		fmt.Println("Cluster monitor endpoint '"+pme+"' is unreachable!", err.Error())
		os.Exit(1)
	}
	request.ToJSON(&response)
	err = httpResponse.Body.Close()
	if err != nil {
		fmt.Println("Http response error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("***** PaasMonitor-CLI *****\n")
	fmt.Println("Host:", response.Host)
	fmt.Println("HostType:", response.HostType)
	fmt.Println("Status:", response.Status)
	fmt.Println("Description:", response.HostDescription)

	fmt.Println("CPU:")
	fmt.Printf("    %20s  %-10d\n", "Cores:", response.CpuCores)
	fmt.Printf("    %20s  %-10s\n", "CPUUsage:", strconv.FormatFloat(response.CpuUsagePercent, 'f', 3, 64)+"%")

	fmt.Println("Memory:")
	fmt.Printf("    %20s  %-20s\n", "Memory Usage/Limit:", response.MemoryUsed+"/"+response.MemoryTotal)
	fmt.Printf("    %20s  %-20s\n", "MemoryUsage:", strconv.FormatFloat(response.MemoryUsagePercent, 'f', 3, 64)+"%")

	fmt.Println("Disk:")
	fmt.Printf("    %20s  %-20s\n", "Disk Usage/Total:", response.DiskUsed+"/"+response.DiskSize)
	fmt.Printf("    %20s  %-20s\n", "DiskUsage:", strconv.FormatFloat(response.DiskUsagePercent, 'f', 3, 64)+"%")

	err = json.Unmarshal([]byte(response.Processes), &processesMap)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Process:")
	sorted_keys := make([]string, 0)
	for k, _ := range processesMap {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	for _, v := range sorted_keys {
		fmt.Printf("    %20s  %-10s\n", v, processesMap[v])
	}
}
