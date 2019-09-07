// Copyright 2019 John Ott
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consumer

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/trivago/gollum/core"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

// IMU I2C Consumer
//
// This consumer read IMU values for the LSM9DS1 9DOF IMU over an I2C bus.
//
// Parameters
//
// - Bus: Defines the I2C Bus to connect to.
//
// - AccelerometerAddress: The accelerometer/gyroscope address on the I2C bus.
//
// - MagnetometerAddress: The magnetometer address on the I2C bus.
//
// Examples
//
// This config reads data from I2C2 on default addresses
//
//  IMUIn:
//    Type: consumer.Imu
//    Streams: imu
//	  Bus: I2C2
type Imu struct {
	core.SimpleConsumer `gollumdoc:"embed_type"`
	busName             string `config:"Bus"`
	accelAddr           string `config:"AccelerometerAddress" default:"0x6A`
	magnetoAddr         string `config:"MagnetometerAddress" default:"0x1C`
	bus                 i2c.BusCloser
	accel               *i2c.Dev
	magneto             *i2c.Dev
}

func init() {
	core.TypeRegistry.Register(Imu{})
}

func readReg(dev *i2c.Dev, reg byte, length byte) ([]byte, error) {
	tx := []byte{reg}
	r := make([]byte, length)
	err := dev.Tx(tx, r)
	return r, err
}

func writeReg(dev *i2c.Dev, reg byte, value []byte) error {
	tx := []byte{reg}
	tx = append(tx, value...)
	err := dev.Tx(tx, nil)
	return err
}

func (cons *Imu) Configure(conf core.PluginConfigReader) {
	var err error

	// Initialize periph
	if _, err := host.Init(); err != nil {
		cons.Logger.Error(err)
	}

	// Open I2C Bus
	if cons.bus, err = i2creg.Open(cons.busName); err != nil {
		cons.Logger.Error(err)
	}

	// Parse out the accelerometer address
	accelAddrInt, err := strconv.ParseUint(cons.accelAddr, 0, 16)
	if err != nil {
		cons.Logger.Error(err)
	}
	cons.accel = &i2c.Dev{Bus: cons.bus, Addr: uint16(accelAddrInt)}

	// Parse out the magnetometer address
	magnetoAddrInt, err := strconv.ParseUint(cons.magnetoAddr, 0, 16)
	if err != nil {
		cons.Logger.Error(err)
	}
	cons.magneto = &i2c.Dev{Bus: cons.bus, Addr: uint16(magnetoAddrInt)}

	// Do some configuration for the accelerometer and magnetometer here!

	//TODO: ODR and Scale
	cons.initGyro()
	cons.initAccel()
	cons.initMag()

	// For example:
	//err = writeReg(cons.accel, 0x10, []byte{0x20})
	//if err != nil {
	//	cons.Logger.Error(err)
	//}
}

func (cons *Imu) pollAccel() {
	for cons.IsActive() {
		// Poll the accelerometer status register for data and call cons.Enqueue()

		r, err := readReg(cons.accel, 0x28, 6)
		if err != nil {
			cons.Logger.Error(err)
		}

		valuex := (r[1] << 8) | r[0]
		valuey := (r[3] << 8) | r[2]
		valuez := (r[5] << 8) | r[4]

		// Enqueue new value
		str := fmt.Sprintf("{\"X\":%d}\n", valuex)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"Y\":%d}\n", valuey)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"Z\":%d}\n", valuez)
		cons.Enqueue([]byte(str))

		// Don't spam
		time.Sleep(100 * time.Millisecond)
	}
}

func (cons *Imu) Consume(workers *sync.WaitGroup) {
	// Close I2C bus on exit
	defer cons.bus.Close()

	go cons.pollAccel()

	cons.ControlLoop()
}

func (cons *Imu) initGyro() {
	var err error
	err = writeReg(cons.accel, 0x10, []byte{0xC0})//CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x11, []byte{0x00})//CTRL_REG2_G
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x12, []byte{0x00})//CTRL_REG3_G
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x1E, []byte{0x38})//CTRL_REG4_G
	if err != nil {
		cons.Logger.Error(err)
	}
	//Calibration
	r, err := readReg(cons.accel, 0x10, 1)//CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	var gScl=1
	r[0]&=byte (0xFF^(0x3<<3))
	r[0]|=byte (gScl<<3)
	err = writeReg(cons.accel, 0x10, r)//CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
}
func (cons *Imu) initAccel() {
	var err error
	err = writeReg(cons.accel, 0x1F, []byte{0x38})//CTRL_REG5_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x20, []byte{0x00})//CTRL_REG6_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x21, []byte{0x00})//CTRL_REG6_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	//Calibration
	r, err := readReg(cons.accel, 0x20, 1)//CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	var aScl=1
	r[0]&=byte (0xC7)
	r[0]|=byte (aScl<<3)
	err = writeReg(cons.accel, 0x20, r)//CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
}
func (cons *Imu) initMag() {
	var err error
	err = writeReg(cons.magneto, 0x20, []byte{0x1C})//CTRL_REG1_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x21, []byte{0x00})//CTRL_REG2_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x22, []byte{0x00})//CTRL_REG3_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x23, []byte{0x00})//CTRL_REG4_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x24, []byte{0x00})//CTRL_REG5_M
	if err != nil {
		cons.Logger.Error(err)
	}
}
