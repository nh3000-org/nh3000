echo "root"
cfssl  genkey -initca ca-root-csr.json | cfssljson -bare ca-root
echo "sign"
cfssl gencert -ca ca-root.pem -ca-key ca-root-key.pem -config=ca-config.json -profile=server server-csr.json | cfssljson -bare server
cfssl gencert -ca ca-root.pem -ca-key ca-root-key.pem -config=ca-config.json -profile=client client-csr.json | cfssljson -bare client
openssl x509 -noout -text -in ca-root.pem
openssl x509 -noout -text -in server.pem
openssl x509 -noout -text -in client.pem
