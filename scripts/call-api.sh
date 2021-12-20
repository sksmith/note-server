curl -i -H "content-type:application/json" \
    -XPUT -d'{"id":"1","note":"somenote"}' \
    http://localhost:8080/api/v1/note/1

curl -i -H "content-type:application/json" \
    http://localhost:8080/api/v1/note/1