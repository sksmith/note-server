curl -i -u test:test -H "content-type:application/json" \
    -X PUT -d'{"id":"1","title": "First Note","data":"somenote"}' \
    http://localhost:8080/api/v1/note

curl -i -u test:test -H "content-type:application/json" \
    -X PUT -d'{"id":"2","title": "Second Note","data":"some other note with more info"}' \
    http://localhost:8080/api/v1/note

curl -i -u test:test -H "content-type:application/json" \
    -X DELETE http://localhost:8080/api/v1/note/1

curl -i -u test:test -H "content-type:application/json" \
    http://localhost:8080/api/v1/note/1

curl -i -u test:test -H "content-type:application/json" \
    -X PUT -d'{"id":"2","title": "Second Note - updated","data":"some other note with more info"}' \
    http://localhost:8080/api/v1/note

curl -i -u test:test -H "content-type:application/json" \
    -X PUT -d'{"id":"2","title": "Second Note - updated","data":"changed the contents"}' \
    http://localhost:8080/api/v1/note