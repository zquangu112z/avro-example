JAVRO=java -jar /Users/kilia/Downloads/avro-tools-1.8.2.jar
JAVRO_OUTPUT=./schema-avro
GOGEN=gogen-avro -containers 
GOGEN_OUTPUT=./models
GOGEN_PACKAGE=models

SERIALIZE=go run main.go

MAIN_IDL=./schema-avdl/main.avdl
MAIN_AVRO=$(JAVRO_OUTPUT)/Main.avsc 

idl:
	rm -f $(JAVRO_OUTPUT)/* && $(JAVRO) idl2schemata $(MAIN_IDL) $(JAVRO_OUTPUT)

output:
	rm -f $(GOGEN_OUTPUT)/* && mkdir $(GOGEN_OUTPUT) && $(GOGEN) $(GOGEN_PACKAGE) $(MAIN_AVRO) && $(SERIALIZE)

test:
	rm -f $(JAVRO_OUTPUT)/* && $(JAVRO) idl2schemata $(MAIN_IDL) $(JAVRO_OUTPUT) && rm -f $(GOGEN_OUTPUT)/* && mkdir $(GOGEN_OUTPUT) &&  $(GOGEN) $(MAIN_AVRO) $(MAIN_AVRO) && $(SERIALIZE)

clean:
	rm -f $(GOGEN_OUTPUT)/* && rm -f $(JAVRO_OUTPUT)/*

tojson:
	$(JAVRO) tojson observation.avro > obx.json