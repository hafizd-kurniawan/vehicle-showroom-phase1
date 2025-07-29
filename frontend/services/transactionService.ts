import { apiClient } from './apiClient';
import { Customer } from './customerService';
import { Vehicle } from './vehicleService';

interface User {
  id: number;
  username: string;
  full_name: string;
  role: string;
}

export interface PurchaseTransaction {
  id: number;
  transaction_number: string;
  invoice_number: string;
  vehicle_id: number;
  customer_id: number;
  vehicle_price: number;
  tax_amount: number;
  total_amount: number;
  payment_method: string;
  payment_reference?: string;
  transaction_date: string;
  cashier_id: number;
  status: string;
  notes?: string;
  created_at: string;
  vehicle?: Vehicle;
  customer?: Customer;
  cashier?: User;
}

export interface SalesTransaction {
  id: number;
  transaction_number: string;
  invoice_number: string;
  vehicle_id: number;
  customer_id: number;
  vehicle_price: number;
  tax_amount: number;
  discount_amount: number;
  total_amount: number;
  payment_method: string;
  payment_reference?: string;
  transaction_date: string;
  cashier_id: number;
  status: string;
  notes?: string;
  created_at: string;
  vehicle?: Vehicle;
  customer?: Customer;
  cashier?: User;
}

export interface CreatePurchaseTransactionRequest {
  vehicle_id: number;
  customer_id: number;
  vehicle_price: number;
  tax_amount: number;
  payment_method: 'cash' | 'transfer' | 'check';
  payment_reference?: string;
  notes?: string;
}

export interface CreateSalesTransactionRequest {
  vehicle_id: number;
  customer_id: number;
  vehicle_price: number;
  tax_amount: number;
  discount_amount: number;
  payment_method: 'cash' | 'transfer' | 'check' | 'credit';
  payment_reference?: string;
  notes?: string;
}

export interface TransactionListResponse {
  transactions: PurchaseTransaction[] | SalesTransaction[];
  total: number;
  page: number;
  limit: number;
}

export interface DashboardStats {
  total_vehicles: number;
  vehicles_for_sale: number;
  vehicles_in_repair: number;
  vehicles_sold: number;
  total_customers: number;
  today_purchases: number;
  today_sales: number;
  today_revenue: number;
  monthly_revenue: number;
  total_profit: number;
}

export const transactionService = {
  // Purchase Transactions
  async listPurchases(page = 1, limit = 10, search = ''): Promise<TransactionListResponse> {
    const response = await apiClient.get('/transactions/purchases', {
      params: { page, limit, search }
    });
    return response.data.data;
  },

  async getPurchaseById(id: number): Promise<PurchaseTransaction> {
    const response = await apiClient.get(`/transactions/purchases/${id}`);
    return response.data.data;
  },

  async createPurchase(data: CreatePurchaseTransactionRequest): Promise<PurchaseTransaction> {
    const response = await apiClient.post('/transactions/purchases', data);
    return response.data.data;
  },

  // Sales Transactions
  async listSales(page = 1, limit = 10, search = ''): Promise<TransactionListResponse> {
    const response = await apiClient.get('/transactions/sales', {
      params: { page, limit, search }
    });
    return response.data.data;
  },

  async getSalesById(id: number): Promise<SalesTransaction> {
    const response = await apiClient.get(`/transactions/sales/${id}`);
    return response.data.data;
  },

  async createSales(data: CreateSalesTransactionRequest): Promise<SalesTransaction> {
    const response = await apiClient.post('/transactions/sales', data);
    return response.data.data;
  },

  // Dashboard
  async getDashboardStats(): Promise<DashboardStats> {
    const response = await apiClient.get('/dashboard/stats');
    return response.data.data;
  },
};
