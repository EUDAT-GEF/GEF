#!/bin/sh

mkdir -p /etc/grid-security/vomsdir/fedcloud.egi.eu

cat > /etc/grid-security/vomsdir/fedcloud.egi.eu/voms1.grid.cesnet.cz.lsc << EOF 
/DC=org/DC=terena/DC=tcs/C=CZ/ST=Hlavni mesto Praha/L=Praha 6/O=CESNET/CN=voms1.grid.cesnet.cz
/C=NL/ST=Noord-Holland/L=Amsterdam/O=TERENA/CN=TERENA eScience SSL CA 3
EOF

cat > /etc/grid-security/vomsdir/fedcloud.egi.eu/voms2.grid.cesnet.cz.lsc << EOF 
/DC=org/DC=terena/DC=tcs/C=CZ/ST=Hlavni mesto Praha/L=Praha 6/O=CESNET/CN=voms2.grid.cesnet.cz
/C=NL/ST=Noord-Holland/L=Amsterdam/O=TERENA/CN=TERENA eScience SSL CA 3
EOF

cat >> /etc/vomses << EOF 
"fedcloud.egi.eu" "voms1.grid.cesnet.cz" "15002" "/DC=org/DC=terena/DC=tcs/C=CZ/ST=Hlavni mesto Praha/L=Praha 6/O=CESNET/CN=voms1.grid.cesnet.cz" "fedcloud.egi.eu" "24"
"fedcloud.egi.eu" "voms2.grid.cesnet.cz" "15002" "/DC=org/DC=terena/DC=tcs/C=CZ/ST=Hlavni mesto Praha/L=Praha 6/O=CESNET/CN=voms2.grid.cesnet.cz" "fedcloud.egi.eu" "24"
EOF

# CA Certificates
echo "Downloading CA certificates"
mkdir -p /etc/grid-security/certificates
CA_BUNDLE=https://dist.eugridpma.info/distribution/igtf/current/accredited/igtf-preinstalled-bundle-classic.tar.gz
curl -s $CA_BUNDLE | tar -xz -C /etc/grid-security/certificates