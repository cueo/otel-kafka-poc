# Kafka

## Setting up Kafka

### Download and untar Kafka
```shell
wget https://dlcdn.apache.org/kafka/3.5.0/kafka_2.13-3.5.0.tgz -O kafka.tgz
tar -xzf kafka.tgz
cd kafka_2.13-3.5.0
```

### Start Kafka server
```shell
KAFKA_CLUSTER_ID="$(bin/kafka-storage.sh random-uuid)"
bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c config/kraft/server.properties
bin/kafka-server-start.sh config/kraft/server.properties
```

### Create a topic
```shell
export KAFKA_METRICS_TOPIC="metrics"
bin/kafka-topics.sh --create --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
```
### Read metrics from the topic
```shell
bin/kafka-console-consumer.sh --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
```


## Handy aliases
```shell
KAFKA_DIR="~/dev/kafka/kafka_2.13-3.5.0"
kafka_start() {
  KAFKA_CLUSTER_ID="$($KAFKA_DIR/bin/kafka-storage.sh random-uuid)"
  $KAFKA_DIR/bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c $KAFKA_DIR/config/kraft/server.properties
  $KAFKA_DIR/bin/kafka-server-start.sh $KAFKA_DIR/config/kraft/server.properties
}

kafka_create_topic() {
  if [ -z ${1+x} ]
  then
    $KAFKA_DIR/bin/kafka-topics.sh --create --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
  else
    $KAFKA_DIR/bin/kafka-topics.sh --create --topic $1 --bootstrap-server localhost:9092
  fi
}

kafka_delete_topic() {
  if [ -z ${1+x} ]
  then
    $KAFKA_DIR/bin/kafka-topics.sh --delete --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
  else
    $KAFKA_DIR/bin/kafka-topics.sh --delete --topic $1 --bootstrap-server localhost:9092
  fi
}

kafka_describe_topic() {
  if [ -z ${1+x} ]
  then
    $KAFKA_DIR/bin/kafka-topics.sh --describe --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
  else
    $KAFKA_DIR/bin/kafka-topics.sh --describe --topic $1 --bootstrap-server localhost:9092
  fi
}

kafka_consume() {
  if [ -z ${1+x} ]
  then
    $KAFKA_DIR/bin/kafka-console-consumer.sh --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
  else
    $KAFKA_DIR/bin/kafka-console-consumer.sh --topic $1 --bootstrap-server localhost:9092
  fi
}

kafka_produce() {
  if [ -z ${1+x} ]
  then
    $KAFKA_DIR/bin/kafka-console-producer.sh --topic $KAFKA_TOPIC --bootstrap-server localhost:9092
  else
    $KAFKA_DIR/bin/kafka-console-producer.sh --topic $1 --bootstrap-server localhost:9092
  fi
}

kafka_kill() {
  lsof -t -i:9092 | xargs kill -9
}

kafka_nuke() {
  kafka_kill
  rm -rf /tmp/kafka-logs /tmp/zookeeper /tmp/kraft-combined-logs
}
```
