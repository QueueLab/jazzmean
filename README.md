# Middleware for AI SDK Integration

This project implements middleware to facilitate communication between environment-aware agents and an AI SDK using Next.js. The middleware ensures complete contextual synchronicity between the generated UI and generated AI for each query. It also implements multi-modal magnetic highly concurrent retrieval-augmented generation (RAG) with PostgreSQL using an ORM tool.

## Features

- Communication between environment-aware agents and AI SDK using Next.js
- Complete contextual synchronicity between generated UI and generated AI for each query
- Multi-modal magnetic highly concurrent retrieval-augmented generation (RAG) with PostgreSQL using an ORM tool
- User authentication and authorization with role-based access control
- Real-time query monitoring and analytics with a dashboard for visualization
- Subtraction feature for handling subtraction queries

## Prerequisites

- Go 1.16 or later
- Node.js 14 or later
- PostgreSQL 12 or later

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/githubnext/workspace-blank.git
   cd workspace-blank
   ```

2. Set up environment variables:

   ```sh
   export AI_SDK_API_KEY=your_ai_sdk_api_key
   export POSTGRES_URL=postgres://user:password@localhost:5432/dbname
   export GO_AGENT_URL=http://localhost:8080
   ```

3. Install dependencies:

   ```sh
   go mod tidy
   npm install
   ```

4. Compile the Go WebAssembly module:

   ```sh
   GOOS=js GOARCH=wasm go build -o public/agent.wasm agent.go
   ```

5. Run the middleware server:

   ```sh
   go run main.go
   ```

6. Run the Next.js application:

   ```sh
   npm run dev
   ```

## Usage

1. Send a POST request to the middleware server with a query:

   ```sh
   curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query": "your_query"}'
   ```

2. The middleware will process the query with the AI SDK and PostgreSQL, and return the response.

3. Send a POST request to the middleware server for subtraction:

   ```sh
   curl -X POST http://localhost:8080/subtract -H "Content-Type: application/json" -d '{"a": 10, "b": 5}'
   ```

4. The middleware will process the subtraction query and return the result.

## Overview

The middleware is designed to handle queries from the generated UI and process them using the AI SDK and PostgreSQL. It ensures complete contextual synchronicity between the generated UI and generated AI for each query. Highly concurrent retrieval-augmented generation (RAG) with PostgreSQL using an ORM tool provides efficient and scalable data processing.

## License

This project is licensed under the MIT License.
