package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Address structure represents the address format used in the application.
type Address struct {
	Cep        string `json:"cep"`        // Postal code
	Logradouro string `json:"logradouro"` // Street name
	Bairro     string `json:"bairro"`     // Neighborhood
	Localidade string `json:"localidade"` // City
	Uf         string `json:"uf"`         // State
}

// Structure for the response from BrasilAPI
type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

// Structure for the response from ViaCEP
type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// fetchAddress fetches the address information from the given API URL.
// It sends the result to the result channel if successful or an error to the errChan channel if it fails.
// It uses the isBrasilAPI flag to determine which response structure to decode into.
func fetchAddress(ctx context.Context, url string, result chan<- Address, errChan chan<- error, isBrasilAPI bool) {
	// Create a new HTTP request with the given context
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req) // Send the request
	if err != nil {
		errChan <- err // Send error to error channel if the request fails
		return
	}
	defer resp.Body.Close()

	var address Address

	if isBrasilAPI {
		// Decode response into BrasilAPIResponse struct if it's from BrasilAPI
		var brasilAPIResponse BrasilAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&brasilAPIResponse); err != nil {
			errChan <- err // Send error to error channel if decoding fails
			return
		}
		// Map BrasilAPIResponse to the Address struct
		address = Address{
			Cep:        brasilAPIResponse.Cep,
			Logradouro: brasilAPIResponse.Street,
			Bairro:     brasilAPIResponse.Neighborhood,
			Localidade: brasilAPIResponse.City,
			Uf:         brasilAPIResponse.State,
		}
	} else {
		// Decode response into ViaCEPResponse struct if it's from ViaCEP
		var viaCEPResponse ViaCEPResponse
		if err := json.NewDecoder(resp.Body).Decode(&viaCEPResponse); err != nil {
			errChan <- err // Send error to error channel if decoding fails
			return
		}
		// Map ViaCEPResponse to the Address struct
		address = Address{
			Cep:        viaCEPResponse.Cep,
			Logradouro: viaCEPResponse.Logradouro,
			Bairro:     viaCEPResponse.Bairro,
			Localidade: viaCEPResponse.Localidade,
			Uf:         viaCEPResponse.Uf,
		}
	}

	// Send the result to the result channel
	result <- address
}

func main() {
	// Set a default postal code (CEP)
	cep := "01153000"

	// If a postal code is provided as a command-line argument, use that instead
	if len(os.Args) > 1 {
		cep = os.Args[1]
	}

	// URLs for BrasilAPI and ViaCEP with the provided postal code
	brasilAPI := "https://brasilapi.com.br/api/cep/v1/" + cep
	viaCEP := "http://viacep.com.br/ws/" + cep + "/json/"

	// Create channels to receive the address result or an error
	brasilAPIResult := make(chan Address)
	brasilAPIErr := make(chan error)
	viaCEPResult := make(chan Address)
	viaCEPErr := make(chan error)

	// Create a context with a 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is canceled when we're done

	// Start fetching address from BrasilAPI in a separate goroutine
	go fetchAddress(ctx, brasilAPI, brasilAPIResult, brasilAPIErr, true)

	// Start fetching address from ViaCEP in a separate goroutine
	go fetchAddress(ctx, viaCEP, viaCEPResult, viaCEPErr, false)

	// Wait for the first successful response or a timeout
	select {
	case addr := <-brasilAPIResult:
		// Print the address from BrasilAPI if it's successful
		fmt.Printf("Endereço (Brasil API): %s, %s, %s - %s\n", addr.Logradouro, addr.Bairro, addr.Localidade, addr.Uf)
	case err := <-brasilAPIErr:
		// Print the error from BrasilAPI if it occurs
		fmt.Printf("Erro (Brasil API): %v\n", err)

	case addr := <-viaCEPResult:
		// Print the address from ViaCEP if it's successful
		fmt.Printf("Endereço (ViaCEP): %s, %s, %s - %s\n", addr.Logradouro, addr.Bairro, addr.Localidade, addr.Uf)
	case err := <-viaCEPErr:
		// Print the error from ViaCEP if it occurs
		fmt.Printf("Erro (ViaCEP): %v\n", err)

	case <-ctx.Done():
		// Print a timeout message if neither API responds within the timeout period
		fmt.Println("Timeout: Nenhuma das APIs respondeu em tempo hábil.")
	}
}
