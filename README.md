# Receipt Processor  

This project is a Go-based web service that processes receipt data and returns an ID for retrieving points based on internal logic.  

## Features  
- Accepts JSON `POST` requests at `http://localhost:8080/receipts/process`  
- Returns a receipt ID  
- Allows querying points with `http://localhost:8080/receipts/{id}/points`  
- Docker support for easy deployment  

## Setup and Usage  

### Running with Docker  
Download the project and build the Docker image:  

```sh
docker build -t go-web-service .
```

Run the service with:  

```sh
docker run -p 8080:8080 go-web-service
```

Now the service will be available at `http://localhost:8080`.

## API Endpoints  

### Submit a Receipt  
- **Endpoint:** `POST /receipts/process`  
- **Request Body (JSON):**  

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    }
  ],
  "total": "6.49"
}
```

- **Response:**  

```json
{
  "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2"
}
```

### Retrieve Points  
- **Endpoint:** `GET /receipts/{id}/points`  
- **Response:**  

```json
{
  "points": 100
}
```

## Testing with Postman  

1. Open Postman  
2. Set method to **POST** and enter URL:  
   ```
   http://localhost:8080/receipts/process
   ```
3. Go to the **Body** tab, select **raw**, and paste the JSON request body.  
4. Set **Headers**:  
   - `Key`: `Content-Type`  
   - `Value`: `application/json`  
5. Send the request and copy the returned `id`.  
6. Use **GET** with the following URL:  
   ```
   http://localhost:8080/receipts/{id}/points
   ```
   (Replace `{id}` with the actual receipt ID)

## Troubleshooting  

If you encounter errors, try the following:  

1. Ensure dependencies are installed:  
   ```sh
   go get github.com/google/uuid
   ```  
2. Rebuild the Docker image:  
   ```sh
   docker build -t go-web-service .
   ```
