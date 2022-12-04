package services

import (
	"gl-farming/app/helper"
	"gl-farming/app/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s CurrencyServiceImpl) GetCurrencyRates() (map[string]float64, error) {

	var currencyRate models.ValuteRate

	requestDate := time.Now().Format("02-01-2006")

	if err := s.RequestCurrencyRate(&currencyRate, requestDate); err != nil {
		return nil, err
	}

	currencyRates := make(map[string]float64)

	for _, currency := range currencyRate.ValuteList {

		strVal := strings.Replace(currency.Value, ",", ".", -1)
		value, err := strconv.ParseFloat(strVal, 64)

		if err != nil {
			continue
		}
		value = helper.RoundFloat(value, 2)

		currencyRates[currency.CharCode] = value

	}

	return currencyRates, nil
}

func (srvc CurrencyServiceImpl) RequestCurrencyRate(currencyRates *models.ValuteRate, requestDate string) error {

	baseUrl := "http://www.cbr.ru/scripts/XML_daily.asp?date_req="

	ulrPath := baseUrl + requestDate

	request, err := http.NewRequest(http.MethodGet, ulrPath, nil)

	if err != nil {
		return err
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	defer response.Body.Close()

	decoder := helper.NewDecoderXML(response.Body)

	if err := decoder.Decode(&currencyRates); err != nil {
		return err
	}

	return nil
}
