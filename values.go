package main

import (
	"crypto/rand"
	"fmt"
	"github.com/capnm/sysinfo"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"time"
)

//-----------------------------------------------
//                     DATA
//-----------------------------------------------
//Read and store data in a local database
func dataCare() {
	//CPU
	cpuUsage()
	//RAM
	ramUsage()
	//Random Values - sensors
	sens = append(sens, SensorValues{values: getSensorVals()})
	clValues.sensors = sens
}

// Get dummy sensors values and save them
func getSensorVals() []int64 {
	sensorVals := make([]int64, 4)

	sensorVals[0], sensorVals[1], sensorVals[2], sensorVals[3] = getRandVals()

	return sensorVals
}

// Generate random Integers for dummy sensors
func getRandVals() (int64, int64, int64, int64) {

	big1 := big.NewInt(100)
	big2 := big.NewInt(150)
	big3 := big.NewInt(50)

	val1, err := rand.Int(rand.Reader, big1)
	if err != nil {
		fmt.Println(err)
	}

	val2, err := rand.Int(rand.Reader, big2)
	if err != nil {
		fmt.Println(err)
	}

	val3, err := rand.Int(rand.Reader, big3)
	if err != nil {
		fmt.Println(err)
	}

	val4, err := rand.Int(rand.Reader, big1)
	if err != nil {
		fmt.Println(err)
	}

	return val1.Int64(), val2.Int64(), (val3.Int64() - 30), val4.Int64()
}

// Get the ticks of the CPU
func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

// Calculate the usage of CPU
func cpuUsage() {
	idle0, total0 := getCPUSample()
	time.Sleep(250 * time.Millisecond)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	clValues.cpu = append(clValues.cpu, cpuUsage)
	//fmt.Printf("CPU usage is %.2f%% [busy: %.2f, total: %.2f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
}

// Calculate the amount of free RAM
func ramUsage() {
	si := sysinfo.Get()
	freeRamPer := float64(100 * (float64(si.FreeRam) / float64(si.TotalRam)))

	//usedRam := si.TotalRam - si.FreeRam
	//uint64(math.Round(freeRamPer*100)/100)
	clValues.ram = append(clValues.ram, freeRamPer)
	//fmt.Printf("Total RAM (MB): %d | Used RAM (MB): %d | Free RAM (MB): %d, (%%) %.2f |\n", (si.TotalRam / 1000), (usedRam / 1000), (si.FreeRam / 1000), freeRamPer)
}

//---------------------------------------------------
//                    PRINTS
//---------------------------------------------------
func printValues(n int, tipo int, vr int) {
	switch tipo {
	case 1: //Print for all variables
		dataCare()
		if n <= len(clValues.sensors) {
			printAllValues(n)
		} else {
			fmt.Println("Don't have enough values to print.")
		}
		break
	case 2: //Print just for one variable
		dataCare()
		if n <= len(clValues.sensors) {
			printSingleValue(n, vr)
		} else {
			fmt.Println("Don't have enough values to print.")
		}
		break
	case 3: //Print for more then one variable
		dataCare()
		if n <= len(clValues.sensors) {
			printMeanValue(vr)
		} else {
			fmt.Println("Don't have enough values to print.")
		}
		break
	default:
		break
	}
}

func printAllValues(n int) {
	fmt.Println("\n| CPU(%) | RAM(%) | Sensor1 | Sensor2 | Sensor3 | Sensor4 |")
	for i, _ := range clValues.sensors {
		fmt.Printf("|  %2.1f  |  %3.2f  |  %05d  |  %05d  |  %05d  |  %05d  |\n", clValues.cpu[i], clValues.ram[i], clValues.sensors[i].values[0], clValues.sensors[i].values[1], clValues.sensors[i].values[2], clValues.sensors[i].values[3])

		if i == n {
			fmt.Println()
			break
		}
	}
}

func printMeanValue(vr int) {
	total := 0.0
	var tot int64
	switch vr {
	case 1:
		fmt.Println("\nThe average of CPU usage.")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			total = total + clValues.cpu[i]
		}
		avgCpu := total / float64(len(clValues.cpu))
		fmt.Printf("Average of CPU usage %2.1f %%\n", avgCpu)
		break
	case 2:
		fmt.Printf("\nThe average of free RAM.\n")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			total = total + float64(clValues.ram[i])
		}
		avgRam := total / float64(len(clValues.cpu))
		fmt.Printf("Average of free RAM %.1f %%\n", avgRam)
		break
	case 3:
		fmt.Println("\nThe average of Sensor 1 values.")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			tot = tot + clValues.sensors[i].values[0]
		}
		avgSens1 := tot / int64(len(clValues.cpu))
		fmt.Printf("Sensor 1: %d\n", avgSens1)
		break
	case 4:
		fmt.Println("\nThe average of Sensor 2 values.")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			total = total + float64(clValues.sensors[i].values[1])
		}
		avgSens2 := total / float64(len(clValues.cpu))
		fmt.Printf("Average of Sensor 2: %.1f\n", avgSens2)
		break
	case 5:
		fmt.Println("\nThe average of Sensor 3 values.")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			total = total + float64(clValues.sensors[i].values[2])
		}
		avgSens3 := total / float64(len(clValues.cpu))
		fmt.Printf("Average of Sensor 3: %.1f\n", avgSens3)
		break
	case 6:
		fmt.Println("\nThe average of Sensor 4 values.")
		for i := 0; i <= len(clValues.cpu)-1; i++ {
			total = total + float64(clValues.sensors[i].values[3])
		}
		avgSens4 := total / float64(len(clValues.cpu))
		fmt.Printf("Average of Sensor 4: %.1f\n", avgSens4)
		break
	}
}

func printSingleValue(n int, vr int) {
	switch vr {
	case 1:
		fmt.Printf("\nThe Last %d CPU values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %2.1f%% |\n", clValues.cpu[i])
		}
		break
	case 2:
		fmt.Printf("\nThe Last %d RAM values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %3.2f%% |\n", clValues.ram[i])
		}
		break
	case 3:
		fmt.Printf("\nThe Last %d Sensor 1 values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %4d |\n", clValues.sensors[i].values[0])
		}
		break
	case 4:
		fmt.Printf("\nThe Last %d Sensor 2 values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %4d |\n", clValues.sensors[i].values[1])
		}
		break
	case 5:
		fmt.Printf("\nThe Last %d Sensor 3 values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %4d |\n", clValues.sensors[i].values[2])
		}
		break
	case 6:
		fmt.Printf("\nThe Last %d Sensor 4 values.\n", n)
		for i := 0; i <= n; i++ {
			fmt.Printf("| %4d |\n", clValues.sensors[i].values[3])
		}
		break
	}
}
