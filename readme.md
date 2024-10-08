Grotto
=======
distributed file storage

# Demo

## Run

```sh
docker-compose up --build
```

## Upload file

```sh
curl -X POST --data-binary @lena30.jpg http://localhost:8000/v1/file/lena30.jpg
```

## Read catalog

```sh
curl http://localhost:8000/v1/files/all
```

## Download file

```sh
curl http://localhost:8000/v1/file/lena30.jpg -o result.jpg
```

## Delete file

```sh
curl -X DELETE http://localhost:8000/v1/file/lena30.jpg 
```