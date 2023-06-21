// Copyright (c) 2017-2022 The lxi developers. All rights reserved.
// Project site: https://github.com/twsnmp/lxi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package lxi

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// Device models an LXI device, which is currently just a TCPIP socket
// interface. An LXI Device also implements the ivi.Driver interface.
type Device struct {
	conn    net.Conn
	timeout int
}

// NewDevice opens a TCPIP Device using the given VISA address resource string.
// timeout in milisecond.
func NewDevice(address string, timeout int) (*Device, error) {
	var d Device
	v, err := NewVisaResource(address)
	if err != nil {
		return &d, err
	}
	tcpAddress := fmt.Sprintf("%s:%d", v.hostAddress, v.port)
	c, err := net.Dial("tcp", tcpAddress)
	if err != nil {
		return &d, err
	}
	d.conn = c
	d.timeout = timeout
	return &d, nil
}

// Write writes the given data to the network connection.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.conn.Write(p)
}

// Read reads from the network connection into the given byte slice.
func (d *Device) Read(p []byte) (n int, err error) {
	if d.timeout > 0 {
		d.conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(d.timeout)))
	} else {
		d.conn.SetReadDeadline(time.Time{})
	}
	return d.conn.Read(p)
}

// Close closes the underlying network connection.
func (d *Device) Close() error {
	return d.conn.Close()
}

// WriteString writes a string using the underlying network connection.
func (d *Device) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Command sends the SCPI/ASCII command to the underlying network connection. A
// newline character is automatically added to the end of the string.
func (d *Device) Command(format string, a ...interface{}) error {
	cmd := format
	if a != nil {
		cmd = fmt.Sprintf(format, a...)
	}
	_, err := d.WriteString(strings.TrimSpace(cmd) + "\n")
	return err
}

// Query writes the given string to the underlying network connection and
// returns a string. A newline character is automatically added to the query
// command sent to the instrument.
func (d *Device) Query(cmd string) (string, error) {
	if cmd != "" {
		if err := d.Command(cmd); err != nil {
			return "", err
		}
	}
	if d.timeout > 0 {
		d.conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(d.timeout)))
	} else {
		d.conn.SetReadDeadline(time.Time{})
	}
	return bufio.NewReader(d.conn).ReadString('\n')
}

// SetTimeout set timeout in milisecond.
func (d *Device) SetTimeout(timeout int) {
	d.timeout = timeout
}
