# test_assignment

---
## Run local compose command:

make compose-test-up

---

## Request examples:

curl 'localhost:8081/sites&search=ya.ru'

curl 'localhost:8081/sites/stats?from=2022-03-01T00:04:05Z&to=2022-03-02T23:04:05Z'

---
## Config example:


{
    "http_server": {
        "port": "8081",
        "timeout": 30,   
        "ttl": 30
    },
    "data_base":{
        "host":"postgres",
        "port":"5432",
        "user":"myuser",
        "password":"secret",
        "db_name":"mydb"
    },
    "zap_logger_mode": "development"
}

---