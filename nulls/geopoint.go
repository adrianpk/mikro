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
	"errors"
	"fmt"
)

// Point - Spatial point
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// MakeGeoPoint - Returns a valid
func MakeGeoPoint(lat, long float64, valid bool) *GeoPoint {
	return &GeoPoint{Point{lat, long}, false}
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

// Match - Custom model comparator.
func (p Point) Match(tc Point) bool {
	r := p.Lat == tc.Lat &&
		p.Lng == tc.Lng
	return r
}

// GeoPoint - Spatial nullable point
type GeoPoint struct {
	Point Point
	Valid bool
}

// Latitude - Returns latitude.
func (ngp *GeoPoint) Latitude() float64 {
	if ngp.Valid {
		return ngp.Point.Lat
	}
	return 0.0
}

// Longitude - Returns latitude.
func (ngp *GeoPoint) Longitude() float64 {
	if ngp.Valid {
		return ngp.Point.Lng
	}
	return 0.0
}

// Scan implements the Scanner interface.
func (ngp *GeoPoint) Scan(val interface{}) error {
	if val == nil {
		ngp.Point, ngp.Valid = Point{}, false
		return nil
	}

	point := &Point{}
	err := point.Scan(val)
	if err != nil {
		ngp.Point, ngp.Valid = Point{}, false
		return nil
	}
	ngp.Point = Point{
		Lat: point.Lat,
		Lng: point.Lng,
	}
	ngp.Valid = true

	return nil
}

// Value implements the driver Valuer interface.
func (ngp GeoPoint) Value() (driver.Value, error) {
	if !ngp.Valid {
		return nil, errors.New("Invalid GeoPoint data")
	}
	return ngp.Point, nil
}

// NullZeroPoint - GeoPoint zero value.
func NullZeroPoint() GeoPoint {
	return GeoPoint{
		Point: Point{0, 0},
		Valid: false,
	}
}

// Match - Custom model comparator.
func (ngp GeoPoint) Match(tc GeoPoint) bool {
	r := ngp.Point.Match(tc.Point) &&
		ngp.Valid == tc.Valid
	return r
}
