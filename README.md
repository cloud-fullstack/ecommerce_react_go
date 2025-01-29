# Full-Stack E-Commerce Platform Excerpt

## Overview

Excerpt of full-stack e-commerce platform [demo featuring a Go backend with Gin, PostgreSQL, and AWS S3 for storage, along with a React-based frontend. This platform allows users to browse products, view discounts, interact with blogs, and manage their shopping experience seamlessly.

## Features

- **Backend (Go + Gin):**

  - REST API with authentication and authorization.
  - PostgreSQL database integration.
  - AWS S3 integration for image storage.
  - CORS support for secure cross-origin requests.
  - Logging system with PostgreSQL storage.

- **Frontend (React + TypeScript):**

  - User-friendly UI with React Router.
  - Product previews and detailed pages.
  - Shopping cart management.
  - Cookie consent modal.
  - Blog interaction and product filtering.

## Tech Stack

### Backend

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL
- **Storage:** AWS S3 (DigitalOcean Spaces)
- **Authentication:** Middleware-based authentication
- **Hosting:** Azure Container Registry, Render.com

### Frontend

- **Language:** TypeScript
- **Framework:** React
- **UI Library:** TailwindCSS 
- **Routing:** React Router
- **API Communication:** Axios
- **State Management:** useState, useEffect

## Setup & Installation

### Prerequisites

- Go 1.19+
- Node.js 18+
- PostgreSQL database
- AWS S3 (or DigitalOcean Spaces) 
- Docker (for containerized deployment)

### Backend Setup

1. Clone the repository:
   ```sh
   git clone https://github.com/cloud-fullstack/ecommerce_react_go.git
   cd ./backend
   ```
2. Set environment variables:
   ```sh
   export API_PORT=8080
   export DB_HOST=your_db_host
   export DB_PORT=5432
   export DB_USER=your_db_user
   export DB_PASSWORD=your_db_password
   export DB_NAME=your_db_name
   export DB_SSLMODE=require
   export SPACES_ACCESS_KEY=your_aws_key
   export SPACES_SECRET_KEY=your_aws_secret
   export BACKEND_APP_DOMAIN_NAME=https://yourbackend.com
   ```
3. Run the backend:
   ```sh
   go run main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```sh
   cd ./frontend
   ```
2. Install dependencies:
   ```sh
   npm install
   ```
3. Run the frontend:
   ```sh
   npm start
   ```

## API Endpoints

### Public Routes

- `GET /api/frontpage-product-previews/` - Get product previews for homepage
- `GET /api/store-details/:storeID` - Fetch details of a store
- `GET /api/most-loved-recent-blogs/` - Get recent popular blog posts
- `GET /api/discounted-products-frontpage/` - Get discounted products
- `POST /api/gen-token/` - Generate authentication token

### Authenticated Routes (Require JWT)

- `POST /api/upload-picture/:pictureType` - Upload images to S3
- `DELETE /api/answer/:answerID` - Delete a user’s answer
- `POST /api/create-order/` - Create a new order
- `POST /api/cancel-order/` - Cancel an order

## Deployment

### Docker Deployment

1. Build and run the Docker container:
   ```sh
   docker-compose up --build
   ```
2. Push images to Azure Container Registry:
   ```sh
   docker tag backend youracr.azurecr.io/backend
   docker push youracr.azurecr.io/backend

### Render Deployment

1. Set up the service on Render.com.
2. Connect the repository and configure environment variables.
3. Deploy and monitor logs in Render.

