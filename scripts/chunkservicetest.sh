printf "save test: "
curl -X POST --data-binary @lena30.jpg http://localhost:8001/v1/chunk/lena30.jpg/0

printf "load test: "
curl http://localhost:8001/v1/chunk/lena30.jpg/0  -o result.jpg
