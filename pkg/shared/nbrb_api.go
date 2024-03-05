package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Nbrb struct {
	ApiUrl string `mapstructure:"api_url"`
	RubId  string `mapstructure:"rub_id"`
	UsdId  string `mapstructure:"usd_id"`
}
type Currency struct {
	Rub      float64 `json:"rateRub"`
	RubScale float64 `json:"rateRubScale"`
	Usd      float64 `json:"rateUsd"`
}

func (nbrb *Nbrb) GetCurrencyRate(date string) (*Currency, error) {

	currency := Currency{}

	// IDs to retrieve data for
	ids := []string{nbrb.RubId, nbrb.UsdId}

	for i, id := range ids {
		// URL of the API with the current ID
		url := fmt.Sprintf(nbrb.ApiUrl + id + "?ondate=" + date + "T00:00:00")

		data, err := nbrb.New(url)
		if err != nil {
			return &Currency{}, err
		}

		// Access the data
		rate := data["Cur_OfficialRate"].(float64)
		scale := data["Cur_Scale"].(float64)
		if i == 0 {
			currency.Rub = rate
			currency.RubScale = scale
		} else if i == 1 {
			currency.Usd = rate
		}
	}

	// Return the currency rates
	return &currency, nil
}

func (nbrb *Nbrb) New(url string) (map[string]interface{}, error) {
	// Send GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	// Parse the JSON response
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return data, nil
}
