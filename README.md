# Buscador de Endereços
Aplicação criada para buscar endereços a partir do cep

## Como usar
#### Para passar o cep, você pode passar uma flag cep com uma string contendo todos os ceps, separados por vírgula:
    -cep 11000000,22000000
#### É possível passar uma flag com um arquivo contendo um cep por linha:
    -fileName cep.txt

#### Para rodar a aplicação existem 2 opções:
    1º - Fazer o clone da aplicação e gerar o binário com o comando go build e executá-lo ./main -cep 11000000,22000000
    
    2º - Rodar via a imagem docker gerada: docker run hugovallada/get-address:latest -cep 11000000,22000000