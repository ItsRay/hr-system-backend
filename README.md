# hr-system-backend

## Overview
This project is a practice backend service designed to simulate an HR management system.
It provides a Go-based backend service with MySQL and Redis, orchestrated using Docker Compose.

## Prerequisites
- Docker & Docker Compose

## Quick Start

Start the application with MySQL and Redis:
```bash
make compose-up
```
Shut down all services:
```bash 
make compose-down
```

See the Makefile for more commands.

## API Testing

Use hr-system.postman_collection.json to import the Postman collection for testing the API.

## API introduction

#### 1. Create Employee
- Method: POST
- Path: /api/employees
- Description: Creates a new employee and saves it to the database.

#### 2. Get Employee by ID
- Method: GET
- Path: /api/employees/{id} 
- Description: Retrieves the details of an employee by their ID.

#### 3. Get Employees (Paginated)
- Method: GET
- Path: /api/employees?page={page}&page_size={page_size}
- Description: Retrieves a paginated list of employees and their total count.

#### 4. Create Leave
- Method: POST
- Path: /api/leaves
- Description: Submits a new leave request for an employee.

#### 5. Get Leave by ID
- Method: GET
- Path: /api/leaves/{id}
- Description: Retrieves details of a specific leave request by its ID.

#### 6. Get Leaves
- Method: GET
- Path: /api/leaves?employee_id={employee_id}
- Description: Retrieves leaves of an employee.

#### 7. Review a Leave
- Method: POST
- Path: /api/v1/leaves/{id}/review
- Description: Review a leave request.
