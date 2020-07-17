package geo

import (
	"math"

	"github.com/mustafadubul/product/internal/domain"
)

const (
	earthRadius = 6371000
)

// https://www.movable-type.co.uk/scripts/latlong.html
// https://stackoverflow.com/questions/3695224/sqlite-getting-nearest-locations-with-latitude-and-longitude
func Destination(lat, lng, distance, bearing float64) domain.Point {

	dr := distance / earthRadius
	bearing2 := float64(bearing * (math.Pi / 180.0))

	lat1 := (lat * (math.Pi / 180.0))
	lng1 := (lng * (math.Pi / 180.0))

	sinφ1 := math.Sin(lat1) * math.Cos(dr)
	sinφ2 := math.Cos(lat1) * math.Sin(dr) * math.Cos(bearing)

	lat2 := math.Asin(sinφ1 + sinφ2)

	y := math.Sin(bearing2) * math.Sin(dr) * math.Cos(lat1)
	x := math.Cos(dr) - (math.Sin(lat1) * math.Sin(lat2))

	lng2 := lng1 + math.Atan2(y, x)
	lng2 = math.Mod((lng2+3*math.Pi), (2*math.Pi)) - math.Pi

	return domain.Point{
		X: lat2 * (180.0 / math.Pi),
		Y: lng2 * (180.0 / math.Pi),
	}
}

func BoundingBox(lat, lng, distance float64) []domain.Point {

	top := Destination(lat, lng, distance, 0)
	right := Destination(lat, lng, distance, 90)
	bottom := Destination(lat, lng, distance, 180)
	left := Destination(lat, lng, distance, 270)

	points := []domain.Point{top, right, bottom, left}
	return points
}
