import { apiClient } from './apiClient';
import { Customer } from './customerService';

export interface Vehicle {
  id: number;
  vehicle_code: string;
  chassis_number: string;
  license_plate?: string;
  brand: string;
  model: string;
  variant?: string;
  year: number;
  color?: string;
  mileage?: number;
  fuel_type?: string;
  transmission?: string;
  purchase_price?: number;
  total_repair_cost: number;
  suggested_selling_price?: number;
  approved_selling_price?: number;
  final_selling_price?: number;
  status: 'purchased' | 'in_repair' | 'ready_to_sell' | 'reserved' | 'sold';
  purchased_from_customer_id?: number;
  sold_to_customer_id?: number;
  purchased_by_cashier?: number;
  sold_by_cashier?: number;
  price_approved_by_admin?: number;
  purchased_at?: string;
  sold_at?: string;
  created_at: string;
  updated_at: string;
  purchase_notes?: string;
  condition_notes?: string;
  purchased_from_customer?: Customer;
  sold_to_customer?: Customer;
}

export interface CreateVehicleRequest {
  chassis_number: string;
  license_plate?: string;
  brand: string;
  model: string;
  variant?: string;
  year: number;
  color?: string;
  mileage?: number;
  fuel_type?: string;
  transmission?: string;
  purchase_price?: number;
  purchased_from_customer_id?: number;
  purchase_notes?: string;
  condition_notes?: string;
}

export interface UpdateVehicleRequest {
  license_plate?: string;
  brand: string;
  model: string;
  variant?: string;
  year: number;
  color?: string;
  mileage?: number;
  fuel_type?: string;
  transmission?: string;
  suggested_selling_price?: number;
  purchase_notes?: string;
  condition_notes?: string;
}

export interface VehicleListResponse {
  vehicles: Vehicle[];
  total: number;
  page: number;
  limit: number;
}

export const vehicleService = {
  async list(page = 1, limit = 10, search = '', status = ''): Promise<VehicleListResponse> {
    const response = await apiClient.get('/vehicles', {
      params: { page, limit, search, status }
    });
    return response.data.data;
  },

  async getById(id: number): Promise<Vehicle> {
    const response = await apiClient.get(`/vehicles/${id}`);
    return response.data.data;
  },

  async create(data: CreateVehicleRequest): Promise<Vehicle> {
    const response = await apiClient.post('/vehicles', data);
    return response.data.data;
  },

  async update(id: number, data: UpdateVehicleRequest): Promise<Vehicle> {
    const response = await apiClient.put(`/vehicles/${id}`, data);
    return response.data.data;
  },

  async updateStatus(id: number, status: string): Promise<Vehicle> {
    const response = await apiClient.put(`/vehicles/${id}/status`, { status });
    return response.data.data;
  },

  async delete(id: number): Promise<void> {
    await apiClient.delete(`/vehicles/${id}`);
  },
};
