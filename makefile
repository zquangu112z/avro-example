idl:
	rm -rf output && mkdir output && java -jar /Users/kilia/Downloads/avro-tools-1.8.2.jar idl2schemata idl/main.avdl output

output:
	rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro output/Main.avsc && go run main.go

test:
	rm -rf output && mkdir output && java -jar /Users/kilia/Downloads/avro-tools-1.8.2.jar idl2schemata idl/main.avdl output && rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro output/Main.avsc && go run main.go

clean:
	rm -rf ./avro && rm -rf ./output