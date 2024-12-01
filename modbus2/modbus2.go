package modbus2

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/simonvetter/modbus"
)

type Conf struct {
	Name      []string
	Address   []int
	Size_data []string
	Bit       []int
	Type_data []string
}

type Res struct {
	Name []string
	Res  []int
}

type WriteReq struct {
	Register int
	DataSize string
	DataType string
	Value    int
}

// Création du client modbus
func CreateModbusClient(host string) (*modbus.ModbusClient, error) {

	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://" + host,
		Timeout: 1 * time.Second,
	})

	if err != nil {
		fmt.Printf("failed to create modbus client: %v\n", err)
		return nil, err
	}

	err = client.Open()
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		return nil, err
	}

	return client, nil
}

func (c *Conf) Read(mc *modbus.ModbusClient, r *Res) {

	r.Res = make([]int, len(c.Address))

	for i, j := range c.Address {
		dataSize := c.Size_data[i]
		dataSizeRead := 1
		dataType := c.Type_data[i]
		dataTypeRead := modbus.HOLDING_REGISTER
		bitNumber := c.Bit[i]

		switch dataSize {
		case "int16", "uint16":
			dataSizeRead = 1

		case "int32", "uint32":
			dataSizeRead = 2

		default:
			fmt.Println("Data size must be int16, uint16, int32 or uint32")
			os.Exit(1)
		}

		switch dataType {
		case "input":
			dataTypeRead = modbus.INPUT_REGISTER
		case "holding":
			dataTypeRead = modbus.HOLDING_REGISTER
		case "coil":
			regs, err := mc.ReadCoil(uint16(j))
			if err != nil {
				fmt.Printf("failed to read coil registers %d: %v\n", j, err)
				r.Res[i] = 404
			} else if regs == true {
				r.Res[i] = 1
			} else {
				r.Res[i] = 0
			}
		default:
			fmt.Println("Register Type must input, holding, or coil ")
			os.Exit(1)
		}

		if dataType != "coil" {
			regs, err := mc.ReadRegisters(uint16(j), uint16(dataSizeRead), dataTypeRead)
			if err != nil {
				fmt.Printf("failed to read registers %d: %v\n", j, err)
				r.Res[i] = 404

			} else if dataSizeRead == 2 {
				int32Value := int32(regs[0])<<16 | int32(uint16(regs[1]))

				if bitNumber < 32 {
					r.Res[i] = int(int32Value >> bitNumber & 1)
				} else {
					r.Res[i] = int(int32Value)
				}

			} else {
				if bitNumber < 32 {
					r.Res[i] = int(regs[0] >> bitNumber & 1)
				} else {
					r.Res[i] = int(regs[0])
				}
			}
		}
	}
	r.Name = c.Name
}

func (c *Conf) Decode(path string) {

	// il est probable que l'on puisse faire beaucoup plus simple en utilisant des modules de lecture csv pré existant.

	file, err := os.Open(path)
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

		//decode le nom
		c.Name = append(c.Name, record[0])

		//decode le registre

		if strings.HasPrefix(record[1], "0x") {

			k, err := strconv.ParseInt(record[1], 0, 0)
			if err != nil {
				fmt.Println("Erreur lors de la conversion du registre hexadécimal:", err)
			}
			c.Address = append(c.Address, int(k))

		} else {
			k, err := strconv.Atoi(record[1])
			if err != nil {
				fmt.Printf("failed to convert string to int %v\n", err)
			}
			c.Address = append(c.Address, k)
		}

		//decode la taille de la donnée requise
		c.Size_data = append(c.Size_data, record[2])

		//decode le bit si requis
		j, err := strconv.Atoi(record[3])
		if err != nil {
			fmt.Printf("failed to convert string to int %v\n", err)
		}
		c.Bit = append(c.Bit, j)

		//decode le type de la donnée requise (input, coil ou holding)
		data_type := record[4]

		switch data_type {
		case "coil", "input", "holding":
			c.Type_data = append(c.Type_data, data_type)
		default:
			fmt.Println("wrong input data, must be coil, input or holding")
			os.Exit(1)
		}
	}
}

func Write(c *modbus.ModbusClient, r WriteReq) {

	fmt.Println("write function called")
	var err error

	switch r.DataType {
	case "coil register":
		err = c.WriteCoil(uint16(r.Register), r.Value == 1)
	case "holding register":
		switch r.DataSize {
		case "int16", "uint16":
			err = c.WriteRegister(uint16(r.Register), uint16(r.Value))
		case "int32", "uint32":
			err = c.WriteUint32(uint16(r.Register), uint32(r.Value))
		default:
			fmt.Println("Unsupported DataSize for register:", r.DataSize)
			return
		}
	default:
		fmt.Println("Unsupported DataType:", r.DataType)
		return
	}

	if err != nil {
		fmt.Println("Error during write operation:", err)
	} else {
		fmt.Println("Write operation successful")
	}
}
