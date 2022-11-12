package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var cepsData string
var fileName string

func init() {
	flag.StringVar(&cepsData, "cep", "", "Cep das regiões, separados por vírgula.")
	flag.StringVar(&fileName, "file", "", "Arquivo contendo os ceps, 1 por linha.")
	flag.Parse()
}

type Endereco struct {
	Cep    string `json:"cep"`
	Rua    string `json:"logradouro"`
	Bairro string `json:"bairro"`
	Cidade string `json:"localidade"`
	Estado string `json:"uf"`
}

func (e Endereco) Stringfy() string {

	return fmt.Sprintf("Cep: %s\n%s, %s\n%s - %s\n", e.Cep, e.Rua, e.Bairro, e.Cidade, e.Estado)

}

func main() {
	ceps, err := processCeps()
	if err != nil {
		log.Fatal("Um erro aconteceu:", err)
		return
	}
	enderecos := getAddresses(ceps...)
	for _, endereco := range enderecos {
		fmt.Println(endereco.Stringfy())
	}
}

func processCeps() ([]string, error) {
	var ceps []string

	if fileName != "" {
		file, err := os.Open(fileName)
		if err != nil {
			return ceps, err
		}
		defer file.Close()
		fileScanner := bufio.NewScanner(file)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			ceps = append(ceps, fileScanner.Text())
		}
		return ceps, nil
	}

	ceps = append(ceps, strings.Split(cepsData, ",")...)
	return ceps, nil
}

func getAddresses(ceps ...string) []Endereco {
	var wg sync.WaitGroup
	end := make(chan Endereco)
	var enderecos []Endereco

	for _, cep := range ceps {
		wg.Add(1)
		go getAddress(cep, &wg, end)
	}

	go func() {
		for endereco := range end {
			enderecos = append(enderecos, endereco)
		}
	}()

	wg.Wait()
	close(end)

	return enderecos
}

func getAddress(cep string, wg *sync.WaitGroup, end chan<- Endereco) {
	defer wg.Done()
	request, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep))

	if err != nil {
		panic(err)
	}

	defer request.Body.Close()

	var endereco Endereco

	json.NewDecoder(request.Body).Decode(&endereco)

	end <- endereco
}
