package server

import (
	"fmt"
	"os"
	"time"

	"github.com/simonvetter/modbus"
)

var (
	K int = 0
)

func updateInput() {
	for {
		time.Sleep(1 * time.Second)
		K += 1
		fmt.Println(K)
	}
}

func Server() {
	go updateInput()
	var server *modbus.ModbusServer
	var mh *modbusHandler
	var err error

	mh = &modbusHandler{}
	mh.holdingReg[4] = 40
	mh.coils[10] = true

	host := "tcp://localhost:1502"

	server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL:        host,
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, mh)

	if err != nil {
		fmt.Printf("failed to create server: %v\n", err)
		os.Exit(1)
	}

	err = server.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("server started on : %s", host)

	for {
	}

}

type modbusHandler struct {
	coils      [100]bool
	holdingReg [5]uint16
}

func (mh *modbusHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {

	if req.UnitId != 1 {
		err = modbus.ErrIllegalFunction
		return
	}

	if int(req.Addr)+int(req.Quantity) > len(mh.coils) {
		err = modbus.ErrIllegalDataAddress
		return
	}
	for i := 0; i < int(req.Quantity); i++ {

		if req.IsWrite {
			mh.coils[int(req.Addr)+i] = req.Args[i]
		}
		res = append(res, mh.coils[int(req.Addr)+i])
	}
	return
}

func (mh *modbusHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction
	return
}

func (mh *modbusHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {

	var (
		regAddr uint16
	)

	for i := 0; i < int(req.Quantity); i++ {
		regAddr = req.Addr + uint16(i)
		switch regAddr {
		case 1000:
			if req.IsWrite {
				mh.holdingReg[4] = req.Args[i]
			}
			res = append(res, mh.holdingReg[4])
		case 1100:
			if req.IsWrite {
				mh.holdingReg[0] = req.Args[i]
			}
			res = append(res, mh.holdingReg[0])

		case 1101:
			if req.IsWrite {
				mh.holdingReg[1] = req.Args[i]
			}
			res = append(res, mh.holdingReg[1])

		case 1102:
			if req.IsWrite {
				mh.holdingReg[2] = req.Args[i]
			}
			res = append(res, mh.holdingReg[2])

		case 1103:
			if req.IsWrite {
				mh.holdingReg[3] = req.Args[i]
			}
			res = append(res, mh.holdingReg[3])

		default:
			res = append(res, 0)
		}
	}
	return
}

// loop:
// 	for i := 0; i < int(req.Quantity); i++ {
// 		regAddr = req.Addr + uint16(i)
// 		for j := 0; j < len(mh.holdingReg); j++ {
// 			//fmt.Println(regAddr, baseAddr)

// 			if regAddr == baseAddr {
// 				if req.IsWrite {
// 					mh.holdingReg[j] = req.Args[i]
// 				}
// 				res = append(res, mh.holdingReg[j])
// 				continue loop
// 			}
// 			baseAddr++
// 		}

// 		res = append(res, 0)
// 	}

func (mh *modbusHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	var regAddr uint16

	for i := 0; i < int(req.Quantity); i++ {
		regAddr = req.Addr + uint16(i)
		switch regAddr {
		case 150:
			err = modbus.ErrIllegalDataAddress
			return
		case 100:
			res = append(res, uint16(K))
		case 101:
			res = append(res, 15)
		default:
			res = append(res, 0)
			return
		}
	}
	return
}
