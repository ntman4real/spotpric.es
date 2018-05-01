package main

import (
	"encoding/xml"
	"net/http"
	"time"

	"gopkg.in/amz.v3/aws"
)

type priceHistoryResponse struct {
	Items []struct {
		InstanceType string  `xml:"instanceType"`
		Description  string  `xml:"productDescription"`
		Price        float32 `xml:"spotPrice"`
		Zone         string  `xml:"availabilityZone"`
	} `xml:"spotPriceHistorySet>item"`
}

func getPrices() (priceHistoryResponse, error) {
	req, err := http.NewRequest("GET", "https://ec2.amazonaws.com/", nil)
	if err != nil {
		return priceHistoryResponse{}, err
	}

	query := req.URL.Query()
	query.Add("Action", "DescribeSpotPriceHistory")
	query.Add("Version", "2016-11-15")

	var currentTime = time.Now().In(time.UTC).Format(time.RFC3339)
	query.Add("Timestamp", currentTime)
	query.Add("StartTime", currentTime)
	query.Add("Filter.1.Name", "product-description")
	query.Add("Filter.1.Value", "Linux/UNIX (Amazon VPC)")

	req.URL.RawQuery = query.Encode()

	auth, err := aws.EnvAuth()
	if err != nil {
		return priceHistoryResponse{}, err
	}

	aws.SignV2(req, auth)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return priceHistoryResponse{}, err
	}

	defer r.Body.Close()
	var resp priceHistoryResponse
	xml.NewDecoder(r.Body).Decode(&resp)
	return resp, nil
}
