/**
 * Copyright (c) 2019 Adrian P.K. <apk@kuguar.io>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package nulls

import (
	"time"

	"github.com/gobuffalo/pop/nulls"
	uuid "github.com/gofrs/uuid"
)

var (
	// Double colons in order to avoid error message
	// "unexpected `:` while reading named param at"
	// in sqlx NamedExec.
	dateFormat = "2006-01-02T15::04::05.999Z"
)

// ToZeroUUID - Return a nullable version of zero value UUID.
func ToZeroUUID() nulls.UUID {
	return ToNullsUUID(uuid.Nil)
}

// ToNullsUUID - Return a nullable version of the passed argument.
func ToNullsUUID(u uuid.UUID) nulls.UUID {
	return nulls.NewUUID(u)
}

// ToNullsString - Return a nullable version of the passed argument.
func ToNullsString(str string) nulls.String {
	return nulls.NewString(str)
}

// ToZeroInt64 - Return a nullable version of zero Int64.
func ToZeroInt64() nulls.Int64 {
	return nulls.NewInt64(0)
}

// ToNullsInt64 - Return a nullable version of the passed argument.
func ToNullsInt64(i int64) nulls.Int64 {
	return nulls.NewInt64(i)
}

// ToZeroFloat64 - Return a nullable version of zero Float64.
func ToZeroFloat64() nulls.Float64 {
	return nulls.NewFloat64(0)
}

// ToFoat64 - Return a nullable version of the passed argument.
func ToFoat64(f float64) nulls.Float64 {
	return nulls.NewFloat64(f)
}

// ToNullsBool - Return a nullable version of the passed argument.
func ToNullsBool(bln bool) nulls.Bool {
	return nulls.NewBool(bln)
}

// ToTime - Return a nullable version of the passed argument.
func ToTime(t time.Time) nulls.Time {
	return nulls.NewTime(t)
}

// EmptyString - Return a nullable version of "" string.
func EmptyString() nulls.String {
	return ToNullsString("")
}

// TrueBool - Return a nullable version of bool value 'true'.
func TrueBool() nulls.Bool {
	return ToNullsBool(true)
}

// FalseBool - Return a nullable version of bool value 'false'.
func FalseBool() nulls.Bool {
	return ToNullsBool(false)
}

// ToZeroTime - Return a nullable version of zero time.
func ToZeroTime() nulls.Time {
	return ToTime(time.Time{})
}

// NowTime - Return a nullable version of current time.
func NowTime() nulls.Time {
	return ToTime(time.Now())
}

// ToZeroPoint - Return a nullable version of geolocation zero point.
func ToZeroPoint() nulls.Time {
	return ToTime(time.Now())
}

// FormatDate - Format date to use in SQL queries and statements.
func FormatDate(time nulls.Time) string {
	return time.Time.Format(dateFormat)
}
