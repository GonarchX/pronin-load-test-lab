from pyspark.sql import SparkSession
from pyspark.sql.types import StructField, IntegerType, StringType, StructType

def main():
    spark = SparkSession.builder \
        .appName("LoadTestSparkWriter") \
        .getOrCreate()

    schema = StructType([
        StructField("id", IntegerType(), False),
        StructField("name", StringType(), False),
        StructField("salary", IntegerType(), False)
    ])

    # Читаем поток данных из директории с массивами JSON файлами
    df = spark.readStream \
        .schema(schema) \
        .option("multiLine", True) \
        .json("/employees")  # Замените на ваш путь

    # Преобразуем датафрейм с добавлением ключа
    json_df = df.selectExpr("to_json(struct(*)) AS value", "CAST(id AS STRING) AS key")

    # Запись данных в Kafka
    query = json_df.writeStream \
        .format("kafka") \
        .option("kafka.bootstrap.servers", "kafka:29092") \
        .option("topic", "load-test-topic") \
        .option("checkpointLocation", "/tmp/checkpoints") \
        .start()

    query.awaitTermination()

if __name__ == "__main__":
    main()
