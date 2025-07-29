# Vehicle Showroom Management System - Backend

## Phase 2: Customer & Vehicle Management ✅

### Features Implemented:
- ✅ Clean Architecture (Entity, Repository, UseCase, Handler)
- ✅ PostgreSQL Database with SQLX
- ✅ JWT Authentication
- ✅ User Management (Admin, Mechanic, Cashier roles)
- ✅ Session Management
- ✅ Complete Database Schema (11 tables)
- ✅ CORS Support for Flutter integration
- ✅ **Customer Management CRUD**
- ✅ **Vehicle Management CRUD**
- ✅ **Demo Data Seeding**
- ✅ **Vehicle Status Management**
- ✅ **Customer-Vehicle Relationships**

### API Endpoints:

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get user profile

#### Customer Management
- `GET /api/v1/customers` - List customers (with pagination & search)
- `POST /api/v1/customers` - Create new customer
- `GET /api/v1/customers/:id` - Get customer by ID
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

#### Vehicle Management
- `GET /api/v1/vehicles` - List vehicles (with pagination, search & status filter)
- `POST /api/v1/vehicles` - Create new vehicle
- `GET /api/v1/vehicles/:id` - Get vehicle by ID
- `PUT /api/v1/vehicles/:id` - Update vehicle
- `PUT /api/v1/vehicles/:id/status` - Update vehicle status
- `DELETE /api/v1/vehicles/:id` - Delete vehicle

### Database Tables:
1. users ✅
2. user_sessions ✅
3. customers ✅
4. vehicles ✅
5. vehicle_images
6. spare_parts
7. purchase_transactions
8. sales_transactions
9. repairs
10. repair_parts
11. stock_movements

### Demo Data:
- ✅ 3 Demo Users (admin, cashier, mechanic)
- ✅ 5 Demo Customers (individual & corporate)
- ✅ 5 Demo Vehicles (various brands & statuses)

### Setup Instructions:
1. Install PostgreSQL
2. Create database: `vehicle_showroom`
3. Copy `.env` file and update database credentials
4. Run: `go mod tidy`
5. Run: `go run cmd/main.go`
6. Demo data will be automatically seeded

### Demo Users:
- Admin: admin / admin123
- Cashier: cashier / cashier123  
- Mechanic: mechanic / mechanic123

### Next Phase:
Phase 3: Transaction System (Purchase & Sales)
