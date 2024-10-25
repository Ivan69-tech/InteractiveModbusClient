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

type Conf struct {
	Name      []string
	Address   []int
	Type_data []string
	Bit       []int
}

type Res struct {
	Name []string
	Res  []int
}

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

func (c *Conf) Read(mc *modbus.ModbusClient, r *Res) {

	r.Res = make([]int, len(c.Address))
	fmt.Print("ok")
	for i, j := range c.Address {
		regs, err := mc.ReadRegisters(uint16(j), 1, modbus.HOLDING_REGISTER)

		if err != nil {
			fmt.Printf("failed to read registers %d: %v\n", j, err)
		}
		fmt.Println(regs)
		r.Res[i] = int(regs[0])

	}

	r.Name = c.Name
}

func (c *Conf) Decode() {

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
		c.Name = append(c.Name, record[0])
		k, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Printf("failed to convert string to int %v\n", err)
		}
		c.Address = append(c.Address, k)
		c.Type_data = append(c.Type_data, record[2])

		j, err := strconv.Atoi(record[3])
		if err != nil {
			fmt.Printf("failed to convert string to int %v\n", err)
		}
		c.Bit = append(c.Bit, j)
	}
}
