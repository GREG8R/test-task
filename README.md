# test-task
this is repository for test task - api for getting transactions and show history by hours.

How to launch server:
1) clone repository to your local computer
2) install docker https://docs.docker.com/get-docker/
3) run this command in the terminal from the project directory /test-task/: docker-compose up --build
   after this command server was launched on the 8080 port on the localhost.
4) now you can use this request for send some money:
    `curl --location --request POST 'localhost:8080/send' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "datetime": "2019-10-05T15:48:01+01:00",
    "amount": 1.2
    }'`
   
    this request for getting history:
    `curl --location --request GET 'localhost:8080/history' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "startDatetime": "2019-10-05T12:48:01+01:00",
    "endDatetime": "2019-10-05T16:48:01+01:00"
    }'`
    
    and this one for getting metrics:
    `curl --location --request GET 'localhost:8080/metrics'`

For using these curls you can install Postman and import requests there, 
or you can use curl console util in the terminal
