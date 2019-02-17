// Copyright (c) 2019 Adrian P.K. <apk@kuguar.io>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

/**
 *Copyright (c) 2018 Adrian P.K. <apk@kuguar.io>
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package nulls

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"

	"fmt"
)

// Point - Spatial point
type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// MakeGeoPoint - Returns a valid
func MakeGeoPoint(lat, long float64, valid bool) GeoPoint {
	return GeoPoint{Point{lat, long}, valid}
}

func (p *Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)
}

// Scan implements the Scanner interface.
func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// Value implements the driver Valuer interface.
func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}

// GeoPoint - Spatial nullable point
type GeoPoint struct {
	Point Point
	Valid bool
}

// Scan implements the Scanner interface.
func (np *GeoPoint) Scan(val interface{}) error {
	if val == nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}

	point := &Point{}
	err := point.Scan(val)
	if err != nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}
	np.Point = Point{
		Lat: point.Lat,
		Lng: point.Lng,
	}
	np.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (np GeoPoint) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}
	return np.Point, nil
}
