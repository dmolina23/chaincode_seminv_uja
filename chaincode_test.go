package chaincodetest

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

import "github.com/hyperledger/fabric-protos-go/ledger/queryresult"


// MockStub simula el stub de transacción
type MockStub struct {
	mock.Mock
	shim.ChaincodeStubInterface
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)
	return args.Error(0)
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

// MockClientIdentity simula la identidad del cliente
type MockClientIdentity struct {
	mock.Mock
}

func (mci *MockClientIdentity) GetMSPID() (string, error) {
	args := mci.Called()
	return args.String(0), args.Error(1)
}

func (mci *MockClientIdentity) GetID() (string, error) {
	args := mci.Called()
	return args.String(0), args.Error(1)
}

func (mci *MockClientIdentity) GetAttributeValue(attrName string) (string, bool, error) {
	args := mci.Called(attrName)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (mci *MockClientIdentity) AssertAttributeValue(attrName, attrValue string) error {
	args := mci.Called(attrName, attrValue)
	return args.Error(1)
}

// GetX509Certificate implements cid.ClientIdentity
func (mci *MockClientIdentity) GetX509Certificate() (*x509.Certificate, error) {
	args := mci.Called()
	if cert, ok := args.Get(0).(*x509.Certificate); ok {
		return cert, args.Error(1)
	}
	return nil, args.Error(1)
}

// MockContext simula el contexto de transacción
type MockContext struct {
	mock.Mock
	contractapi.TransactionContextInterface
	stub     *MockStub
	clientID *MockClientIdentity
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	return mc.stub
}

func (mc *MockContext) GetClientIdentity() cid.ClientIdentity {
	return mc.clientID
}

// Pruebas unitarias
func TestIssueTitleToStudent(t *testing.T) {
	// Configurar mocks
	mockStub := new(MockStub)
	mockClientID := new(MockClientIdentity)
	mockContext := new(MockContext)
	mockContext.stub = mockStub
	mockContext.clientID = mockClientID

	contract := new(TitleContract)

	// Datos de prueba
	testTitle := AcademicTitle{
		TitleID:      "TITLE001",
		StudentID:    "STUDENT001",
		StudentName:  "Juan Pérez",
		Degree:       "Ingeniería Informática",
		EmissionDate: "2025-03-05",
	}

	// Test case: Usuario con permisos correctos
	t.Run("Emisión exitosa con permisos correctos", func(t *testing.T) {
		// Configurar expectativas para GetMSPID en lugar de HasRole
		mockClientID.On("GetMSPID").Return("universityMSP", nil)
		mockStub.On("PutState", testTitle.TitleID, mock.Anything).Return(nil)

		// Ejecutar función
		err := contract.IssueTitleToStudent(mockContext, testTitle)

		// Verificar resultados
		assert.NoError(t, err)
		mockClientID.AssertExpectations(t)
		mockStub.AssertExpectations(t)
	})

	// Test case: Usuario sin permisos
	t.Run("Emisión fallida sin permisos", func(t *testing.T) {
		// Reset mocks
		mockClientID = new(MockClientIdentity)
		mockStub = new(MockStub)
		mockContext.clientID = mockClientID
		mockContext.stub = mockStub

		// Configurar expectativas para GetMSPID con un valor incorrecto
		mockClientID.On("GetMSPID").Return("studentMSP", nil)

		// Ejecutar función
		err := contract.IssueTitleToStudent(mockContext, testTitle)

		// Verificar resultados
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Solo la universidad puede emitir títulos")
		mockClientID.AssertExpectations(t)
		// No debería llegar a PutState
		mockStub.AssertNotCalled(t, "PutState")
	})

	// Test case: Error al obtener MSP ID
	t.Run("Error al obtener MSP ID", func(t *testing.T) {
		// Reset mocks
		mockClientID = new(MockClientIdentity)
		mockStub = new(MockStub)
		mockContext.clientID = mockClientID
		mockContext.stub = mockStub

		// Configurar expectativas para GetMSPID con error
		mockClientID.On("GetMSPID").Return("", fmt.Errorf("error MSP ID"))

		// Ejecutar función
		err := contract.IssueTitleToStudent(mockContext, testTitle)

		// Verificar resultados
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR [mientras se intentaba obtener el MSP ID]")
		mockClientID.AssertExpectations(t)
		// No debería llegar a PutState
		mockStub.AssertNotCalled(t, "PutState")
	})
}

func TestVerifyTitle(t *testing.T) {
	// Configurar mocks
	mockStub := new(MockStub)
	mockClientID := new(MockClientIdentity)
	mockContext := new(MockContext)
	mockContext.stub = mockStub
	mockContext.clientID = mockClientID

	contract := new(TitleContract)

	// Datos de prueba
	testTitle := AcademicTitle{
		TitleID:      "TITLE001",
		StudentID:    "STUDENT001",
		StudentName:  "Juan Pérez",
		Degree:       "Ingeniería Informática",
		EmissionDate: "2025-03-05",
		ValidationHash: generateValidationHash(AcademicTitle{
			StudentID:    "STUDENT001",
			Degree:       "Ingeniería Informática",
			EmissionDate: "2025-03-05",
		}),
	}

	titleBytes, _ := json.Marshal(testTitle)

	// Test case: Verificación exitosa
	t.Run("Verificación exitosa de título existente", func(t *testing.T) {
		mockStub.On("GetState", "TITLE001").Return(titleBytes, nil)

		title, err := contract.VerifyTitle(mockContext, "TITLE001")

		assert.NoError(t, err)
		assert.Equal(t, testTitle.StudentID, title.StudentID)
		assert.Equal(t, testTitle.ValidationHash, title.ValidationHash)
		mockStub.AssertExpectations(t)
	})

	// Test case: Título no encontrado
	t.Run("Verificación de título inexistente", func(t *testing.T) {
		// Reset mock
		mockStub = new(MockStub)
		mockContext.stub = mockStub

		// Devolvemos un empty byte array para simular que no existe
		mockStub.On("GetState", "NONEXISTENT").Return([]byte{}, nil)

		_, err := contract.VerifyTitle(mockContext, "NONEXISTENT")

		// El error ocurrirá en el unmarshal de un JSON vacío
		assert.Error(t, err)
		// title no será nil, pero será un struct vacío después del error unmarshal
		mockStub.AssertExpectations(t)
	})

	// Test case: Error al obtener estado
	t.Run("Error al obtener estado", func(t *testing.T) {
		// Reset mock
		mockStub = new(MockStub)
		mockContext.stub = mockStub

		mockStub.On("GetState", "ERROR_TITLE").Return([]byte{}, fmt.Errorf("error al leer del ledger"))

		title, err := contract.VerifyTitle(mockContext, "ERROR_TITLE")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error al leer del ledger")
		assert.Nil(t, title)
		mockStub.AssertExpectations(t)
	})
}
func TestTransfer(t *testing.T) {
	// Configurar mocks
	mockStub := new(MockStub)
	mockClientID := new(MockClientIdentity)
	mockContext := new(MockContext)
	mockContext.stub = mockStub
	mockContext.clientID = mockClientID

	contract := new(TitleContract)

	// Test case: Intento de transferencia
	t.Run("Intento de transferencia de Soulbound Token", func(t *testing.T) {
		err := contract.Transfer(mockContext, "TITLE001", "NEW_OWNER")

		// Verificar resultados
		assert.Error(t, err)
		assert.Equal(t, "ERROR: Soulbound Tokens no son transferibles", err.Error())
	})
}

func TestGetTitlesByStudent(t *testing.T) {
	// Configurar mocks
	mockStub := new(MockStub)
	mockClientID := new(MockClientIdentity)
	mockContext := new(MockContext)
	mockContext.stub = mockStub
	mockContext.clientID = mockClientID

	contract := new(TitleContract)

	// Datos de prueba
	testTitles := []AcademicTitle{
		{
			TitleID:      "TITLE001",
			StudentID:    "STUDENT001",
			StudentName:  "Juan Pérez",
			Degree:       "Ingeniería Informática",
			EmissionDate: "2025-03-05",
			ValidationHash: generateValidationHash(AcademicTitle{
				StudentID:    "STUDENT001",
				Degree:       "Ingeniería Informática",
				EmissionDate: "2025-03-05",
			}),
		},
		{
			TitleID:      "TITLE002",
			StudentID:    "STUDENT001",
			StudentName:  "Juan Pérez",
			Degree:       "Máster en IA",
			EmissionDate: "2026-07-15",
			ValidationHash: generateValidationHash(AcademicTitle{
				StudentID:    "STUDENT001",
				Degree:       "Máster en IA",
				EmissionDate: "2026-07-15",
			}),
		},
	}

	// Test case: Obtener títulos exitosamente
	t.Run("Obtener títulos exitosamente", func(t *testing.T) {
		// Crear resultados simulados para el iterador
		mockResults := []shim.StateQueryIteratorInterface{
			&MockStateQueryIterator{
				Results: [][]byte{
					mustMarshal(testTitles[0]),
					mustMarshal(testTitles[1]),
				},
			},
		}

		mockStub.On("GetQueryResult", `{"selector":{"studentId":"STUDENT001"}}`).Return(mockResults[0], nil)

		titles, err := contract.GetTitlesByStudent(mockContext, "STUDENT001")

		assert.NoError(t, err)
		assert.Len(t, titles, 2)
		assert.Equal(t, testTitles[0].TitleID, titles[0].TitleID)
		assert.Equal(t, testTitles[1].TitleID, titles[1].TitleID)
		mockStub.AssertExpectations(t)
	})

	// Test case: Error al obtener resultados del ledger
	t.Run("Error al obtener resultados del ledger", func(t *testing.T) {
		mockStub.On("GetQueryResult", `{"selector":{"studentId":"STUDENT001"}}`).Return(nil, fmt.Errorf("error al ejecutar consulta"))

		titles, err := contract.GetTitlesByStudent(mockContext, "STUDENT001")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR al obtener títulos")
		assert.Nil(t, titles)
		mockStub.AssertExpectations(t)
	})

	// Test case: Error al iterar resultados
	t.Run("Error al iterar resultados", func(t *testing.T) {
		mockResults := []shim.StateQueryIteratorInterface{
			&MockStateQueryIterator{
				Results: [][]byte{
					mustMarshal(testTitles[0]),
				},
				NextError: fmt.Errorf("error al iterar"),
			},
		}

		mockStub.On("GetQueryResult", `{"selector":{"studentId":"STUDENT001"}}`).Return(mockResults[0], nil)

		titles, err := contract.GetTitlesByStudent(mockContext, "STUDENT001")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR al iterar resultados")
		assert.Nil(t, titles)
		mockStub.AssertExpectations(t)
	})

	// Test case: Error al deserializar título
	t.Run("Error al deserializar título", func(t *testing.T) {
		mockResults := []shim.StateQueryIteratorInterface{
			&MockStateQueryIterator{
				Results: [][]byte{
					[]byte("invalid-json"),
				},
			},
		}

		mockStub.On("GetQueryResult", `{"selector":{"studentId":"STUDENT001"}}`).Return(mockResults[0], nil)

		titles, err := contract.GetTitlesByStudent(mockContext, "STUDENT001")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR al deserializar título")
		assert.Nil(t, titles)
		mockStub.AssertExpectations(t)
	})
}

// Helper para serializar datos
func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

// MockStateQueryIterator simula un iterador de consulta de estado
type MockStateQueryIterator struct {
	mock.Mock
	Results   [][]byte
	NextError error
	Index     int
}

func (ms *MockStub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	args := ms.Called(query)
	return args.Get(0).(shim.StateQueryIteratorInterface), args.Error(1)
}

func (m *MockStateQueryIterator) HasNext() bool {
	return m.Index < len(m.Results)
}

func (m *MockStateQueryIterator) Next() (*queryresult.KV, error) {
	if m.NextError != nil {
		return nil, m.NextError
	}
	result := m.Results[m.Index]
	m.Index++
	return &queryresult.KV{Value: result}, nil
}

func (m *MockStateQueryIterator) Close() error {
	return nil
}