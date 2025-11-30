# Implementation Plan - Multi Inventory Management System MVP

## Goal Description
Create a mobile-first web application for inventory management. The system will support multiple users, inventory tracking (with barcode scanning), and sales management.
Targeting MVP status with a clean, modern mobile UI.

## User Review Required
> [!IMPORTANT]
> **Tech Stack Selection**:
> - **Frontend**: Vue 3 + Vite + Vant UI (Mobile-first components).
> - **Backend**: Go (Golang) using Chi router, following DDD principles.
> - **Database**: PostgreSQL (running in Docker).
> - **Containerization**: Docker Compose for local development (App + DB).

## Proposed Changes

### Project Structure
#### [NEW] [docker-compose.yml](file:///home/nansuri/Git/multi-inventory/docker-compose.yml)
- Defines `postgres` service.
- Defines `backend` and `frontend` services (optional for dev, but good for full stack run).

### Backend (Go)
Directory: `backend/`
- **Architecture**: DDD (Domain-Driven Design).
    - `cmd/server/`: Entry point.
    - `internal/domain/`: Entities (User, Item, SalesOrder) and Repository Interfaces.
    - `internal/application/`: Use Cases / Services (AuthService, InventoryService, SalesService).
    - `internal/infrastructure/`: Implementation of repositories (Postgres), HTTP Handlers (Chi), Auth Middleware (JWT).

#### [NEW] Backend Files
- `go.mod`, `go.sum`
- `cmd/server/main.go`
- `internal/domain/user.go`, `item.go`, `order.go`
- `internal/infrastructure/postgres/`: DB connection and repository implementations.
- `internal/infrastructure/http/`: REST API handlers.

### Frontend (Vue 3)
Directory: `frontend/`
- **Tech**: Vue 3, Vite, Vant (UI Kit), Vue Router, Pinia (State Management).
- **Features**:
    - **Layout**: Mobile sidebar/bottom nav.
    - **Views**:
        - `Login.vue`
        - `Dashboard.vue`
        - `InventoryList.vue` / `InventoryEdit.vue`
        - `SalesList.vue` / `SalesCreate.vue`
    - **Components**:
        - `BarcodeScanner.vue`: Wrapper around a JS barcode scanning library (e.g., `html5-qrcode`).

#### [NEW] Frontend Files
- `package.json`, `vite.config.js`
- `src/App.vue`, `src/main.js`
- `src/views/...`
- `src/components/...`

## Verification Plan

### Automated Tests
- **Backend**: Unit tests for Domain/Application logic.
    - Run: `go test ./...` in `backend/` directory.
- **Frontend**: Basic component mount tests (if time permits, otherwise manual).

### Manual Verification
1.  **Setup**: Run `docker-compose up -d postgres`. Run backend `go run cmd/server/main.go`. Run frontend `npm run dev`.
2.  **Auth**: Register a new user, Login. Verify JWT token.
3.  **Inventory**:
    - Add item manually.
    - Scan barcode (simulate with camera or file input if on desktop) to fill field.
    - Edit/Delete item.
4.  **Sales**:
    - Create order.
    - Scan item to add to order.
    - Verify totals.
    - "Order Checker" mode: check off items.
