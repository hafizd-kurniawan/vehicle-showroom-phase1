# Vehicle Showroom Management System - Frontend

## Phase 1: Foundation & Authentication ✅

### Features Implemented:
- ✅ React + TypeScript + Vite
- ✅ Clean Architecture (Services, Contexts, Hooks)
- ✅ JWT Authentication with Context API
- ✅ Axios HTTP Client with interceptors
- ✅ shadcn/ui components
- ✅ Tailwind CSS styling
- ✅ React Router for navigation

### Screens Implemented:
- ✅ Login Page with form validation
- ✅ Dashboard Page with role-based UI
- ✅ Authentication flow management
- ✅ Toast notifications

### Architecture:
```
frontend/
├── contexts/           # React Context for state management
├── hooks/             # Custom React hooks
├── services/          # API services
├── pages/             # Page components
├── components/        # Reusable UI components (shadcn/ui)
└── App.tsx           # Main app component
```

### Setup Instructions:
1. Ensure Node.js is installed
2. Backend should be running on http://localhost:8080
3. Frontend will run on the configured port

### Demo Credentials:
- Admin: admin / admin123
- Cashier: cashier / cashier123
- Mechanic: mechanic / mechanic123

### Next Phase:
Phase 2: Customer & Vehicle Management UI
