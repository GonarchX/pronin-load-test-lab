./wrk -d 10s -t 1 -c 1 -R 20000 -s post.lua  http://localhost:9080/v0/entry

.build server
docker build -t load-test-server-go -f ./cmd/server/Dockerfile .

.build kafka-consumer
docker build -t load-test-kafka-consumer -f ./cmd/kafka_consumer/Dockerfile .

.build spark-writer
docker build -t load-test-spark-writer  .

.run containers