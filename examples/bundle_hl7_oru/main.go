package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/actgardner/gogen-avro/container"
	avro "github.com/zquangu112z/avro-example/models"
	hl7 "github.com/zquangu112z/hl7-parser"
)

var (
	message = []byte(`MSH|^~\&|Epic|SIH|PH Application^2.16.840.1.113883.3.72.7.3^HL7|PH facility^2.16.840.1.113883.3.72.7.4^HL7|20171028000018|LABBACKGROUND|ORU^R01^ORU_R01|669115|P|2.5.1|||||||||LRI_Common_Component^^2.16.840.1.113883.9.16^ISO~LRI_NG_Component^^2.16.840.1.113883.9.13^ISO~LRI_RN_Component^^2.16.840.1.113883.9.15^ISO||
	SFT|Epic Systems Corporation^L^^^^ANSI&1.2.840&ISO^XX^^^1.2.840.114350|Epic 2015 |Bridges|8.2.2.0||20151217091530
	PID|1||1071782^^^SIHEPIC^MR||Sargent^Dakota^Aaron^^^^D||19970910|M||2106-3|400 S 23RD ST^^HERRIN^IL^62948-2116^USA^L^^WILLIAMSON|||||||511452378|333-94-5527|||N||||||||N|||20170803143227|1|
	PD1|||SIH HERRIN HOSPITAL^^102001|1548261449^CHURCH^SHOSHANA^^^^^^NPI^^^^NPI||||||||||||||
	NK1|1|ULMER^GERMANE^^|GRD|400 S 23RD ST^^HERRIN^IL^62948-2116^USA^^^WILLIAMSON|(618)727-0645^^H^^^618^7270645|||||||||||||||||||||||||||
	PV1|1|E|HH ED^HER11^11^HH^^^102001^SIH Herrin Hospital^HH EMERGENCY DEPARTMENT^^|ER|||1760764153^MASON^BENJAMIN^D^^^^^NPI^^^^NPI|||ER|||||||||200000081838|||||||||||||||||||||||||20171027223500|
	ORC|RE|13817799^|17H-300C0402^Beaker|17H-300C0402^Beaker||||||||1760764153^MASON^BENJAMIN^D^^^^^NPI^^^^NPI||(618)529-0721^^^^^618^5290721||||||||||MHC EMERGENCY DEPARTMENT 405 WEST JACKSON^^CARBONDALE^IL^62901^^C|||||||
	OBR|1|13817799^|17H-300C0402^Beaker|LAB3002^Drug Screen, Urine^SIH_EAP^^^^^^DRUG SCREEN, URINE|||20171027234100||||Unit Collect|||||1760764153^MASON^BENJAMIN^D^^^^^NPI^^^^NPI|(618)529-0721^^^^^618^5290721|||||20171028000000|||F|||||||&Baker&Christopher&&||||||||||||||||||
	NTE|1|L|This assay provides a preliminary analytical test result.  A more specific, alternate chemical method must be ordered to obtain a confirmed analytical result.|
	TQ1|1||||||20171027223900|20171027235959|S
	OBX|1|ST|1534341^Amphetamine^SIH_LRR^14308-1^Amphetamines Ur Ql Scn>1000 ng/mL^LN^^^AMPHETAMINE, CUT-OFF 1000 NG/ML||Detected||Cut-off 1000 ng/mL|A|||F|||20171027234100||CHBAKER^BAKER^CHRISTOPHER^^|||20171028000013||||HERRIN HOSPITAL^D|201 SOUTH 14TH STREET^^HERRIN^IL^62948^^B|1801826425^GONZALEZ^JUAN^G^^^^^NPI^^^^NPI
	OBX|2|ST|1534345^Barbituate^SIH_LRR^70155-7^Barbiturates Ur Ql Scn>200 ng/mL^LN^^^BARBITUATE, CUT-OFF 300 NG/ML||Undetected||Cut-off 300 ng/mL||||F|||20171027234100||CHBAKER^BAKER^CHRISTOPHER^^|||20171028000013||||HERRIN HOSPITAL^D|201 SOUTH 14TH STREET^^HERRIN^IL^62948^^B|1801826425^GONZALEZ^JUAN^G^^^^^NPI^^^^NPI
	OBX|3|ST|1558232^Benzodiazepine^SIH_LRR^3390-2^Benzodiaz Ur Ql^LN^^^BENZODIAZEPINE, CUT-OFF 300 NG/ML||Undetected||Cut-off 300 ng/mL||||F|||20171027234100||CHBAKER^BAKER^CHRISTOPHER^^|||20171028000013||||HERRIN HOSPITAL^D|201 SOUTH 14TH STREET^^HERRIN^IL^62948^^B|1801826425^GONZALEZ^JUAN^G^^^^^NPI^^^^NPI
	SPM|1|||^^^^^^^^Urine||||^^^^^^^^Urine, Voided|||||||||20171027234100|20171027234633||||||
	`)
	generalHeaders, specificHeaders []string
	generalNames                    []string
	mappingSection                  [][]string

	mapping = [][]string{
		// Patient
		{">PID.3.0.0.0", "Patient.Identifier.Value"},
		{">PID.5.0.1.0", "Patient.HumanName.Given"},
		{">PID.5.0.0.0", "Patient.HumanName.Family"},

		// Observation
		{">OBR.21.0.0.0", "Observation.Issued"},
		{"OBX.3.0.0.0", "Observation.Code"},
	}
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

// Return 2 lists of header: one for general and one for specific data (section)
func getHeaderLists(mapping [][]string) ([]string, []string) {
	var generalHeaders []string
	var specificHeaders []string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], ">") {
			generalHeaders = appendIfMissing(generalHeaders, tuple[1])
		} else {
			specificHeaders = appendIfMissing(specificHeaders, tuple[1])
		}
	}
	return generalHeaders, specificHeaders
}

// Return 2 lists of segments' name: one for general and one for specific data (section)
func getNameLists(mapping [][]string) ([]string, string) {
	var generalNames []string
	var specificName string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], ">") {
			newItem := strings.TrimPrefix(strings.Split(tuple[0], ".")[0], ">")
			generalNames = appendIfMissing(generalNames, newItem)
		} else {
			if len(specificName) == 0 {
				newItem := strings.Split(tuple[0], ".")[0]
				specificName = newItem
			}
		}
	}
	return generalNames, specificName
}

// Return the mapping for section only (not include general information)
func getMappingSection(mapping [][]string, specificName string) [][]string {
	var mappingSpecific [][]string
	for _, tuple := range mapping {
		if strings.HasPrefix(tuple[0], specificName) {
			mappingSpecific = append(mappingSpecific, tuple)
		}
	}
	return mappingSpecific
}

// Return a slice with unique items
func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
