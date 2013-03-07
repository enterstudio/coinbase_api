package coinbase_api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// An AuthenticatedRequest needs to have a way to load the API key.
type AuthenticatedRequest interface {
	SetApiKey()
}

// PostAuthenticatedRequest makes an authenticated JSON POST request to the API.
// The structure containing the result should be passed into res.
func PostAuthenticatedRequest(data AuthenticatedRequest, endpoint string, res interface{}) (err error) {
	if ApiKey == "" {
		return ErrNotAuthenticated
	}
	data.SetApiKey()

	request_json_body, err := json.Marshal(data)
	if err != nil {
		return
	}
	request_body := bytes.NewBuffer(request_json_body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", api_base+endpoint, request_body)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.Unmarshal(response_body, &res)
	return
}

// Make an authenticated GET request to the API.
func GetAuthenticatedRequest(data AuthenticatedRequest, endpoint string, res interface{}) (err error) {
	if ApiKey == "" {
		return ErrNotAuthenticated
	}
	data.SetApiKey()
	request_json_body, err := json.Marshal(data)
	if err != nil {
		return
	}
	request_body := bytes.NewBuffer(request_json_body)

	client := &http.Client{}
	req, err := http.NewRequest("GET", api_base+endpoint, request_body)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.Unmarshal(response_body, &res)
	return
}

// Make an unauthenticated GET request.
func GetUnauthenticatedRequest(data interface{}, endpoint string, res interface{}) (err error) {
	var req *http.Request
	if data != nil {
		var request_json_body []byte
		request_json_body, err = json.Marshal(data)
		if err != nil {
			return
		}
		request_body := bytes.NewBuffer(request_json_body)
		req, err = http.NewRequest("GET", api_base+endpoint, request_body)
		req.Header.Add("content-type", "application/json")
	} else {
		req, err = http.NewRequest("GET", api_base+endpoint, nil)
	}

	client := &http.Client{}
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	response_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.Unmarshal(response_body, &res)
	return
}

// Retrieve exchange rates for a list of currencies.
func GetExchangeRates(currencies []string) (ExchangeRate, error) {
	exch := make(ExchangeRate, 0)
	endpoint := "currencies/exchange_rates"
	err := GetUnauthenticatedRequest(nil, endpoint, &exch)
	if err == nil {
		for k, _ := range exch {
			var keep bool
			for _, currency := range currencies {
				if k == currency {
					keep = true
					break
				}
			}
			if !keep {
				delete(exch, k)
			}
		}

	}
	return exch, err
}

// PurchaseBTC attempts to purchase the specified quantity of bitcoins.
func PurchaseBTC(qty float64) (p *Transaction, err error) {
	if qty < MinimumPurchase {
		err = ErrMinimumSubtotal
		return
	}
	pr := new(TransactionRequest)
	pr.Qty = qty

	p = new(Transaction)
	endpoint := "buys"

	err = PostAuthenticatedRequest(pr, endpoint, &p)
	return
}

// SellBTC attempts to sell the specified quantity of bitcoins.
func SellBTC(qty float64) (p *Transaction, err error) {
	if qty < MinimumPurchase {
		err = ErrMinimumSubtotal
		return
	}
	pr := new(TransactionRequest)
	pr.Qty = qty

	p = new(Transaction)
	endpoint := "buys"

	err = PostAuthenticatedRequest(pr, endpoint, &p)
	return
}

// GetAccountBalance retrieves the number of bitcoins in the user's account.
func GetAccountBalance() (b *Balance, err error) {
	get := new(GetAuthenticated)
	endpoint := "account/balance"

	b = new(Balance)
	err = GetAuthenticatedRequest(get, endpoint, &b)
	return
}

// GetBitcoinAddress retrieves the number of bitcoins in the user's account.
func GetReceiveAddress() (a *ReceiveAddress, err error) {
	get := new(GetAuthenticated)
	endpoint := "account/receive_address"

	a = new(ReceiveAddress)
	err = GetAuthenticatedRequest(get, endpoint, &a)
	return
}

// GetSellPrice returns the total returns from selling a certain number
// of bitcoins, accounting for transaction fees and market depth.
func GetSellPrice(qty float64) (b *Balance, err error) {
	var quantity struct {
		Qty float64 `json:"qty"`
	}
	quantity.Qty = qty
	endpoint := "prices/sell"
	b = new(Balance)
	err = GetUnauthenticatedRequest(quantity, endpoint, &b)
	return
}

// GetBuyPrice returns the total cost for purchasing a certain number of
// bitcoins, accounting for transaction fees and market depth.
func GetBuyPrice(qty float64) (b *Balance, err error) {
	var quantity struct {
		Qty float64 `json:"qty"`
	}
	quantity.Qty = qty
	endpoint := "prices/buy"
	b = new(Balance)
	err = GetUnauthenticatedRequest(quantity, endpoint, &b)
	return
}