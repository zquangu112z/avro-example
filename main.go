package main

import (
	"fmt"
	"os"

	"github.com/actgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"github.com/zquangu112z/avro-example/avro"
)

var (
	codec *goavro.Codec
)

func main() {

	//Sample Data
	coding := avro.Coding{
		System:       "okie",
		Version:      "okie",
		Code:         "okie",
		Display:      "okie",
		UserSelected: true,
	}
	// var coding_array []*avro.Coding
	// coding_array = append(coding_array, &coding)
	codeableConcept := avro.CodeableConcept{
		Text:   "hihi",
		Coding: []*avro.Coding{&coding},
	}

	reference := avro.Reference{
		Reference: "ref1",
		Display:   "Dis1",
	}
	period := avro.Period{
		Start: "1",
		End:   "2",
	}

	quantity := avro.Quantity{
		Value:      6.9,
		Comparator: "placeholader",
		Unit:       "placeholader",
		System:     "placeholader",
		Code:       "placeholader",
	}

	arange := avro.Range{
		Low:  &quantity,
		High: &quantity,
	}

	observation := avro.Observation{
		Identifier:           nil,
		BaseOn:               []*avro.Reference{&reference},
		Status:               "fine",
		Category:             &codeableConcept,
		Code:                 &codeableConcept,
		Subject:              &reference,
		Context:              &reference,
		EffectiveDateTime:    "today",
		EffectivePeriod:      &period,
		Issued:               "ABX",
		Performer:            &reference,
		ValueQuantity:        &quantity,
		ValueCodeableConcept: &codeableConcept,
		ValueString:          "nil",
		ValueBoolean:         "nil",
		ValueRange:           &arange,
		ValueTime:            "nil",
		ValueDateTime:        "nil",
		ValuePeriod:          &period,
	}

	// fmt.Println(observation.Schema())

	// Open a file to write
	fileWriter, err := os.Create("observation.avro")
	if err != nil {
		fmt.Printf("Error opening file writer: %v\n", err)
		return
	}

	containerWriter, err := avro.NewObservationWriter(fileWriter, container.Null, 10)
	if err != nil {
		fmt.Printf("Error opening container writer: %v\n", err)
		return
	}

	// Write the record to the container file
	err = containerWriter.WriteRecord(&observation)
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
