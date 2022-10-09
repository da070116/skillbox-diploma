package collectors

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func validateCountry(rawData string) (country string, isOK bool) {
	availableCountries := []string{"RU", "US", "GB", "FR", "BL", "AT", "BG", "DK", "CA", "ES", "CH", "TR", "PE", "NZ", "MC"}
	if !Contains(availableCountries, rawData) {
		return
	}
	country = rawData
	isOK = true
	return
}

func validateBandwidth(rawData string) (bandwidth string, isOK bool) {
	bandwidthPercent, err := strconv.Atoi(rawData)
	if err != nil {
		return
	}
	if bandwidthPercent >= 0 && bandwidthPercent <= 100 {
		bandwidth = rawData
		isOK = true
	}
	return

}

func validateResponseTime(rawData string) (responseTime string, isOK bool) {
	_, err := strconv.Atoi(rawData)
	if err != nil {
		return
	}
	responseTime = rawData
	isOK = true
	return
}

func validateProvider(rawData string) (provider string, isOK bool) {
	allowedProviders := []string{"Topolo", "Rond", "Kildy"}
	if !Contains(allowedProviders, rawData) {
		return "", false
	}
	return rawData, true
}

func parseData(data [][]string) []SMSData {
	result := make([]SMSData, 0)
	i := 0
	for _, raw := range data {

		var country, bandwidth, responseTime, provider string
		var isOk bool

		parameters := strings.Split(raw[0], ";")
		if len(parameters) < 4 {
			continue
		}
		i++

		country, isOk = validateCountry(parameters[0])
		if !isOk {
			continue
		}

		bandwidth, isOk = validateBandwidth(parameters[1])
		if !isOk {
			continue
		}
		responseTime, isOk = validateResponseTime(parameters[2])
		if !isOk {
			continue
		}
		provider, isOk = validateProvider(parameters[3])
		if !isOk {
			continue
		}
		result = append(result, SMSData{
			Country:      country,
			Bandwidth:    bandwidth,
			ResponseTime: responseTime,
			Provider:     provider,
		})

	}
	return result
}

func SMSCollector() {
	csv_file, err := os.Open("sms.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csv_file.Close()

	csvReader := csv.NewReader(csv_file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	smsList := parseData(data)
	fmt.Println(smsList)
}
