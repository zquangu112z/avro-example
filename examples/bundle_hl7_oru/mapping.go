package main

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
	mapping                         = [][]string{
		// Patient
		{">PID.3.0.0.0", "Patient.Identifier.Value"},
		{">PID.5.0.1.0", "Patient.HumanName.Given"},
		{">PID.5.0.0.0", "Patient.HumanName.Family"},

		// Observation
		{">OBR.21.0.0.0", "Observation.Issued"},
		{"OBX.3.0.0.0", "Observation.Code"},
	}
)
