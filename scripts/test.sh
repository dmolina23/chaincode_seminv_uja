#!/bin/bash

echo "   ________          _                      __   ______          __ "
echo "  / ____/ /_  ____ _(_)___  _________  ____/ /__/_  __/__  _____/ /_"
echo " / /   / __ \/ __ \`/ / __ \/ ___/ __ \/ __  / _ \/ / / _ \/ ___/ __/"
echo "/ /___/ / / / /_/ / / / / / /__/ /_/ / /_/ /  __/ / /  __(__  ) /_  "
echo "\____/_/ /_/\__,_/_/_/ /_/\___/\____/\__,_/\___/_/  \___/____/\__/  "
echo "                                                                    "
echo "                    (Developed by: Mz3r0on3)                         "
echo "---------------------------------------------------------------------"

# Función para mostrar el menú
show_menu() {
    echo
    echo "1. Instalar requisitos previos (Docker, Go, Hyperledger Fabric)"
    echo "2. Configurar proyecto"
    echo "3. Iniciar red de pruebas"
    echo "4. Desplegar ChainCode"
    echo "5. Realizar pruebas básicas"
    echo "6. Limpiar entorno y detener redes"
    echo "7. Salir"
    echo
    echo "Seleccione una opción: "
}

# Función para instalar requisitos previos
install_prerequisites() {
    echo "Instalando requisitos previos..."
    
    # Docker y Docker Compose
    sudo apt update
    sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    sudo usermod -aG docker $USER
    
    # Go
    wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
    echo 'export GOPATH=$HOME/go' >> ~/.profile
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile
    source ~/.profile
    
    # Hyperledger Fabric
    mkdir -p ~/hyperledger/fabric
    cd ~/hyperledger/fabric
    curl -sSL https://bit.ly/2ysbOFE -o bootstrap.sh
    chmod +x bootstrap.sh
    ./bootstrap.sh 2.2.5 1.5.2
    echo 'export PATH=$PATH:$HOME/hyperledger/fabric/fabric-samples/bin' >> ~/.profile
    source ~/.profile
    
    echo "Instalación completada. Por favor, reinicie su terminal."
}

# Función para configurar el proyecto
setup_project() {
    echo "Configurando proyecto..."
    mkdir -p ~/academic-titles-project/{chaincode,tests,scripts}
    cd ~/academic-titles-project
    mkdir -p chaincode/academic-titles
    cd chaincode/academic-titles
    go mod init github.com/yourusername/academic-titles
    echo "Configuración del proyecto completada."
}

# Función para iniciar la red de pruebas
start_network() {
    echo "Iniciando red de pruebas..."
    cd ~/hyperledger/fabric/fabric-samples/test-network
    ./network.sh down
    ./network.sh up createChannel -c academictitles
    
    # Configurar variables de entorno
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    
    echo "Red de pruebas iniciada y configurada."
}

# Función para desplegar el chaincode
deploy_chaincode() {
    echo "Desplegando ChainCode..."
    cd ~/hyperledger/fabric/fabric-samples/test-network
    ./network.sh deployCC -ccn academic-titles -ccp ~/academic-titles-project/chaincode/academic-titles -ccl go
    echo "ChainCode desplegado."
}

# Función para realizar pruebas básicas
run_basic_tests() {
    echo "Realizando pruebas básicas..."

    echo "Prueba de emisión de título..."

    # Emisión de título
    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C academictitles -n academic-titles --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"IssueTitleToStudent","Args":["{\"titleId\":\"TITLE001\",\"studentId\":\"STUDENT001\",\"studentName\":\"Ana López\",\"degree\":\"Ingeniería Informática\",\"emissionDate\":\"2025-03-06\", \"validationHash\":\"\"}"]}'

    # Esperar 3 segundos para que la transacción se procese
    sleep 3

    echo "Prueba de validación de título..."
    
    # Verificación de título
    peer chaincode query -C academictitles -n academic-titles -c '{"function":"VerifyTitle","Args":["TITLE001"]}'
    
    echo "Pruebas básicas completadas."
}

# Función para limpiar el entorno
clean_environment() {
    echo "Limpiando el entorno..."
    
    # Detener la red de Hyperledger si está corriendo
    cd ~/hyperledger/fabric/fabric-samples/test-network
    ./network.sh down

    # Eliminar todos los contenedores Docker relacionados con Hyperledger
    echo "Eliminando contenedores Docker..."
    docker ps -a | grep "dev-peer\|peer\|orderer\|ca\|couchdb" | awk '{print $1}' | xargs -r docker rm -f
    
    # Eliminar volúmenes Docker no utilizados
    echo "Eliminando volúmenes Docker..."
    docker volume prune -f
    
    # Eliminar redes Docker no utilizadas
    echo "Eliminando redes Docker..."
    docker network prune -f
    
    # Eliminar imágenes Docker relacionadas con Hyperledger
    echo "Eliminando imágenes Docker de Hyperledger..."
    docker images | grep "dev-peer\|hyperledger" | awk '{print $3}' | xargs -r docker rmi -f
    
    # Limpiar directorios de Hyperledger
    echo "Limpiando directorios de Hyperledger..."
    rm -rf ~/hyperledger/fabric/fabric-samples/test-network/organizations/ordererOrganizations
    rm -rf ~/hyperledger/fabric/fabric-samples/test-network/organizations/peerOrganizations
    rm -rf ~/hyperledger/fabric/fabric-samples/test-network/channel-artifacts/*
    rm -rf ~/hyperledger/fabric/fabric-samples/test-network/system-genesis-block/*
    
    echo "Limpieza completada. El entorno está listo para una nueva instalación."
}

# Bucle principal del menú
while true; do
    show_menu
    read -r opt
    case $opt in
        1) install_prerequisites ;;
        2) setup_project ;;
        3) start_network ;;
        4) deploy_chaincode ;;
        5) run_basic_tests ;;
        6) clean_environment ;;
        7) echo "Saliendo..."; exit 0 ;;
        *) echo "Opción inválida" ;;
    esac
    echo "Presione ENTER para continuar..."
    read -r
done