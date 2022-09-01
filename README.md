## musics-features.go
Pega as musics features a partir de um arquivo txt com o link de várias músicas.

## playlist-features.go
Pega as musics features de todas as músicas salvas pelo "get-musics-ids.go".

## get-musics-ids.go
Pega os IDs de todas as músicas de alguma playlist, basta inserir o link no arquivo "playlists.txt" e salva num txt.

## get-genres.go
Pega os generos dos artistas já buscados anteriormente e busca pelos gêneros deles.

## get-artists.go
Pega os IDs dos artistas a partir dos IDs das músicas e salva tanto num txt.

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