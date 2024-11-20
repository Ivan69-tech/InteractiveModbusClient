package database

import (
	"gotools2/modbus2"
	"time"
)

type Database struct {
	Time    []time.Time
	Signaux []string
	Data    [][]int
}

func (d *Database) Save(m modbus2.Res) {

	d.Time = append(d.Time, time.Now())
	d.Signaux = m.Name

	if len(d.Data) == 0 {
		d.Data = make([][]int, len(d.Signaux))
	}
	for k := 0; k < len(m.Name); k++ {
		d.Data[k] = append(d.Data[k], m.Res[k])
	}

}
