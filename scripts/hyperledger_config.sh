#!/bin/bash

# Script para configurar red de prueba para chaincode de títulos académicos

# Requisitos:
# - Hyperledger Fabric instalado (v2.2+)
# - Docker y Docker Compose

# Directorio para test-network
cd $HOME/fabric-samples/test-network

# Limpiar red anterior si existe
./network.sh down

# Iniciar red con un canal
./network.sh up createChannel -c testChannel

# Configurar variables para la orgnización (Universidad)
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="UniversityMSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/university.example.com/peers/peer0.university.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/perrOrganizations/university.example.com/users/Admin@university.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Empaquetar chaincode
cd ../chaincode/academic-titles/
GO111MODULE=on go mod vendor
cd ../../test-network
peer lifecycle chaincode package academic-titles.tar.gz --path ../chaincode/academic-titles/ --lang golang --label academic-titles_1.0

# Instalar chaincode en la red
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID testChannel --name academic-titles --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Preparar variables para segunda organización (Verificador)
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# Aprobar chaincode para organización Verificador
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID titleschannel --name academic-titles --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Solicitar commit del chaincode
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID titleschannel --name academic-titles --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

echo "Chaincode instalado y listo para pruebas"