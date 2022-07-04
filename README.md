## get-token.js Deployment

To deploy get-token endpoint. Create a docker.env file with your credentials (CLIENT_ID and CLIENT_SECRET). 
Then run this commands below with docker.

```bash
  docker-compose build
```
```bash
  docker-compose up
```


## get-token.js Endpoint

#### Get Token - This endpoint will cat your credentials at docker.env file inside container and return json with your Bearer Token.

```http
  GET localhost:8080/token
```