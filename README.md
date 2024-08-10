# Address Fetcher Service

This project is a Go application that retrieves address information based on a provided Brazilian postal code (CEP) by querying two different APIs concurrently. The program returns the address data from the API that responds the fastest, ensuring a quick and efficient result.

## How it Works

The application sends concurrent requests to the following APIs:

- [BrasilAPI](https://brasilapi.com.br/api/cep/v1/01153000)
- [ViaCEP](https://viacep.com.br/ws/01153000/json/)

It uses the API that provides the first response, ignoring the slower one. The program also includes a timeout mechanism to ensure the entire process completes within 1 second. If neither API responds within the time limit, a timeout message is displayed.

## Requirements

- Go 1.16 or later.

## How to Run the Application

1. Clone the repository or save the code to a file named main.go
2. Ensure you have Go installed on your system.
3. Run the application:

```sh
go run main.go [CEP]
```

Replace [CEP] with the desired Brazilian postal code. If no CEP is provided, the application defaults to 01153000.

## Output

The application prints the fetched address to the command line with the following format:

>Endereço (API): Logradouro, Bairro, Localidade - UF

If the request exceeds the 1-second timeout or an error occurs, the program will print a corresponding error message.

## Example

```sh
go run main.go 01001000
```

Expected output:

```less
// BrasilAPI
Endereço (Brasil API): Praça da Sé, Sé, São Paulo - SP
```

or

```less
// ViaCEP
Endereço (ViaCEP): Praça da Sé, Sé, São Paulo - SP
```

If a timeout occurs:

```less
// ViaCEP
Timeout: Nenhuma das APIs respondeu em tempo hábil.
```

## Aditional Information

- **Timeout:** The application enforces a 1-second timeout to ensure responsiveness.
- **APIs:** The application chooses the first responding API to minimize the wait time for the user.
