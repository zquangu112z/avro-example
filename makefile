test:
	rm -rf avro &&  mkdir avro && ~/go/bin/gogen-avro --containers  avro schemas/*.avsc && go run main.go