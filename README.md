# hr-system-backend

## Overview
This project provides a Go-based backend service with MySQL and Redis, orchestrated using Docker Compose. 
The Docker configurations are stored in the folder deployments.

## Getting Started

1. Build and Run Containers
In the deployments folder, run:

    ```bash
    docker-compose up --build
    ```

2. Environment Variables

    The app uses the following environment variables (defined in docker-compose.yml):
    
    MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_PASSWORD, MYSQL_DB_NAME
    REDIS_HOST, REDIS_PORT

3. Stopping the Containers

   To stop and remove containers, run:
   ```bash
   docker-compose down
   ```
   
4. Access the Application

    The application is accessible at http://localhost:8080.
