package main

import (
	"fmt"

	"github.com/linkedin/goavro"
	"github.com/zquangu112z/avro-example/avro"
)

var (
	codec *goavro.Codec
)

func main() {

	//Sample Data
	bundle := avro.Observation{}
	fmt.Println(twitter_schema.Schema())

	// Open a file to write
	// fileWriter, err := os.Create("bundle.avro")
	// if err != nil {
	// 	fmt.Printf("Error opening file writer: %v\n", err)
	// 	return
	// }

	// containerWriter, err := avro.NewTwitter_schemaWriter(fileWriter, container.Null, 10)
	// if err != nil {
	// 	fmt.Printf("Error opening container writer: %v\n", err)
	// 	return
	// }

	// // Write the record to the container file
	// err = containerWriter.WriteRecord(&twitter_schema)
	// if err != nil {
	// 	fmt.Printf("Error writing record to file: %v\n", err)
	// 	return
	// }

	// // Flush the buffers to ensure the last block has been written
	// err = containerWriter.Flush()
	// if err != nil {
	// 	fmt.Printf("Error flushing last block to file: %v\n", err)
	// 	return
	// }

	// fileWriter.Close()

}
