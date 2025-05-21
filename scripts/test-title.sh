#!/bin/bash

# Configurar entorno
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/

# Configurar variables para la organización (Universidad)
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Función para emitir título
function emitTitle() {
  echo "Emitiendo título: $1 para estudiante: $3"
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C academictitles -n academic-titles --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c "{\"function\":\"IssueTitleToStudent\",\"Args\":[\"$1\", \"$2\", \"$3\", \"$4\", \"$5\"]}"
}

# Función para verificar título
function verifyTitle() {
  echo "Verificando título: $1"
  peer chaincode query -C academictitles -n academic-titles -c "{\"function\":\"VerifyTitle\",\"Args\":[\"$1\"]}"
}

# Parámetros por defecto
TITLE_ID="TITLE$(date +%s)"
STUDENT_ID="STUDENT001"
STUDENT_NAME="María García"
DEGREE="Ingeniería Informática"
DATE=$(date +%Y-%m-%d)

# Procesar argumentos
if [ "$1" == "emit" ]; then
  emitTitle "$2" "$3" "$4" "$5" "$6"
elif [ "$1" == "verify" ]; then
  verifyTitle "$2"
else
  echo "Uso: ./test-title.sh emit TITLE_ID STUDENT_ID STUDENT_NAME DEGREE DATE"
  echo "  o: ./test-title.sh verify TITLE_ID"
  echo "Ejemplo: ./test-title.sh emit TITLE001 STUDENT001 \"Ana López\" \"Ingeniería Informática\" 2025-03-06"
  echo "Ejemplo: ./test-title.sh verify TITLE001"
fi