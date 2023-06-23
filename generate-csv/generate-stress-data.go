package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Driver struct {
	DriverID string `json:"driver_id"`
	ImageURL string  `json:"imageURL"`
	SelfieURL string  `json:"selfieURL"`
	ChallengeID string  `json:"challengeID"`
}

type Challenge struct {
	Name string `json:"name"`
}

type ChallengeRequest struct {
	FlowToken string `json:"flow_token"`
	Context   string `json:"context"`
}


type ChallengeResponseData  struct {
	Challenges []Challenge`json:"challenges"`
	ChallengeToken string `json:"challenge_token"`
}

type ChallengeResponse struct {
	Success bool  `json:"success"`
	Data ChallengeResponseData `json:"data"`
}

func main() {
	fileName := "input.csv"
	drivers :=readFromCsv(fileName)
	newDrivers := generateChallenges(drivers, 30)
	writetoCSV(newDrivers)
}

func writetoCSV(drivers []Driver) {

	// create a file
	file, err := os.Create("/Users/abhishek/Desktop/personal/machine-coding-practice/generate-csv/result.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	dataStr := [][]string{
		{"id" , "imageURL", "selfieURL", "challengeID"},
	}

	for _, driver := range drivers {
		 dataStr = append(dataStr, []string{
			 driver.DriverID , driver.ImageURL, driver.SelfieURL, driver.ChallengeID})
	 }


	// initialize csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// write all rows at once
	err = writer.WriteAll(dataStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	
}

func readFromCsv(fileName string ) []Driver {
	f, err := os.Open("/Users/abhishek/Desktop/personal/machine-coding-practice/generate-csv/input.csv")
	var drivers []Driver
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	firstRec := true
	csvReader := csv.NewReader(f)
	for {
		rec, err:= csvReader.Read()
		if err == io.EOF {
			break
		}
		if err!= nil {
			log.Fatal(err)
		}
		if !firstRec {
			drivers = append(drivers , Driver{
				DriverID:    rec[0],
				ImageURL:    rec[1],
				SelfieURL:   rec[2],
				ChallengeID: "",
			})
		}
		firstRec = false

	}
	return drivers

}


func generateChallenges(drivers []Driver, count int) []Driver{
	flowToken := "0d42a648-69fd-48a2-90a2-72c0462d1731"
	clientID := "driver-mfa-service"
	passKey := "06a61c88-baf9-479f-8323-1d5478ba6a3c"
	var newDrivers []Driver
	host := "driver-platform-face-recognition-service-internal.golabs.io"
	for _, driver := range drivers {
		for  i:=0; i< count; i++ {
			url := fmt.Sprintf("http://%s/internal/drivers/%s/challenge", host, driver.DriverID)
			chalengeID := createChallenge(driver.DriverID, flowToken ,  clientID, passKey, url)
			newDrivers = append(newDrivers, Driver{
				DriverID:    driver.DriverID,
				ImageURL:    driver.ImageURL,
				SelfieURL:   driver.SelfieURL,
				ChallengeID: chalengeID,
			})
		}
	}

	return newDrivers
}

func createChallenge(driverID, flowToken, clientID, passKey, url  string) string{
	client := http.Client{}

	challengeReq := ChallengeRequest{
		FlowToken: flowToken,
		Context:   "{\"profile\":{\"phone_number\":\"940193897\",\"country_code\":\"+91\",\"id\":\"940412600\",\"payload\":{\"all_session_count\":\"2\",\"category\":\"4W\", \"countryCode\": \"ID\", \"email\":\"Vikashcardriver@google.com\",\"name\":\"Vikash\",\"same_device_session_count\":\"0\"},\"scopes\":[\"d:c:*\"]},\"account_id\":\"d74870c3-0f36-4250-80f7-c1950254225a\"}",
	}
	jsonStr, err := json.Marshal(challengeReq)
	responseBody := bytes.NewBuffer(jsonStr)

	req, err := http.NewRequest("POST", url,  responseBody )
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Pass-Key", passKey)
	req.Header.Set("Accept-Language", "en")

	if err != nil{
		fmt.Errorf(err.Error())
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	var challengeRes ChallengeResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &challengeRes)
	if err!= nil {
		fmt.Errorf(err.Error())
	}
	return challengeRes.Data.ChallengeToken
}
