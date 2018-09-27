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
	observation := avro.Observation{
		Identifier:           nil,
		BaseOn:               nil,
		Status:               "fine",
		Category:             nil,
		Code:                 nil,
		Subject:              nil,
		Context:              nil,
		EffectiveDateTime:    "today",
		EffectivePeriod:      nil,
		Issued:               "ABX",
		Performer:            nil,
		ValueQuantity:        nil,
		ValueCodeableConcept: nil,
		ValueString:          "nil",
		ValueBoolean:         "nil",
		ValueRange:           nil,
		ValueTime:            "nil",
		ValueDateTime:        "nil",
		ValuePeriod:          nil,
	}
	// fmt.Println(observation.Schema())
	fmt.Println("-----------")

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
