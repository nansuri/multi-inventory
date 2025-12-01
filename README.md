# Multi Inventory Management System

A mobile-web application for managing inventory in a multi-user environment with barcode scanning, sales order management, and real-time inventory tracking.

## 🚀 Features

### User Authentication & Authorization
- Secure user registration and login
- JWT-based authentication
- Role-based access control

### Inventory Management
- **Barcode/QR Code Scanning**: Integrated barcode scanner for quick item lookup and entry
- **Item Management**: Add, update, and delete inventory items
- **Item Details**: Track price, location, halal/non-halal status, and more
- **Real-time Updates**: Keep inventory synchronized across all users

### Sales Management
- **Sales Order Creation**: 
  - Scan barcodes to add items to orders
  - Automatic inventory verification
  - Real-time quantity and total price calculation
- **Order Fulfillment**: Track order completion status
- **Order Checker**: Verify individual items in orders are fulfilled
- **Order History**: View, update, and delete sales orders

### User Interface
- Responsive mobile-first design
- Intuitive sidebar navigation
- Dashboard with key metrics
- User profile management

## 🛠️ Tech Stack

### Frontend
- **Framework**: Vue.js 3.5
- **UI Library**: Vant 4.9 (Mobile UI components)
- **State Management**: Pinia
- **Routing**: Vue Router
- **Barcode Scanner**: html5-qrcode
- **Build Tool**: Vite 7.2
- **Styling**: SASS

### Backend
- **Language**: Go 1.23
- **Architecture**: Domain-Driven Design (DDD)
- **HTTP Router**: Chi v5
- **Database ORM**: GORM
- **Database Driver**: pgx/v5
- **Authentication**: JWT with bcrypt
- **CORS Handling**: go-chi/cors

### Database
- **RDBMS**: PostgreSQL

### DevOps
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Deployment**: Multi-container setup

## 📋 Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- Node.js 18+ and npm (for local frontend development)
- PostgreSQL (if running outside Docker)

## 🚦 Getting Started

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd multi-inventory
   ```

2. **Configure environment variables**
   
   The application connects to an existing PostgreSQL database. Update `docker-compose.yml` with your database credentials if needed:
   ```yaml
   environment:
     - DB_HOST=qa_platform_db
     - DB_USER=postgres
     - DB_PASSWORD=your_password
     - DB_NAME=qa_platform
     - DB_PORT=5432
   ```

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

### Local Development

#### Backend

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the backend directory:
   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=qa_platform
   DB_PORT=5432
   JWT_SECRET=your_jwt_secret
   ```

4. **Run migrations**
   ```bash
   # Apply database migrations
   psql -h localhost -U postgres -d qa_platform -f migrations/001_initial_schema.sql
   ```

5. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

   The backend will start on http://localhost:8080

#### Frontend

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

   The frontend will start on http://localhost:5173

4. **Build for production**
   ```bash
   npm run build
   ```

## 📁 Project Structure

```
multi-inventory/
├── backend/
│   ├── cmd/
│   │   └── server/          # Application entry point
│   ├── internal/
│   │   ├── application/     # Application services & business logic
│   │   ├── domain/          # Domain models & entities
│   │   └── infrastructure/  # External dependencies (HTTP, DB)
│   │       ├── http/        # HTTP handlers
│   │       └── postgres/    # Database repositories
│   ├── migrations/          # Database migrations
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/      # Reusable Vue components
│   │   ├── views/           # Page components
│   │   ├── router/          # Vue Router configuration
│   │   ├── assets/          # Static assets
│   │   ├── App.vue
│   │   └── main.js
│   ├── public/              # Public static files
│   ├── Dockerfile
│   ├── package.json
│   └── vite.config.js
├── planning/                # Project planning documents
├── docker-compose.yml
└── README.md
```

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout

### Inventory
- `GET /api/inventory` - List all inventory items
- `GET /api/inventory/:id` - Get item details
- `POST /api/inventory` - Create new item
- `PUT /api/inventory/:id` - Update item
- `DELETE /api/inventory/:id` - Delete item
- `GET /api/inventory/barcode/:code` - Search by barcode

### Sales
- `GET /api/sales` - List all sales orders
- `GET /api/sales/:id` - Get order details
- `POST /api/sales` - Create new order
- `PUT /api/sales/:id` - Update order
- `DELETE /api/sales/:id` - Delete order
- `POST /api/sales/:id/fulfill` - Mark order as fulfilled

## 🏗️ Architecture

The backend follows **Domain-Driven Design (DDD)** principles:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases and service orchestration
- **Infrastructure Layer**: External interfaces (HTTP handlers, database repositories)

This separation ensures:
- Clean architecture with clear boundaries
- Testability and maintainability
- Business logic independence from frameworks
- Easy to extend and modify

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm run test
```

## 🔒 Security

- Passwords hashed using bcrypt
- JWT-based authentication
- CORS configured for cross-origin requests
- SQL injection prevention via parameterized queries (GORM)
- Input validation on all endpoints

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📧 Contact

For questions or support, please open an issue in the repository.

## 🗺️ Roadmap

- [ ] Mobile app (Flutter/React Native)
- [ ] Advanced reporting and analytics
- [ ] Multi-location support
- [ ] Batch inventory import/export
- [ ] Email notifications
- [ ] Advanced user permissions
- [ ] Inventory alerts and low-stock warnings
- [ ] Integration with external POS systems

## 🙏 Acknowledgments

- Vue.js team for the amazing framework
- Go community for excellent tooling and libraries
- All contributors and users of this project
