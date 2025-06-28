# Project Backend Go - E Wallet

This project was made by Muhammad Davinda Rinaldy in Kodacademy Training Program. This project uses Go Language to make a backend application for e-wallet and PostgreSQL for the database.

Endpoint included in this project:
1. Auth Flow (Register User & Login User)
2. Users (Update Profile, Get User By Email, Get All User & Search User by Phone Number and Name)
3. Transaction (Transfer, Top Up, History Transaction, Balance, Income & Expense)

This project also protect endpoint for users and transaction with token (generated when user login) and utilizing middleware to verify the token.

## Prerequisites

Make sure you already install Go to run this project

## How to Run this Project

1. Create a new empty directory for the project and navigate into it
2. Clone this project into the empty current directory:
```
git clone https://github.com/mdavindarinaldy/fgo24-be-weekly.git .
``` 
3. Install dependencies
```
go mod tidy
```
4. Run the project
```
go run main.go
```

## Entity Relationship Diagram (ERD)

```mermaid
erDiagram
    direction LR
    users {
        int id PK
        string name
        string email
        string phone_number
        string password
        string pin
    }
    user_balance {
        int id PK
        int id_user FK
        decimal balance
        timestamp created_at
    }
    transactions {
        int id PK
        timestamp transaction_date
        decimal nominal
        enum type
        text notes
        string id_user FK
        string id_other_user FK
    }
    users ||--o{ transactions : performs
    users ||--o{ transactions : receives
    users ||--o{ user_balance : has
```

## Dependencies
This project use:
1. gin-gonic from github.com/gin-gonic/gin : for handling HTTP request/response data (gin.Context), for defining middleware and route handlers (gin.HandlerFunc), for organizing routes into groups (gin.RouterGroup) and for managing HTTP routing and server configuration (gin.Engine)
2. jwt v5 from github.com/golang-jwt/jwt/v5 : for creating, parsing and validating JSON Web Tokens (JWT) for authentication and authorization
3. pgx from github.com/jackc/pgx/v5 : for direct database interactions (PostgreSQL)
4. godotenv from github.com/joho/godotenv : for loading environment variables from a .env file into the application

## Basic Information
This project is part of training in Kodacademy Bootcamp Batch 2 made by Muhammad Davinda Rinaldy