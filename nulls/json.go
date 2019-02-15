/**
 *Copyright (c) 2018 Adrian P.K. <apk@kuguar.io>
 *
 *This software is released under the MIT License.
 *https://opensource.org/licenses/MIT
 */

package nulls

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSON is a property map representation of JSON
// type JSON map[string]interface{}

// JSONB is a nullable JSON
type JSONB struct {
	ByteArray []byte
	Valid     bool
}

// Scan is an sql.Scanner interface implementation for JSONB type
func (nj *JSONB) Scan(val interface{}) error {
	js, ok := val.([]byte)
	if !ok || js == nil {
		nj.ByteArray, nj.Valid = []byte{}, false
		return nil
	}
	var i interface{}
	err := json.Unmarshal(js, &i)
	if err != nil {
		nj.ByteArray, nj.Valid = []byte{}, false
		return nil
	}
	nj.ByteArray = js
	nj.Valid = true
	return nil
}

// Value is a driver.Valuer interface implementation for JSONB type.
func (nj JSONB) Value() (driver.Value, error) {
	// byteArray, err := json.Marshal(nj.ByteArray)
	if !nj.Valid {
		return nil, errors.New("Invalid JSONB data")
	}
	return nj.ByteArray, nil
}

func (nj JSONB) String() string {
	return string(nj.ByteArray)
}

// Match - Custom model comparator.
func (nj JSONB) Match(tc JSONB) bool {
	r := string(nj.ByteArray) == string(tc.ByteArray)
	return r
}
