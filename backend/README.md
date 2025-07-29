# Vehicle Showroom Management System - Backend

## Phase 1: Foundation & Authentication ✅

### Features Implemented:
- ✅ Clean Architecture (Entity, Repository, UseCase, Handler)
- ✅ PostgreSQL Database with SQLX
- ✅ JWT Authentication
- ✅ User Management (Admin, Mechanic, Cashier roles)
- ✅ Session Management
- ✅ Complete Database Schema (11 tables)
- ✅ CORS Support for Flutter integration

### API Endpoints:
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Get user profile

### Database Tables Created:
1. users
2. user_sessions
3. customers
4. vehicles
5. vehicle_images
6. spare_parts
7. purchase_transactions
8. sales_transactions
9. repairs
10. repair_parts
11. stock_movements

### Setup Instructions:
1. Install PostgreSQL
2. Create database: `vehicle_showroom`
3. Copy `.env` file and update database credentials
4. Run: `go mod tidy`
5. Run: `go run cmd/main.go`

### Demo Users (will be created in Phase 2):
- Admin: admin / admin123
- Cashier: cashier / cashier123  
- Mechanic: mechanic / mechanic123

### Next Phase:
Phase 2: Customer & Vehicle Management
