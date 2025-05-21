# Documentación para la instalación y configuración del entorno de pruebas para la implementación de un NFT en Hyperledger Fabric
Este documento es una guía paso a paso de instalación y configuración de entornos de prueba para el testeo de la implementación de un NFT en una red blockchain permisionada basada en Hyperledger Fabric.

La finalidad de este documento es obtener un entorno de pruebas funcional para testear nuestra chaincode.

Este documento cubre desde la instalación de requisitos previos hasta la ejecución de pruebas básicas con el chaincode desplegado en una red simulada en Hyperledger Fabric.

## Consideraciones de la versión actual
> 1. El script automático y la automatización de la generación de nodos ficticios usando Docker no están testeados y pueden causar errores
> 2. La configuración actual no contempla el uso de las políticas de transacción definidas en los últimos cambios: para aplicar el archivo ACL (Access Control List) debemos realizar algunos cambios en el paso de la generación del canal
> 3. Todavía no se han realizado pruebas de integración sobre la limitación de transferencias ni sobre la obtención de títulos asignados a un alumno. Esto quiere decir que puede la implementación puede dar fallos al intentar realizar pruebas de integración (posibles errores en la interacción con el chaincode para realizar pruebas básicas)
> 4. La parte de pruebas con NodeJS no se ha llevado a cabo, en un futuro se implementará una API wallet (utilizando Express y Node) que nos permitirá obtener un MVP de este proyecto 

## 1. Instalación de requisitos previos
### 1.1 Instalación de Docker y Docker Compose
```bash
# Actualizar repositorios
sudo apt update

# Instalar paquetes necesarios para permitir a apt usar repositorios HTTPs
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# Añadir la clave gpg oficial de Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add
-
# Añadir el repositorio Docker a las fuentes de APT
sudo add-apt-repository "deb [arch=amd64]
https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# Actualizar la base de datos de paquetes
sudo apt update

# Instalar Docker CE
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Instalar Docker Compose
sudo curl -L
"https://github.com/docker/compose/releases/download/1.29.2/docker-
compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Añadir tu usuario al grupo docker para evitar usar sudo
sudo usermod -aG docker $USER

# Verificar instalaciones
docker --version
docker-compose --version
```

Después de ejecutar estos comandos, es recomendable cerrar la terminal y volver a iniciarla de modo que se apliquen todos los cambios.

### 1.2 Instalación de Go (Golang)
```bash
# Descargar Go
wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz

# Extraer en /usr/local
sudo tar -C /usr/local -xzf go1.19.5.linux-amd64.tar.gz

# Configurar variables de entorno
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
echo 'export GOPATH=$HOME/go' >> ~/.profile
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile

# Aplicar cambios
source ~/.profile

# Verificar instalación
go version
```

### 1.3 Instalación de Hyperledger Fabric
```bash
# Crear directorio
mkdir -p ~/hyperledger/fabric
cd ~/hyperledger/fabric

# Descargar el script de bootstrap
curl -sSL https://bit.ly/2ysbOFE -o bootstrap.sh

# Dar permisos de ejecución
chmod +x bootstrap.sh

# Ejecutar script para descargar binarios, imágenes y muestras
./bootstrap.sh 2.2.5 1.5.2

# Añadir binarios de Fabric al PATH
echo 'export PATH=$PATH:$HOME/hyperledger/fabric/fabric-samples/bin' >>
~/.profile

source ~/.profile
```

## 2. Configuración del proyecto
```bash
# Crear estructura de directorios para el proyecto
mkdir -p ~/academic-titles-project/{chaincode,tests,scripts}
cd ~/academic-titles-project

# Preparar directorio para el chaincode
mkdir -p chaincode/academic-titles
cd chaincode/academic-titles

# Inicializar módulo Go
go mod init github.com/yourusername/academic-titles
```

> NOTA: La inicialización de go mod presupone que el código está subido a un repositorio en GitHub

## 3. Testeo de red y despliegue del ChainCode
```bash
# Ir al directorio de test-network
cd ~/hyperledger/fabric/fabric-samples/test-network

# Detener cualquier red anterior
./network.sh down

# Iniciar nueva red con canal
./network.sh up createChannel -c academictitles

# Configurar variables para la organización (Universidad)
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export
CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.exa
mple.com/peers/peer0.org1.example.com/tls/ca.crt
export
CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example
.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Desplegar el chaincode
./network.sh deployCC -ccn academic-titles -ccp ~/academic-titles-
project/chaincode/academic-titles -ccl go
```

## 4. Interacción con el ChainCode para realizar pruebas básicas
```bash
# Invocar la función de emisión de título
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride
orderer.example.com --tls --cafile
${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.exam
ple.com/msp/tlscacerts/tlsca.example.com-cert.pem -C academictitles -n
academic-titles --peerAddresses localhost:7051 --tlsRootCertFiles
${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.ex
ample.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles
${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.ex
ample.com/tls/ca.crt -c '{"function":"IssueTitleToStudent","Args":
["TITLE001", "STUDENT001", "Ana López", "Ingeniería Informática", "2025-03-
06"]}'

# Consultar un título
peer chaincode query -C academictitles -n academic-titles -c
'{"function":"VerifyTitle","Args":["TITLE001"]}'
```

## 5. Descripción del código
El archivo `chaincode.go` representa nuestro SmartContract. Ha sido desarrollado en Go utilizando la API que proporciona Hyperledger Fabric.

También nos encontramos con el archivo `chaincode_test.go`, que simplemente implementa tests unitarios para comprobar que no existen fallos en la implementación.

El archivo `app.test.js` es un test de integración implementado en NodeJS. 

## 6. Paso a paso de realización de pruebas con NodeJS (NO IMPLEMENTADO)
Para poder comprobar su funcionamiento necesitaremos realizar los siguientes pasos:

### 6.1 Verificar e instalar dependencias
El script usa `fabric-network` y `fabric-ca-client`. Asegúrate de tenerlas instaladas:

```
npm install fabric-network fabric-ca-client
```

### 6.2 Asegurar que la red Hyperledger Fabric está corriendo
Antes de ejecutar la prueba, necesitamos verificar que la red Fabric esté levantada. Si usamos `fabric-samples` para configurar la red, podemos asegurarnos de que esté corriendo utilizando:

```sh
./network.sh up createChannel -ca
```

Si ya está corriendo, revisaremos que el chaincode esté instalado y aprobado en el canal titleschannel.

### 6.3 Registrar el usuario `appUser`
El test requiere una identidad (appUser) en la wallet. Si no está registrada, necesitamos inscribirla. Usamos el siguiente comando en el directorio donde configuramos la red:

```sh
node enrollAdmin.js
node registerUser.js
```

Estos scripts deberían estar en la carpeta `fabric-samples/asset-transfer-basic/application-javascript/`. Si no los tienes, crea un nuevo script que inscriba un usuario en la CA.

### 6.4 Verificar el Archivo connection-org1.json
El test carga `connection-org1.json`, que debe contener la configuración de conexión a Fabric. Asegúrate de que el archivo existe en el mismo directorio que `app.test.js`. Si no, cópialo desde `fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json`

### 6.5 Ejecutar la Prueba
Ahora puedes ejecutar la prueba con:

```sh
node app.test.js
```
Si todo está configurado correctamente, deberías ver mensajes de éxito en la consola.

### 6.6 Depuración en Caso de Error
Si encuentras errores como:

* "El usuario appUser debe estar registrado e inscrito antes de ejecutar pruebas"
    * Ejecuta `node registerUser.js` nuevamente.

* "Error: Unable to connect to Fabric network"
    * Verifica que Fabric esté corriendo (`docker ps`).

* "Contract not found"
    * Asegúrate de que el chaincode academic-titles está desplegado en titleschannel.

## 7. Script automático
Para automatizar y simpliﬁcar el desarrollo de las pruebas, hemos generado un script en bash que automatiza la conﬁguración del entorno y las diferentes pruebas a realizar sobre el smartContract desarrollado (`test-title.sh`).

Podemos ejecutar este script de la siguiente forma:
```bash
# Dar permisos de ejecución
chmod +x test-title.sh

# Navegar al directorio de test-network
cd ~/hyperledger/fabric/fabric-samples/test-network

# Copiar el script aquí para facilitar el uso
cp ~/academic-titles-project/scripts/test-title.sh .

# Ejecutar pruebas
./test-title.sh emit TITLE002 STUDENT002 "Carlos Gómez" "Medicina" "2025-
03-06"
./test-title.sh verify TITLE002
```

## 8. Automatización de la generación de nodos ficticios utilizando Docker
Para simplificar y optimizar la generación de los nodos de la red ficticios utilizando Docker, hemos decidido implementar un script en bash que genera y despliega los contenedores de forma automática.