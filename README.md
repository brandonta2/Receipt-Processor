# Receipt-Processor
Receipt Processor

This project uses Go to host a local web server that can accept json POST data at 
http://localhost:8080/receipts/process 
and return an ID for 
http://localhost:8080/receipts/<ID>/points
which returns the price based on internal logic.

For testing I used Postman to send a raw json body to the process endpoint.

For running with docker, download the project, then using terminal from the project root enter:
docker run -p 8080:8080 go-web-service

which will run the service on http://localhost:8080.

Using Postman send a POST to http://localhost:8080/receipts/process with the test json data in raw body. Set Headers to Key = Content-Type Value = application/json

Put the returned ID in http://localhost:8080/receipts/<ID>/points and get returned the price.


Troubleshooting:

If errors then might need to run: go get github.com/google/uuid 
rebuild with: docker build -t go-web-service
