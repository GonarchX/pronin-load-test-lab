#!/bin/bash

readonly LOAD_REPORT_FILE=load_report.txt

./wrk/wrk -d 10 -t 6 -c 12 -R 500 -s ./wrk/post.lua -L http://localhost:9080 > $LOAD_REPORT_FILE

if grep -q "Non-2xx or 3xx responses" $LOAD_REPORT_FILE; then
        echo "При проведении нагрузочного тестирования присутствовали запросы, которые вернули ошибки"
elif grep -q "Socket errors:" $LOAD_REPORT_FILE; then
        echo "При проведении нагрузочного тестирования присутствовали запросы, которые вернули ошибки"
else
        echo "Все запросы к серверу З/П работников успешно завершились"
fi

# Смотрим, что в логах нету ошибок
readonly LOG_DIR="../artifacts/logs/err"

if [ -s "$LOG_DIR/server" ]; then
    echo "Сервер с З/П работников содержит ошибки:"
    cat "$LOG_DIR/server"
fi

if [ -s "$LOG_DIR/kafka_consumer" ]; then
    echo "Сервис kafka_consumer содержит ошибки:"
    cat "$LOG_DIR/kafka_consumer"
fi

# Смотрим, что все контейнеры живы
containers_to_check=(
    "load-server-postgres-local"
    "load-test-server-go"
    "load-test-kafka-ui"
    "load-test-kafka"
    "mongo_container"
    "load-test-mongo-express"
    "load-test-spark-writer"
    "load-test-kafka-consumer"
)

running_containers=$(docker ps --format "{{.Names}}")

stopped_containers=()

for container in "${containers_to_check[@]}"; do
    if ! echo "$running_containers" | grep -q "$container"; then
        stopped_containers+=("$container")
    fi
done

if [ ${#stopped_containers[@]} -ne 0 ]; then
    echo "Список упавших контейнеров:"
    for container in "${stopped_containers[@]}"; do
        echo "* $container"
    done
else
    echo "Все контейнеры остались в рабочем состоянии"
fi
