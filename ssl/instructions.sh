# changes the cn's to match your hosts in your enviroment if needed
SERVER_CN=localhost

# generate certificate authority + trust certificate ca.cry
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# genre the server private key 
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# get a certificate signing request from the ca
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# sign the certificate with the ca we created ('its called self signing') - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -extfile /etc/ssl/openssl.cnf -extensions Server -out server.crt

# convert the server cetificate to .pem format -usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem