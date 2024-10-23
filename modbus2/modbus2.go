package modbus2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/simonvetter/modbus"
)

var (
	Variable1 uint16
	Variable2 uint16
)

func CreateModbusClient(adresse string, port string) *modbus.ModbusClient {

	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + adresse + ":" + port,
		Timeout: 1 * time.Second,
	})

	if err != nil {
		fmt.Printf("failed to create modbus client: %v\n", err)
		os.Exit(1)
	}

	err = client.Open()
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(2)
	}

	return client
}

func Read(mc *modbus.ModbusClient, reg map[string]int) map[string]int {

	res := make(map[string]int)
	for i, j := range reg {
		regs, err := mc.ReadRegisters(uint16(j), 1, modbus.HOLDING_REGISTER)

		if err != nil {
			fmt.Printf("failed to read registers %d: %v\n", j, err)
		}
		res[i] = int(regs[0])
	}
	return res
}

func Decode() map[string]int {

	res := make(map[string]int)

	file, err := os.Open("conf.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		//fmt.Println(record)
		k, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Printf("failed to convert string to int %v\n", err)
		}
		res[record[0]] = k
	}

	fmt.Println(res)

	return res
}
