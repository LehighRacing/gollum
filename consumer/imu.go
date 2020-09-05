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
	accelAddr           string `config:"AccelerometerAddress" default:"0x6B"`
	magnetoAddr         string `config:"MagnetometerAddress" default:"0x1E"`
	bus                 i2c.BusCloser
	accel               *i2c.Dev
	gRes                float32
	aRes                float32
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
	//fmt.Printf("open I2C bus\n");
	if cons.bus, err = i2creg.Open("/dev/i2c-2" /*cons.busName*/); err != nil {
		cons.Logger.Error(err)
	}

	// Parse out the accelerometer address
	//fmt.Printf("parse XL address from:%s\n",cons.accelAddr);
	accelAddrInt, err := strconv.ParseUint(cons.accelAddr, 0, 16)
	if err != nil {
		cons.Logger.Error(err)
	}
	//fmt.Printf("accel=%x\n",accelAddrInt);
	cons.accel = &i2c.Dev{Bus: cons.bus, Addr: uint16(accelAddrInt)}

	// Parse out the magnetometer address
	//fmt.Printf("parse magnet address\n");
	magnetoAddrInt, err := strconv.ParseUint(cons.magnetoAddr, 0, 16)
	if err != nil {
		cons.Logger.Error(err)
	}
	cons.magneto = &i2c.Dev{Bus: cons.bus, Addr: uint16(magnetoAddrInt)}

	// Do some configuration for the accelerometer and magnetometer here!

	cons.initGyro()
	cons.initAccel()
	//must add magnetometer address (0x1E) to device tree, otherwise errors
	//cons.initMag()
}

//OUT_X_XL=0x28-0x29
//OUT_Y_XL=0x2A-0x2B
//OUT_Z_XL=0x2C-0x2D
func (cons *Imu) pollAccel() {
	for cons.IsActive() {
		// Poll the accelerometer status register for data and call cons.Enqueue()

		r, err := readReg(cons.accel, 0x28, 6)
		if err != nil {
			cons.Logger.Error(err)
		}

		valuex := uint16(r[1]<<8) | uint16(r[0])
		valuey := uint16(r[3]<<8) | uint16(r[2])
		valuez := uint16(r[5]<<8) | uint16(r[4])
		fvaluex := cons.calcVal(valuex, cons.aRes)
		fvaluey := cons.calcVal(valuey, cons.aRes)
		fvaluez := cons.calcVal(valuez, cons.aRes)

		//Enqueue raw data
		str := fmt.Sprintf("{\"XL\":%d,%d,%d}\n", valuex, valuey, valuez)
		cons.Enqueue([]byte(str))
		//Enqueue calculated data
		str = fmt.Sprintf("{\"XLCalc\":%f,%f,%f}\n", fvaluex, fvaluey, fvaluez)
		cons.Enqueue([]byte(str))

		/*// Enqueue new value
		str := fmt.Sprintf("{\"X\":%d}\n", valuex)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"Xf\":%f}\n", fvaluex)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"Y\":%d}\n", valuey)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"Yf\":%f}\n", fvaluey)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"Z\":%d}\n", valuez)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"Zf\":%f}\n", fvaluez)
		cons.Enqueue([]byte(str))*/

		// Don't spam
		// collect data every 100 milliseconds (10Hz)
		time.Sleep(100 * time.Millisecond)
	}
}

func (cons *Imu) pollGyro() {
	for cons.IsActive() {
		// Poll the gyrocope status register for data and call cons.Enqueue()

		r, err := readReg(cons.accel, 0x18, 6)
		if err != nil {
			cons.Logger.Error(err)
		}

		valuex := uint16(r[1]<<8) | uint16(r[0])
		valuey := uint16(r[3]<<8) | uint16(r[2])
		valuez := uint16(r[5]<<8) | uint16(r[4])
		fvaluex := cons.calcVal(valuex, cons.gRes)
		fvaluey := cons.calcVal(valuey, cons.gRes)
		fvaluez := cons.calcVal(valuez, cons.gRes)

		//Enqueue raw data
		str := fmt.Sprintf("{\"G\":%d,%d,%d}\n", valuex, valuey, valuez)
		cons.Enqueue([]byte(str))
		//Enqueue calculated data
		str = fmt.Sprintf("{\"GCalc\":%f,%f,%f}\n", fvaluex, fvaluey, fvaluez)
		cons.Enqueue([]byte(str))

		/*// Enqueue new value
		str := fmt.Sprintf("{\"XG\":%d}\n", valuex)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"XGf\":%f}\n", fvaluex)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"YG\":%d}\n", valuey)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"YGf\":%f}\n", fvaluey)
		cons.Enqueue([]byte(str))
		// Enqueue new value
		str = fmt.Sprintf("{\"ZG\":%d}\n", valuez)
		cons.Enqueue([]byte(str))
		str = fmt.Sprintf("{\"ZGf\":%f}\n", fvaluez)
		cons.Enqueue([]byte(str))*/

		// Don't spam
		// collect data every 100 milliseconds (10Hz)
		time.Sleep(100 * time.Millisecond)
	}
}

func (cons *Imu) calcVal(value uint16, res float32) float32 {
	//fmt.Printf("%d\n",value)
	out := float32(value) * res
	return out
}

func (cons *Imu) Consume(workers *sync.WaitGroup) {
	// Close I2C bus on exit
	defer cons.bus.Close()

	go cons.pollAccel()
	go cons.pollGyro()

	cons.ControlLoop()
}

func (cons *Imu) initGyro() {
	fmt.Printf("initGyro\n")
	var err error
	err = writeReg(cons.accel, 0x10, []byte{0xC0}) //CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	cons.Logger.Error(err)
	err = writeReg(cons.accel, 0x11, []byte{0x00}) //CTRL_REG2_G
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x12, []byte{0x00}) //CTRL_REG3_G
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x1E, []byte{0x38}) //CTRL_REG4_G
	if err != nil {
		cons.Logger.Error(err)
	}
	//Set scale
	r, err := readReg(cons.accel, 0x10, 1) //CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	var gScl = 0
	r[0] &= byte(0xFF ^ (0x3 << 3))
	r[0] |= byte(gScl << 3)
	//Set ODR
	var gODR = 0x1 //14.9Hz
	r[0] &= byte(0xFF ^ (0x7 << 5))
	r[0] |= byte(gODR << 5)
	//Write both
	err = writeReg(cons.accel, 0x10, r) //CTRL_REG1_G
	if err != nil {
		cons.Logger.Error(err)
	}
	//calcgRes
	//00=245,01=500,11=2000
	cons.gRes = float32(245) / 32768.0
}
func (cons *Imu) initAccel() {
	fmt.Printf("initAccel\n")
	var err error
	err = writeReg(cons.accel, 0x1F, []byte{0x38}) //CTRL_REG5_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x20, []byte{0x20}) //CTRL_REG6_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.accel, 0x21, []byte{0x00}) //CTRL_REG7_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	//Set scale
	r, err := readReg(cons.accel, 0x20, 1) //CTRL_REG6_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	var aScl = 0 //+-2g
	r[0] &= byte(0xC7)
	r[0] |= byte(aScl << 3)
	//Set ODR
	var aODR = 0x1 //10Hz
	r[0] &= byte(0x1F)
	r[0] |= byte(aODR << 5)
	//Write both
	err = writeReg(cons.accel, 0x20, r) //CTRL_REG6_XL
	if err != nil {
		cons.Logger.Error(err)
	}
	//calcaRes
	//00=2,01=16,10=4,11=8
	cons.aRes = float32(2) / 32768.0
}
func (cons *Imu) initMag() {
	var err error
	err = writeReg(cons.magneto, 0x20, []byte{0x1C}) //CTRL_REG1_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x21, []byte{0x00}) //CTRL_REG2_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x22, []byte{0x00}) //CTRL_REG3_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x23, []byte{0x00}) //CTRL_REG4_M
	if err != nil {
		cons.Logger.Error(err)
	}
	err = writeReg(cons.magneto, 0x24, []byte{0x00}) //CTRL_REG5_M
	if err != nil {
		cons.Logger.Error(err)
	}
}
