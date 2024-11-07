FROM bitnami/spark:latest

COPY load_test_spark_writer.py /opt/bitnami/spark/app/

ENV SPARK_HOME=/opt/bitnami/spark

# Добавляем пакет Kafka для Spark
CMD ["/opt/bitnami/spark/bin/spark-submit", \
     "--packages", "org.apache.spark:spark-sql-kafka-0-10_2.12:3.5.3", \
     "--master", "local[*]", \
     "/opt/bitnami/spark/app/load_test_spark_writer.py"]
