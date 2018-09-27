test:
	rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro schemas/*.avsc && go run main.go

idl:
	rm -rf output && mkdir output && java -jar /Users/kilia/Downloads/avro-tools-1.8.2.jar idl2schemata idl/datatypes.identifier.avdl output

final:
	rm -rf output && mkdir output && java -jar /Users/kilia/Downloads/avro-tools-1.8.2.jar idl2schemata idl/datatypes.identifier.avdl output && rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro output/Identifier.avsc && go run main.go

output:
	rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro output/*.avsc && go run main.go