import { apiClient } from './apiClient';
import { Vehicle } from './vehicleService';
import { SparePart } from './sparePartService';

interface User {
  id: number;
  username: string;
  full_name: string;
  role: string;
}

export interface RepairPart {
  id: number;
  repair_id: number;
  spare_part_id: number;
  quantity_used: number;
  unit_cost: number;
  total_cost: number;
  used_at: string;
  notes?: string;
  spare_part?: SparePart;
}

export interface Repair {
  id: number;
  repair_number: string;
  vehicle_id: number;
  title: string;
  description?: string;
  labor_cost: number;
  total_parts_cost: number;
  total_cost: number;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  mechanic_id?: number;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  work_notes?: string;
  vehicle?: Vehicle;
  mechanic?: User;
  repair_parts?: RepairPart[];
}

export interface CreateRepairRequest {
  vehicle_id: number;
  title: string;
  description?: string;
  mechanic_id?: number;
}

export interface UpdateRepairRequest {
  title: string;
  description?: string;
  labor_cost?: number;
  mechanic_id?: number;
  work_notes?: string;
}

export interface AddPartToRepairRequest {
  spare_part_id: number;
  quantity: number;
}

export interface RepairListResponse {
  repairs: Repair[];
  total: number;
  page: number;
  limit: number;
}

export const repairService = {
  async list(page = 1, limit = 10, search = '', status = ''): Promise<RepairListResponse> {
    const response = await apiClient.get('/repairs', {
      params: { page, limit, search, status }
    });
    return response.data.data;
  },

  async getById(id: number): Promise<Repair> {
    const response = await apiClient.get(`/repairs/${id}`);
    return response.data.data;
  },

  async create(data: CreateRepairRequest): Promise<Repair> {
    const response = await apiClient.post('/repairs', data);
    return response.data.data;
  },

  async update(id: number, data: UpdateRepairRequest): Promise<Repair> {
    const response = await apiClient.put(`/repairs/${id}`, data);
    return response.data.data;
  },

  async updateStatus(id: number, status: string): Promise<Repair> {
    const response = await apiClient.put(`/repairs/${id}/status`, { status });
    return response.data.data;
  },

  async addPart(id: number, data: AddPartToRepairRequest): Promise<RepairPart> {
    const response = await apiClient.post(`/repairs/${id}/parts`, data);
    return response.data.data;
  },

  async removePart(repairId: number, partId: number): Promise<void> {
    await apiClient.delete(`/repairs/${repairId}/parts/${partId}`);
  },
};
