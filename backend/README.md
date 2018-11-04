

# Configurar
Exportar uma variável de ambiente com o endereço do proxy do Wavefront

e.g.: 
export WF_PROXY="192.168.1.100:2878"

docker run -d -p 9090:9090 -e WF_PROXY=$WF_PROXY duboc/cdbackend:1.1

# Executar
go run main.go or build docker image and docker run -d -p 9090:9090 -e WF_PROXY=$WF_PROX
