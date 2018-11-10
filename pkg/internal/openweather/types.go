package openweather

// Result contains the response from the OpenWeather API
type Result struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}
