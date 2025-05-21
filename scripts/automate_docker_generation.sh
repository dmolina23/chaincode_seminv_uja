#!/bin/bash

# Variables
test_network_dir=~/hyperledger/fabric/fabric-samples/test-network
chaincode_dir=~/academic-titles-project/chaincode/academic-titles
channel_name=academictitles
cc_name=academic-titles

# Función para verificar la instalación de Docker
docker_check() {
    if ! command -v docker &> /dev/null; then
        echo "Docker no está instalado. Instalándolo ahora..."
        sudo apt update
        sudo apt install -y docker.io
    fi
    echo "Docker está instalado."
}

docker_compose_check() {
    if ! command -v docker-compose &> /dev/null; then
        echo "Docker Compose no está instalado. Instalándolo ahora..."
        sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
    fi
    echo "Docker Compose está instalado."
}

# Función para iniciar la red Hyperledger Fabric
start_network() {
    cd $test_network_dir
    echo "Deteniendo cualquier red anterior..."
    ./network.sh down
    echo "Iniciando la red de pruebas con el canal $channel_name..."
    ./network.sh up createChannel -c $channel_name
}

# Función para desplegar el chaincode
deploy_chaincode() {
    cd $test_network_dir
    echo "Desplegando el chaincode $cc_name..."
    ./network.sh deployCC -ccn $cc_name -ccp $chaincode_dir -ccl go
}

# Función para ejecutar pruebas básicas
test_chaincode() {
    echo "Ejecutando pruebas básicas..."
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $channel_name -n $cc_name --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt -c '{"function":"IssueTitleToStudent","Args":["TITLE001", "STUDENT001", "Ana López", "Ingeniería Informática", "2025-03-06"]}'
    
    peer chaincode query -C $channel_name -n $cc_name -c '{"function":"VerifyTitle","Args":["TITLE001"]}'
}

# Ejecutar funciones
docker_check
docker_compose_check
start_network
deploy_chaincode
test_chaincode

echo "Pruebas completadas."
