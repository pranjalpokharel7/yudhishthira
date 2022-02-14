package object

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Object struct {
	Name             []byte
	UniquePhysicalID []byte
	BasePrice        uint64
	Brand            []byte
	Category         string
}

type ObjectData struct {
	ObjectHash       []byte
	UniquePhysicalID []byte
}

func (objectInstance *Object) LoadCSVData(objectFilePath string) error {
	file, err := os.Open(objectFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)

	// read and skip first line
	if _, err = csvReader.Read(); err != nil {
		return err
	}

	objectData, err := csvReader.Read()
	if err != nil {
		return err
	}
	objectInstance.Name = []byte(objectData[0])
	objectInstance.UniquePhysicalID = []byte(objectData[1])
	objectInstance.BasePrice, err = strconv.ParseUint(objectData[2], 10, 64) // base 10, uint64
	objectInstance.Brand = []byte(objectData[3])
	objectInstance.Category = objectData[4]

	return err // the remaining error is from base price parsing
}

func (objectInstance *Object) HashObject() []byte {
	var hashBuffer bytes.Buffer

	hashBuffer.Write(objectInstance.Name)
	hashBuffer.Write(objectInstance.UniquePhysicalID)
	hashBuffer.Write(objectInstance.Brand)
	hashBuffer.Write([]byte(objectInstance.Category))
	binary.LittleEndian.PutUint64(hashBuffer.Bytes(), objectInstance.BasePrice)

	objectHash := sha256.Sum256(hashBuffer.Bytes())
	return objectHash[:]
}

func (objectInstance *Object) VerifyObjectHash(providedHash []byte) bool {
	objectHash := objectInstance.HashObject()
	return bytes.Equal(objectHash[:], providedHash)
}

func (objectInstance *Object) String() string {
	var objectDetails []string
	objectDetails = append(objectDetails, "---------- Object ----------")
	objectDetails = append(objectDetails, fmt.Sprintf("Name: %s", objectInstance.Name))
	objectDetails = append(objectDetails, fmt.Sprintf("Unique Physical ID: %s", objectInstance.UniquePhysicalID))
	objectDetails = append(objectDetails, fmt.Sprintf("Base Price: %d", objectInstance.BasePrice))
	objectDetails = append(objectDetails, fmt.Sprintf("Brand: %s", objectInstance.Brand))
	objectDetails = append(objectDetails, fmt.Sprintf("Category: %s\n", objectInstance.Category))
	return strings.Join(objectDetails, "\n")
}
