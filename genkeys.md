### This may only work for localhost

openssl req -new -subj "/C=GB/ST=Devon/CN=localhost" -newkey rsa:2048 -nodes -keyout server.key -out server.csr   
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt