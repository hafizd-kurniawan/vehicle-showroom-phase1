import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import CustomersPage from './pages/CustomersPage';
import CustomerFormPage from './pages/CustomerFormPage';
import VehiclesPage from './pages/VehiclesPage';
import VehicleFormPage from './pages/VehicleFormPage';
import TransactionsPage from './pages/TransactionsPage';
import PurchaseFormPage from './pages/PurchaseFormPage';
import SalesFormPage from './pages/SalesFormPage';
import ReportsPage from './pages/ReportsPage';
import { Toaster } from '@/components/ui/toaster';

function AppRoutes() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <Routes>
      <Route 
        path="/login" 
        element={user ? <Navigate to="/dashboard" replace /> : <LoginPage />} 
      />
      <Route 
        path="/dashboard" 
        element={user ? <DashboardPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/customers" 
        element={user ? <CustomersPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/customers/new" 
        element={user ? <CustomerFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/customers/:id/edit" 
        element={user ? <CustomerFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/vehicles" 
        element={user ? <VehiclesPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/vehicles/new" 
        element={user ? <VehicleFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/vehicles/:id/edit" 
        element={user ? <VehicleFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/transactions" 
        element={user ? <TransactionsPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/transactions/purchase/new" 
        element={user ? <PurchaseFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/transactions/sales/new" 
        element={user ? <SalesFormPage /> : <Navigate to="/login" replace />} 
      />
      <Route 
        path="/reports" 
        element={user && user.role === 'admin' ? <ReportsPage /> : <Navigate to="/dashboard" replace />} 
      />
      <Route 
        path="/" 
        element={<Navigate to={user ? "/dashboard" : "/login"} replace />} 
      />
    </Routes>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="min-h-screen bg-gray-50">
          <AppRoutes />
          <Toaster />
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
