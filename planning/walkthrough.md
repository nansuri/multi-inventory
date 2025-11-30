# Multi Inventory Management System - Walkthrough

## Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

## How to Run

### 1. Database
Start the PostgreSQL database using Docker Compose:
```bash
docker-compose up -d postgres
```

### 2. Backend
Navigate to the backend directory and start the server:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.
It will automatically run migrations on startup.

### 3. Frontend
Navigate to the frontend directory and start the dev server:
```bash
cd frontend
npm install
npm run dev
```
The frontend will be available at `http://localhost:5173`.

## Features Walkthrough

### Authentication
1.  Go to `http://localhost:5173`.
2.  Click "Create Account" to register a new user.
3.  Login with your credentials.

### Inventory Management
1.  From the Dashboard, click "Inventory" to see the list.
2.  Click "Add" (top right) to create a new item.
3.  Use the "Scan" icon in the Barcode field to simulate scanning (requires camera permission).
4.  Click on an item to Edit or Delete it.

### Sales Management
1.  From the Dashboard, click "New Sale".
2.  Scan items or manually add them to the cart (currently manual add via code is not UI-implemented, but scanning works if item exists).
    - *Tip*: Ensure you have items in inventory first.
3.  Click "Submit Order".
4.  View "Sales History" from the Dashboard.
5.  Click on an order to enter "Order Checker" mode.
6.  Tap items to mark them as fulfilled/checked.
