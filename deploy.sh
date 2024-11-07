docker build -t load-test-server-go -f ./cmd/server/Dockerfile .

docker build -t load-test-kafka-consumer -f ./cmd/kafka_consumer/Dockerfile .

docker build -t load-test-spark-writer -f ./spark_writer/Dockerfile .

docker compose up -d