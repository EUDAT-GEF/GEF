# jOCCI-create-resources
This is a Java library to render <a href="http://occi-wg.org/about/specification/">Open Cloud Computing Interface (OCCI)</a> queries.

Detailed documentation is available in the project <a href="https://github.com/EGI-FCTF/jOCCI-api/wiki">wiki</a>.


## Compile and Run

Access the maven project

```cd di4r-training/jOCCI-create-resources/```

Edit your settings in the `src/main/java/it/infn/ct/Exercise4.java` source code to create a new ```compute``` resource:
```
[..]
String OCCI_ENDPOINT_HOST = "https://carach5.ics.muni.cz:11443"; // <= Change here!
String PROXY_PATH = "/tmp/x509up_u1000"; // <= Change here!

// *CREATE* a new virtual appliance (aka VM) with contextualization (public_key)
List<String> RESOURCE = Arrays.asList("compute");

List<String> MIXIN =
Arrays.asList("resource_tpl#medium", // <= Change here!
"http://occi.carach5.ics.muni.cz/occi/infrastructure/os_tpl#uuid_training_centos_6_fedcloud_warg_168"); // <= Change here!

List<String> CONTEXT =
Arrays.asList("public_key=file:/home/userX/.ssh/id_rsa.pub", // <= Change here!
"user_data=file:/home/userX/di4r-training/jOCCI-create-resources/contextualization.txt"); // <= Change here!

List<String> ATTRIBUTES = Arrays.asList("occi.core.title=VM_title"); // <= Change here!

String OCCI_PUBLICKEY_NAME = "centos";
```


Edit your settings in the `src/main/java/it/infn/ct/Exercise4.java` source code to create a new ```storage``` resource:
```
[..]
String OCCI_ENDPOINT_HOST = "https://carach5.ics.muni.cz:11443"; // <= Change here!
String PROXY_PATH = "/tmp/x509up_u1000"; // <= Change here!

// *CREATE* a new block storage
List<String> RESOURCE = Arrays.asList("storage"); 

public static List<String> ATTRIBUTES = 
    Arrays.asList("occi.core.title=VM_volume_1", "occi.storage.size=1"); // <= Change here!
    
// Set to null not used variables.
String OCCI_PUBLICKEY_NAME = "";
List<String> CONTEXT = new ArrayList<String>();
List<String> MIXIN = new ArrayList<String>();
```

Compile and package with maven:
```
$ mvn compile && mvn package
```

Run (you may redirect the output to a file):
```
$ java –jar target/jocci-create-resource-1.0-jar-with-dependencies.jar
```

## Dependencies

jOCCI-create-resources uses:
- jocci-api (v0.2.5)
- slf4j-jdk14 (v1.7.12)

These are already included in the Maven pom.xml file and automatically downloaded when building.

You can also add them to your projects with:

```
    <dependency>
        <groupId>org.slf4j</groupId>
        <artifactId>slf4j-jdk14</artifactId>
        <version>1.7.12</version>
    </dependency>

    <dependency>
        <groupId>cz.cesnet.cloud</groupId>
        <artifactId>jocci-api</artifactId>
        <version>0.2.5</version>
        <scope>compile</scope>
    </dependency>
```
