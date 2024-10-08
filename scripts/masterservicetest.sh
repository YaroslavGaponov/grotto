printf "save test: "
curl -X POST --data-binary @lena30.jpg http://localhost:8000/v1/file/lena30.jpg

printf "load metadata: "
curl  http://localhost:8001/v1/chunk/lena30.jpg.metadata/0

printf "catalog: "
curl http://localhost:8000/v1/files/all

printf "load test: "
curl http://localhost:8000/v1/file/lena30.jpg  -o result.jpg

printf "load test: "
curl -X DELETE http://localhost:8000/v1/file/lena30.jpg 

printf "catalog: "
curl http://localhost:8000/v1/files/all