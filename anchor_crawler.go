package main

import (
	"io"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
	"encoding/json"
	"os"
	"log"
)

const (
	deposit_api = `https://api.anchorprotocol.com/api/v1/deposit`
	beth_api = `https://api.anchorprotocol.com/api/v1/bassets/beth`
	borrow_api = `https://api.anchorprotocol.com/api/v1/borrow`
	bluna_api = `https://api.anchorprotocol.com/api/v1/bassets/bluna`
	ust_api = `https://api.anchorprotocol.com/api/v1/market/ust`
	// ancPrice_api = "https://api.anchorprotocol.com/api/v1/anc"

)

func main(){
	yield_reserve_res := getApiResult(ust_api)
	deposit_ust_api_res := getApiResult(deposit_api)
	beth_api_res := getApiResult(beth_api)
	bluna_api_res := getApiResult(bluna_api)
	borrow_api_res := getApiResult(borrow_api)

	yield_reserve := parseResult(yield_reserve_res,"overseer_ust_balance")/1000000
	deposit_ust := parseResult(deposit_ust_api_res,"total_ust_deposits")/1000000
	borrowed_ust := parseResult(borrow_api_res,"total_borrowed")/1000000
	beth_Collateral := parseResult(beth_api_res,"total_collateral")/1000000
	beth_price := parseResult(beth_api_res,"beth_price")
	bluna_Collateral := parseResult(bluna_api_res,"total_collateral")/1000000
	bluna_price := parseResult(bluna_api_res,"bLuna_price")

	bEth_rate := 4.60/100
	bLuna_rate := 7.9/100
	borrow_rate := 14.58/100
	deposit_apy := parseResult(yield_reserve_res,"deposit_apy")
	//To do 目前把浮點數加工x10000在除以10000來保留小數點後兩位 肯定有更好的方式～
	deposit_apy_float := float64(deposit_apy)/10000

	beth_profit := float64(beth_Collateral*beth_price)*bEth_rate
	bluna_profit := float64(bluna_Collateral*bluna_price)*bLuna_rate
	loan_profit := float64(borrowed_ust)*borrow_rate
	total_profit := beth_profit + bluna_profit + loan_profit
	platform_cost := float64(deposit_ust) * deposit_apy_float
	revenue := beth_profit + bluna_profit + loan_profit - platform_cost

	logFile, err := os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err != nil{
		log.Fatalf("file open error: %v", err)
	}
	defer logFile.Close()
	console := io.MultiWriter(os.Stdout,logFile)
	log.SetOutput(console)
	

	log.Printf("Yield Reserve : %d\nDeposited ust : %d\nbeth Value    : %.02f\nbluna Value   : %.02f\nborrowed ust  : %d\n",yield_reserve, deposit_ust,float64(beth_Collateral*beth_price),float64(bluna_Collateral*bluna_price),borrowed_ust)
	log.Println("-------------------------")
	log.Println("total profit  :",int(total_profit))
	log.Println("platform_cost :",int(platform_cost))
	log.Println("Total Revenue :",int(revenue))
	log.Printf("平台的deposit apy收益應該要是 %.02f ％而不是 %.02f ％才能打平\n",(total_profit/float64(deposit_ust)*100),deposit_apy_float*100)
	log.Printf("照這樣下去的話，在%.02f天Terra就要開始負債囉~", (float64(yield_reserve)/(-revenue/365)))
}

func getApiResult(url string) string{
	req, _ := http.NewRequest("GET",url,nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("fetch %T fail",res)
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	return string(result)
}

func parseResult(result string, key string) int{
	
	var raw map[string]json.RawMessage
	json.Unmarshal([]byte(result),&raw)

	parsed := make(map[string]interface{}, len(raw))
	for key, val := range raw {
		var v interface{}
		err := json.Unmarshal(val, &v)
		if err == nil {
			parsed[key] = v
			continue
		}
		parsed[key] = val
	}
	remaining_ust  := parsed[key].(string)
	answer, err := strconv.Atoi(remaining_ust)

	if err != nil {
		float_answer, _ := strconv.ParseFloat(remaining_ust,64)
		if int(float_answer) == 0 {
			return int(float_answer*10000)
		}
		return int(float_answer)
	}
	return answer
}