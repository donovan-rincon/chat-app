# Chat Application

A simple browser-based chat application using Go, PostgreSQL, and Docker. 

Registered users can log in and chat with other users in a chatroom.
The application also includes a bot that responds to the `/stock=stock_code` command.

## Features

- User registration and login
- Real-time chatroom
- Bot that responds to `/stock=stock_code` command
- Messages associated with users
- Persistent storage with PostgreSQL
- Dockerized setup with Docker Compose

## Prerequisites
- [Go](https://go.dev/doc/install)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Configuration

1. Create the .env file for DB secrets and config and place it in root directory of the project:
    ```env
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=chatapp
    DB_PORT=5432
    DB_SSLMODE=disable
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
    ```

## Running the Application

1. Build and run the Docker containers:

    ```bash
    docker-compose up --build
    ```

2. Access the application:

    - Open your browser and navigate to [`http://localhost:8080/public/login.html`] to log in or register.
    - After logging in, you will be redirected to the chatroom.

## Project Structure

![alt text](image.png)