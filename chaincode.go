package chaincodetest

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AcademicTitle struct {
	TitleID        string `json:"titleId"`
	StudentID      string `json:"studentId"`
	StudentName    string `json:"studentName"`
	Degree         string `json:"degree"`
	EmissionDate   string `json:"emissionDate"`
	ValidationHash string `json:"validationHash"`
}

type TitleContract struct {
	contractapi.Contract
}

func (tc *TitleContract) IssueTitleToStudent(
	ctx contractapi.TransactionContextInterface,
	titleData AcademicTitle) error {

	// Validar que solo la universidad pueda emitir títulos (usando MSP)
	mspID, err := ctx.GetClientIdentity().GetMSPID()	
	if err != nil {
		return fmt.Errorf("ERROR [mientras se intentaba obtener el MSP ID]: %v", err)
	}
	if mspID != "universityMSP" {
		return fmt.Errorf("ERROR: Solo la universidad puede emitir títulos")
	}

	// Generar hash de validación
	titleData.ValidationHash = generateValidationHash(titleData)

	// Almacenar título
	titleBytes, err := json.Marshal(titleData)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(titleData.TitleID, titleBytes)
}

func (tc *TitleContract) VerifyTitle(
	ctx contractapi.TransactionContextInterface,
	titleID string) (*AcademicTitle, error) {
	titleBytes, err := ctx.GetStub().GetState(titleID)
	if err != nil {
		return nil, err
	}

	var title AcademicTitle
	err = json.Unmarshal(titleBytes, &title)
	return &title, err
}

func (tc *TitleContract) Transfer(
    ctx contractapi.TransactionContextInterface,
    tokenID string,
    newOwner string,
) error {
    return fmt.Errorf("ERROR: Soulbound Tokens no son transferibles")
}

func (tc *TitleContract) GetTitlesByStudent(
	ctx contractapi.TransactionContextInterface,
	studentID string) ([]*AcademicTitle, error) {
		queryString := fmt.Sprintf(`{"selector":{"studentId":"%s"}}`, studentID)

		// Creamos un iterador para buscar en el ledger
		resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
		if err != nil {
			return nil, fmt.Errorf("ERROR al obtener títulos: %v", err)
		}
		defer resultsIterator.Close()

		var titles []*AcademicTitle

		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return nil, fmt.Errorf("ERROR al iterar resultados: %v", err)
			}

			var title AcademicTitle
			err = json.Unmarshal(queryResponse.Value, &title)
			if err != nil {
				return nil, fmt.Errorf("ERROR al deserializar título: %v", err)
			}

			titles = append(titles, &title)
		}

		return titles, nil
}

func generateValidationHash(title AcademicTitle) string {
	// Implementar lógica de generación de hash
	// Por ejemplo, hash SHA-256 de los datos del título
	hash := sha256.Sum256([]byte(title.StudentID + title.Degree + title.EmissionDate))

	// Convertir el hash a una cadena hexadecimal
	return fmt.Sprintf("%x", hash)
}
