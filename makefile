JAVRO=java -jar ~/Downloads/avro-tools-1.8.2.jar
JAVRO_OUTPUT=./schema-avro
GOGEN=gogen-avro -containers 
GOGEN_OUTPUT=./models
GOGEN_PACKAGE=models

SERIALIZE=go run main.go

MAIN_IDL=./schema-avdl/main.avdl
MAIN_AVRO=$(JAVRO_OUTPUT)/Main.avsc 

idl:
	rm -rf $(JAVRO_OUTPUT) && mkdir $(JAVRO_OUTPUT) && $(JAVRO) idl2schemata $(MAIN_IDL) $(JAVRO_OUTPUT)

output:
	rm -rf $(GOGEN_OUTPUT) && mkdir $(GOGEN_OUTPUT) && $(GOGEN) $(GOGEN_PACKAGE) $(MAIN_AVRO) && $(SERIALIZE)

tojson:
	$(JAVRO) tojson observation.avro > obx.json

test:
	rm -rf $(JAVRO_OUTPUT) && mkdir $(JAVRO_OUTPUT) && $(JAVRO) idl2schemata $(MAIN_IDL) $(JAVRO_OUTPUT) &&rm -rf $(GOGEN_OUTPUT) && mkdir $(GOGEN_OUTPUT) && $(GOGEN) $(GOGEN_PACKAGE) $(MAIN_AVRO) && $(SERIALIZE) && $(JAVRO) tojson observation.avro > obx.json

clean:
	rm -rf $(GOGEN_OUTPUT) && rm -rf $(JAVRO_OUTPUT)

parser:
	go run examples/bundle_hl7_oru/main.go