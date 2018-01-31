package rupture

import (
	"strconv"
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/stretchr/testify/assert"
)

// vehicle an example struct for testing.
type vehicle struct {
	NumWheels int
	MaxSpeed  int
}

// add a collection of vehicles for testing
func addVehicles(t *testing.T, index func(string, interface{}) error) {
	for i := 0; i < 30; i++ {
		assert.NoError(t, index(strconv.Itoa(i), &vehicle{
			NumWheels: i % 3,
			MaxSpeed:  i,
		}))
	}
}

// search request for vehicles with one wheel, and a max speed of at least 10
func vehicleSearchRequest() *bleve.SearchRequest {
	one := float64(1)
	two := float64(2)
	numWheelsQuery := bleve.NewNumericRangeQuery(&one, &two)
	numWheelsQuery.SetField("NumWheels")

	ten := float64(10)
	maxSpeedQuery := bleve.NewNumericRangeQuery(&ten, nil)
	maxSpeedQuery.SetField("MaxSpeed")

	return bleve.NewSearchRequest(bleve.NewConjunctionQuery(
		numWheelsQuery,
		maxSpeedQuery,
	))
}
