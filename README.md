# Go Notification Service

A notification service built using Go, Kafka.

## Getting Started

- Just look at the `Makefile` and run `make` to see the available commands.

- Ensure that an `.env` file is created in the root directory of the project. You can use the `.env.sample` file as a template.

- Ensure that you either have a running Kafka cluster or you can use the `docker-compose.yml` file to spin up a local Kafka cluster. Use the `Makefile` to spin one up quickly.

## Testing

Send a notification to the service by running the following command:

```bash
curl --location 'http://localhost:6000/producer/send' -d '{"from_id": 2,"to_id": 1,"message": "2 followed 1"}'
```

Read the notifications of a user by running the following command:

```bash
curl --location 'http://localhost:6000/consumer/notifications/1'
```
