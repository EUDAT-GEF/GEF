#!/bin/bash
echo "******************************************************************************************"
echo "You are about to create certificates for the GEF. They will be saved at the current folder"
echo "******************************************************************************************"
echo "Follow the instructions!!!"
read -r -p "Are you sure you want to continue? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY])
        echo ""
        echo "STEP 1: generate CA private and public keys"
        echo "TIP: Common Name is your host name"
        openssl genrsa -aes256 -out ca-key.pem 4096
        openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem

        echo ""
        echo "STEP 2: create a server key and certificate signing request (CSR)"
        openssl genrsa -out server-key.pem 4096
        echo "Please enter the Common Name you have used above:"
        read commonname
        if [[ $commonname == "" ]]
            then
            echo "Common Name cannot be empty. Exiting"
            exit 1
        fi
        openssl req -subj "/CN=$commonname" -sha256 -new -key server-key.pem -out server.csr

        echo ""
        echo "STEP 3: signing the public key with our CA"
        echo "Since TLS connections can be made via IP address as well as DNS name, they need to be specified when creating the certificate."
        echo "Please enter the IP address of your docker server machine. TIP: most likely it is one of those addresses:"
        echo "$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')"
        read ipaddress
        if [[ $ipaddress == "" ]]
            then
            echo "IP address cannot be empty. Exiting"
            exit 1
        fi

        echo subjectAltName = DNS:$commonname,IP:$ipaddress,IP:127.0.0.1 > extfile.cnf
        openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile extfile.cnf

        echo ""
        echo "STEP 4: client key and certificate signing request"
        openssl genrsa -out key.pem 4096
        openssl req -subj '/CN=client' -new -key key.pem -out client.csr


        echo extendedKeyUsage = clientAuth > extfile.cnf
        openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem -extfile extfile.cnf

        echo ""
        echo "STEP 5: postprocessing"
        rm -v client.csr server.csr
        chmod -v 0400 ca-key.pem key.pem server-key.pem
        chmod -v 0444 ca.pem server-cert.pem cert.pem
        echo ""
        echo "Congratulations! Now you have certificates. To make docker daemon more secure you can force it accept connections only form the clients with certificates:"
        echo "dockerd --tlsverify --tlscacert=ca.pem --tlscert=server-cert.pem --tlskey=server-key.pem"
        echo ""
        echo "Now you can access docker from another machine, e.g. by using the command similar to this one (change the ip and port accordingly):"
        echo "docker --tlsverify --tlscacert=ca.pem --tlscert=cert.pem --tlskey=key.pem -H=127.0.0.1:59047"
        ;;
    *)
        echo "Operation has been cancelled"
        ;;
esac


