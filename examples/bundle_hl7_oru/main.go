package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/actgardner/gogen-avro/container"
	avro "github.com/zquangu112z/avro-example/models"
	hl7 "github.com/zquangu112z/hl7-parser"
)

func convertToFhirBundle(message []byte) (*avro.Bundle, error) {
	// Extract messages from file. Each message has several line
	Hl7Message, err := hl7.Decode(message)
	if err != nil {
		panic(err)
	}

	// >>> Bundle
	bundle := avro.NewBundle()

	// >>> MessageHeader
	messageHeader := avro.NewMessageHeader()
	mshSegment := Hl7Message[0] // NOTE: hereby, assume the first segment is a MSH segment
	createMessageHeader(mshSegment, messageHeader)

	messageHeaderBundleEntry := avro.NewBundleEntry()
	messageHeaderBundleEntryResource := avro.UnionObservationPatientMessageHeaderDiagnosticReport{
		MessageHeader: messageHeader,
		UnionType:     avro.UnionObservationPatientMessageHeaderDiagnosticReportTypeEnumMessageHeader,
	}
	(*messageHeaderBundleEntry).Resource = messageHeaderBundleEntryResource

	// MessageHeader is the first entry in the bundle: http://www.hl7.org/implement/standards/fhir/messageheader.html
	(*bundle).Entry = append((*bundle).Entry, messageHeaderBundleEntry)

	generalDataDict := make(map[string]string)

	for _, segment := range Hl7Message[1:] {
		segmentName := segment.AtIndex("0.0.0.0")
		switch segmentName {
		case "OBX": // Found Observation
			// Parse OBX data
			obxDict := make(map[string]string)
			for _, tuple := range mappingSection {
				path := tuple[0]
				key := tuple[1]

				obxDict[key] = segment.AtIndex(strings.SplitN(path, ".", 2)[1])
			}
			// >>> Patient
			patient := avro.NewPatient()

			patientIdentifiers := []*avro.Identifier{}
			// only get 1 Identifier
			patientIdentifier := avro.NewIdentifier()
			patientIdentifier.Value = generalDataDict["Patient.Identifier.Value"]

			patientIdentifiers = append(patientIdentifiers, patientIdentifier)

			patientNames := []*avro.HumanName{}
			// only get 1 HumanName
			patientName := avro.NewHumanName()
			patientName.Given = generalDataDict["Patient.HumanName.Given"]
			patientName.Family = generalDataDict["Patient.HumanName.Family"]

			patientNames = append(patientNames, patientName)

			patient.Identifier = patientIdentifiers
			patient.Name = patientNames

			// @TODO: Send/Check if exist patient
			writePatient(patient)

			observation := avro.NewObservation()
			observation.Issued = obxDict["Observation.Issued"]
			observation.Subject = &avro.Reference{
				Reference: (*patient).Identifier[0].Value,
				Display:   (*patient).Identifier[0].Value,
			}

			observationBundleEntry := avro.NewBundleEntry()
			observationBundleEntryResource := avro.UnionObservationPatientMessageHeaderDiagnosticReport{
				Observation: observation,
				UnionType:   avro.UnionObservationPatientMessageHeaderDiagnosticReportTypeEnumObservation,
			}
			observationBundleEntry.Resource = observationBundleEntryResource

			// Add to bundle
			(*bundle).Entry = append((*bundle).Entry, observationBundleEntry)

		default:
			for _, generalName := range generalNames {
				if segmentName == generalName { // need to update the generalDataDict
					updateGeneralDict(segment, segmentName, generalDataDict)
					break
				}
			}
		}
	}

	return bundle, nil
}

// In some case, the general segment appear more than 1 time. In that case, we need to update to general dict
func updateGeneralDict(segment hl7.Hl7Segment, segmentName string, generalDataDict map[string]string) {
	for _, tuple := range mapping {
		path := tuple[0]
		if strings.Contains(path, segmentName) == false {
			continue
		}
		key := tuple[1]
		actualPath := strings.SplitN(path, ".", 2)[1]
		data := segment.AtIndex(actualPath)
		generalDataDict[key] = data
	}
}

func main() {
	// Extract information from mapping
	generalHeaders, specificHeaders = getHeaderLists(mapping)
	generalNames, _ = getNameLists(mapping)
	mappingSection = getMappingSection(mapping, "OBX")

	bundle, err := convertToFhirBundle(message)
	if err != nil {
		fmt.Println(err)
	}
	writeBundle(bundle)
}

func writeBundle(bundle *avro.Bundle) {
	// Open a file to write
	fileWriter, err := os.Create("bundle.avro")
	if err != nil {
		fmt.Printf("Error opening file writer: %v\n", err)
		return
	}

	containerWriter, err := avro.NewBundleWriter(fileWriter, container.Null, 10)
	if err != nil {
		fmt.Printf("Error opening container writer: %v\n", err)
		return
	}

	// Write the record to the container file
	err = containerWriter.WriteRecord(bundle)
	if err != nil {
		fmt.Printf("Error writing record to file: %v\n", err)
		return
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		fmt.Printf("Error flushing last block to file: %v\n", err)
		return
	}

	fileWriter.Close()
}

func writePatient(patient *avro.Patient) {
	// Open a file to write
	fileWriter, err := os.Create("patient.avro")
	if err != nil {
		fmt.Printf("Error opening file writer: %v\n", err)
		return
	}

	containerWriter, err := avro.NewPatientWriter(fileWriter, container.Null, 10)
	if err != nil {
		fmt.Printf("Error opening container writer: %v\n", err)
		return
	}

	// Write the record to the container file
	err = containerWriter.WriteRecord(patient)
	if err != nil {
		fmt.Printf("Error writing record to file: %v\n", err)
		return
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		fmt.Printf("Error flushing last block to file: %v\n", err)
		return
	}

	fileWriter.Close()
}

func createMessageHeader(mshSegment hl7.Hl7Segment, messageHeader *avro.MessageHeader) error {
	event := avro.NewCoding()
	(*event).Code = "observation-provide" // string(mshSegment[8][0][1][0]),
	(*event).Display = "observation-provide"
	(*event).System = "http://hl7.org/fhir/ValueSet/message-events"

	(*messageHeader).Event = event

	messageHeaderDestination := avro.NewMessageHeaderDestination()
	(*messageHeaderDestination).Name = string(mshSegment[4][0][0][0])
	(*messageHeaderDestination).Endpoint = string(mshSegment[24][0][0][0])

	(*messageHeader).Destination = messageHeaderDestination

	(*messageHeader).Timestamp = string(mshSegment[6][0][0][0])

	messageHeaderSource := avro.NewMessageHeaderSource()
	(*messageHeaderSource).Name = string(mshSegment[2][0][0][0])
	(*messageHeaderSource).Software = string(mshSegment[2][0][0][0])
	(*messageHeaderSource).Endpoint = string(mshSegment[23][0][0][0])

	(*messageHeader).Source = messageHeaderSource

	return nil
}

func createPatient(patientSegment hl7.Hl7Segment, patient *avro.Patient) error {
	return nil
}

func createObservation(patientSegment hl7.Hl7Segment, observation *avro.Observation, patient avro.Patient, provider string) error {
	return nil
}
