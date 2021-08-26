package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 帧类型

const (
	IPv4   uint8 = 0x01
	DIRECT uint8 = 0x02 // 浏览器直接发给远程的类型
	DOMAIN uint8 = 0x03
	IPv6   uint8 = 0x04

	HEAD uint16 = 0x1201
	TAIL uint16 = 0x0825
)

type PDU struct {
	Head   uint16
	Order  uint32
	Type   uint8
	Length uint16
	Data   []byte
	Tail   uint16
}

func NewPDU(order uint32, type_ uint8, data []byte) *PDU {
	return &PDU{
		Head:   HEAD,
		Order:  order,
		Type:   type_,
		Length: uint16(len(data)),
		Data:   data,
		Tail:   TAIL,
	}
}

func (p *PDU) Encode() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.BigEndian, p.Head); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, p.Order); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, p.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, p.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, p.Data); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.BigEndian, p.Tail); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (p *PDU) Decode(bts []byte) error {
	if len(bts) < 10 {
		return fmt.Errorf("pdu is too short")
	}
	dataReader := bytes.NewReader(bts)

	if err := binary.Read(dataReader, binary.BigEndian, &p.Head); err != nil {
		return err
	}
	if p.Head != HEAD {
		return fmt.Errorf("error HEAD")
	}
	if err := binary.Read(dataReader, binary.BigEndian, &p.Order); err != nil {
		return err
	}
	if err := binary.Read(dataReader, binary.BigEndian, &p.Type); err != nil {
		return err
	}
	if err := binary.Read(dataReader, binary.BigEndian, &p.Length); err != nil {
		return err
	}
	if len(bts) != 11+int(p.Length) {
		return fmt.Errorf("error pdu length %v", 10+int(p.Length))
	}
	p.Data = make([]byte, p.Length)
	if err := binary.Read(dataReader, binary.BigEndian, &p.Data); err != nil {
		return err
	}
	if err := binary.Read(dataReader, binary.BigEndian, &p.Tail); err != nil {
		return err
	}
	if p.Tail != TAIL {
		return fmt.Errorf("error TAIL")
	}

	return nil
}
