# Vehicle Showroom Management System - Backend

## Phase 3: Transaction System ✅

### Features Implemented:
- ✅ Clean Architecture (Entity, Repository, UseCase, Handler)
- ✅ PostgreSQL Database with SQLX
- ✅ JWT Authentication
- ✅ User Management (Admin, Mechanic, Cashier roles)
- ✅ Session Management
- ✅ Complete Database Schema (11 tables)
- ✅ CORS Support for Flutter integration
- ✅ Customer Management CRUD
- ✅ Vehicle Management CRUD
- ✅ Demo Data Seeding
- ✅ Vehicle Status Management
- ✅ Customer-Vehicle Relationships
- ✅ **Purchase Transaction Management**
- ✅ **Sales Transaction Management**
- ✅ **Auto Transaction & Invoice Number Generation**
- ✅ **Dashboard Statistics**
- ✅ **Real-time Business Metrics**

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

#### Transaction Management
- `GET /api/v1/transactions/purchases` - List purchase transactions
- `POST /api/v1/transactions/purchases` - Create purchase transaction
- `GET /api/v1/transactions/purchases/:id` - Get purchase transaction by ID
- `GET /api/v1/transactions/sales` - List sales transactions
- `POST /api/v1/transactions/sales` - Create sales transaction
- `GET /api/v1/transactions/sales/:id` - Get sales transaction by ID

#### Dashboard
- `GET /api/v1/dashboard/stats` - Get dashboard statistics

### Business Features:
- ✅ **Purchase Transactions**: Buy vehicles from customers
- ✅ **Sales Transactions**: Sell vehicles to customers
- ✅ **Auto Numbering**: PUR-YYYYMMDD-XXX, SAL-YYYYMMDD-XXX
- ✅ **Invoice Generation**: INV-PUR-YYYYMMDD-XXX, INV-SAL-YYYYMMDD-XXX
- ✅ **Vehicle Status Updates**: Automatic status changes during transactions
- ✅ **Payment Methods**: Cash, Transfer, Check, Credit
- ✅ **Tax & Discount Calculations**: Automatic total calculations
- ✅ **Dashboard Metrics**: Real-time business statistics

### Database Tables:
1. users ✅
2. user_sessions ✅
3. customers ✅
4. vehicles ✅
5. vehicle_images
6. spare_parts
7. purchase_transactions ✅
8. sales_transactions ✅
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
Phase 4: Repair & Parts Management
