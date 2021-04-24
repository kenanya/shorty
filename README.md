## A. Description
This project contains of 1 microservice that aim to shorten urls. This microservice provides 3 endpoint:
- POST /shorten
- GET /:shortcode
- GET /:shortcode/stats

The data will be stored to MongoDB. Tech Stack: Golang, echo, MongoDB, docker

## B. How to Start
Please extract the zip directory and enter to it. 
```bash
unzip test_amartha.zip
cd test_amartha/shorty
```

## C. Configuration
The config file is located at:
- test_amartha/shorty/common/configGlobal.yaml

You can change the values according to your configuration. 

## D. Build and Run Project
#### Build Project
```bash
docker-compose -f docker-compose-local.yml build
```

#### Run Project
```bash
docker-compose -f docker-compose-local.yml up
```

## E. Consume Service
Below are the sample requests and expected responses for each microservice:

### E.1. POST /shorten
| Name | Value |
| ------ | ------ |
| endpoint | localhost:9701/v1/shorten |

#### Positive Scenario 1 - Request
```json
curl --location --request POST 'localhost:9701/v1/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url" : "https://en.wikipedia.org/",
    "shortcode": "Mcdp4W"
}'
```

#### Positive Scenario 1 - Response
```json
201 Created
Content-Type: "application/json"

{
    "shortcode": "Mcdp4W"
}
```

#### Negative Scenario 1 - Request
```json
curl --location --request POST 'localhost:9701/v1/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url" : "https://en.wikipedia.org/",
    "shortcode": "Mcdp4W1"
}'
```

#### Negative Scenario 1 - Response
```json
422 Unprocessable Entity
Content-Type: "application/json"

{
    "description": "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{6}$."
}
```

### E.2. GET /:shortcode
| Name | Value |
| ------ | ------ |
| endpoint | localhost:9701/v1/:shortcode |

#### Positive Scenario 1 - Request
```json
curl --location --request GET 'localhost:9701/v1/Mcdp4W' \
--header 'Content-Type: application/json'
```

#### Positive Scenario 1 - Response
```json
In this case it will be redirected to the https://en.wikipedia.org/ 
```

#### Negative Scenario 1 - Request
```json
curl --location --request GET 'localhost:9701/v1/111111' \
--header 'Content-Type: application/json'
```

#### Negative Scenario 1 - Response
```json
404
Content-Type: "application/json"

{
    "description": "The shortcode cannot be found in the system"
}
```

### E.3. GET /:shortcode/stats
| Name | Value |
| ------ | ------ |
| endpoint | localhost:9701/v1/:shortcode/stats |

#### Positive Scenario 1 - Request
```json
curl --location --request GET 'localhost:9701/v1/Mcdp4W/stats' \
--header 'Content-Type: application/json'
```

#### Positive Scenario 1 - Response
```json
200 OK
Content-Type: "application/json"

{
    "startDate": "2021-04-25T06:30:43+0700",
    "lastSeenDate": "2021-04-25T06:33:09+0700",
    "redirectCount": 2
}
```

#### Negative Scenario 1 - Request
```json
curl --location --request GET 'localhost:9701/v1/14/stats' \
--header 'Content-Type: application/json'
```

#### Negative Scenario 1 - Response
```json
404
Content-Type: "application/json"

{
    "description": "The shortcode cannot be found in the system"   
}
```



## F. Unit Test
These are the steps to run unit test for each microservice:
```bash
APP_ENV=local go test ./pkg/service/v1/ -tags=unit_video_bank_service -v
APP_ENV=local go test ./pkg/service/v1/ -tags=unit_user_service -v
APP_ENV=local go test ./pkg/service/v1/ -tags=unit_analytic_telemetry_service -v
```

For the unit test purpose, I use mongo DB at host localhost and port 32768. 32768 is the external port that used by mongo DB from docker compose. You can check the port which is used by mongo container using `docker ps` command. If the port is not 32768, you can change `db_host_test` variable at config file.
```bash
❯ docker ps
CONTAINER ID    IMAGE   COMMAND                 PORTS       
8880d4b440d3    mongo   "docker-entrypoint.s…"  0.0.0.0:32768->27017/tcp
```

















