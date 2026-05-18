package weather

type Response struct {
	Temperature  string `json:"temperature"`
	Condition    string `json:"condition"`
	Humidity     string `json:"humidity"`
	WindSpeed    string `json:"wind_speed"`
	LocationName string `json:"location_name"`
}
