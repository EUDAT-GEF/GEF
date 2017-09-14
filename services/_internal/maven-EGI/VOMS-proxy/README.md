# Api-VOMS-CANL

Api-VOMS-CANL is sample creation of a Virtual Organization Membership Service (VOMS) proxy using the VOMS client. The application contact the VOMS service to generate RFC or full-legacy X.509 proxy certificates.

## Compile and Run

Access the maven project

```cd di4r-training/Api-VOMS-CANL/```

Edit and specify your settigns in the ```src/main/java/it/infn/ct/VOMSProxyInit.java``` source code:

```
[..]
String VONAME = "training.egi.eu"; // <= Change here!
String VOMS_PROXY_FILEPATH = "/tmp/x509up_u1000"; // <= Change here!
String VOMS_LIFETIME = "24:00";
String VOMSES_DIR = "/etc/vomses/";
String X509_CERT_DIR = "/etc/grid-security/certificates/";
```

Compile and package with maven:

$ mvn compile && mvn package

Run:

$ java –jar target/jVOMS-Proxy-Init-1.0-jar-with-dependencies.jar


## Dependencies

Api-VOMS-CANL uses:
- voms-clients (v3.0.6)
- log4j (v1.2.17)

These are already included in the Maven pom.xml file and automatically downloaded when building.

You can also add them to your projects with:
```
<dependency>
    <groupId>org.italiangrid</groupId>
    <artifactId>voms-clients</artifactId>
    <version>3.0.6</version>
</dependency>

<dependency>
    <groupId>log4j</groupId>
    <artifactId>log4j</artifactId>
    <version>1.2.17</version>
</dependency>
```
