package weather

import "strings"

var base32 = []byte("0123456789bcdefghjkmnpqrstuvwxyz")

// EncodeGeohash encodes lat/lon into a standard geohash with configurable precision.
func EncodeGeohash(lat, lon float64, precision int) string {
	if precision <= 0 {
		precision = 5
	}
	latMin, latMax := -90.0, 90.0
	lonMin, lonMax := -180.0, 180.0
	isEven := true
	bit := 0
	ch := 0
	var out strings.Builder
	out.Grow(precision)

	for out.Len() < precision {
		if isEven {
			mid := (lonMin + lonMax) / 2
			if lon >= mid {
				ch |= 1 << (4 - bit)
				lonMin = mid
			} else {
				lonMax = mid
			}
		} else {
			mid := (latMin + latMax) / 2
			if lat >= mid {
				ch |= 1 << (4 - bit)
				latMin = mid
			} else {
				latMax = mid
			}
		}
		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			out.WriteByte(base32[ch])
			bit = 0
			ch = 0
		}
	}
	return out.String()
}
