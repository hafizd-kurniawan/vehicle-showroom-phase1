import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useToast } from '@/components/ui/use-toast';
import { Search, Plus, Edit, Trash2, ArrowLeft, Car, Calendar, Gauge, Fuel, Settings } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { vehicleService, Vehicle } from '../services/vehicleService';
import { useAuth } from '../contexts/AuthContext';

const statusColors = {
  purchased: 'bg-blue-100 text-blue-800',
  in_repair: 'bg-yellow-100 text-yellow-800',
  ready_to_sell: 'bg-green-100 text-green-800',
  reserved: 'bg-purple-100 text-purple-800',
  sold: 'bg-gray-100 text-gray-800',
};

const statusLabels = {
  purchased: 'Purchased',
  in_repair: 'In Repair',
  ready_to_sell: 'Ready to Sell',
  reserved: 'Reserved',
  sold: 'Sold',
};

export default function VehiclesPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const navigate = useNavigate();
  const { toast } = useToast();
  const { user } = useAuth();

  const canCreate = user?.role === 'admin' || user?.role === 'cashier';
  const canEdit = user?.role === 'admin' || user?.role === 'cashier' || user?.role === 'mechanic';
  const canDelete = user?.role === 'admin';
  const canChangeStatus = user?.role === 'admin' || user?.role === 'cashier' || user?.role === 'mechanic';

  const limit = 10;

  useEffect(() => {
    loadVehicles();
  }, [page, search, statusFilter]);

  const loadVehicles = async () => {
    try {
      setLoading(true);
      const response = await vehicleService.list(page, limit, search, statusFilter);
      setVehicles(response.vehicles);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load vehicles:', error);
      toast({
        title: "Error",
        description: "Failed to load vehicles",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (value: string) => {
    setSearch(value);
    setPage(1);
  };

  const handleStatusFilter = (value: string) => {
    setStatusFilter(value);
    setPage(1);
  };

  const handleStatusChange = async (vehicleId: number, newStatus: string) => {
    try {
      await vehicleService.updateStatus(vehicleId, newStatus);
      toast({
        title: "Success",
        description: "Vehicle status updated successfully",
      });
      loadVehicles();
    } catch (error) {
      console.error('Failed to update vehicle status:', error);
      toast({
        title: "Error",
        description: "Failed to update vehicle status",
        variant: "destructive",
      });
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this vehicle?')) {
      return;
    }

    try {
      await vehicleService.delete(id);
      toast({
        title: "Success",
        description: "Vehicle deleted successfully",
      });
      loadVehicles();
    } catch (error) {
      console.error('Failed to delete vehicle:', error);
      toast({
        title: "Error",
        description: "Failed to delete vehicle",
        variant: "destructive",
      });
    }
  };

  const formatPrice = (price?: number) => {
    if (!price) return 'N/A';
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Button
                variant="ghost"
                onClick={() => navigate('/dashboard')}
                className="mr-4"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <h1 className="text-xl font-semibold text-gray-900">
                Vehicle Management
              </h1>
            </div>
            {canCreate && (
              <Button onClick={() => navigate('/vehicles/new')}>
                <Plus className="h-4 w-4 mr-2" />
                Add Vehicle
              </Button>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search and Filters */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Search Vehicles</CardTitle>
            <CardDescription>Find vehicles by brand, model, code, chassis number, or license plate</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="Search vehicles..."
                  value={search}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10"
                />
              </div>
              <Select value={statusFilter} onValueChange={handleStatusFilter}>
                <SelectTrigger className="w-48">
                  <SelectValue placeholder="Filter by status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All Status</SelectItem>
                  <SelectItem value="purchased">Purchased</SelectItem>
                  <SelectItem value="in_repair">In Repair</SelectItem>
                  <SelectItem value="ready_to_sell">Ready to Sell</SelectItem>
                  <SelectItem value="reserved">Reserved</SelectItem>
                  <SelectItem value="sold">Sold</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Vehicle List */}
        <div className="space-y-4">
          {loading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
              <p className="mt-2 text-gray-500">Loading vehicles...</p>
            </div>
          ) : vehicles.length === 0 ? (
            <Card>
              <CardContent className="text-center py-8">
                <p className="text-gray-500">No vehicles found</p>
                {canCreate && (
                  <Button 
                    onClick={() => navigate('/vehicles/new')}
                    className="mt-4"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Add First Vehicle
                  </Button>
                )}
              </CardContent>
            </Card>
          ) : (
            vehicles.map((vehicle) => (
              <Card key={vehicle.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-6">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-3">
                        <Car className="h-5 w-5 text-blue-600" />
                        <h3 className="text-lg font-semibold">
                          {vehicle.brand} {vehicle.model} {vehicle.variant}
                        </h3>
                        <Badge className={statusColors[vehicle.status]}>
                          {statusLabels[vehicle.status]}
                        </Badge>
                        <span className="text-sm text-gray-500">#{vehicle.vehicle_code}</span>
                      </div>
                      
                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-600 mb-4">
                        <div className="flex items-center gap-2">
                          <Calendar className="h-4 w-4" />
                          <span>{vehicle.year}</span>
                        </div>
                        {vehicle.mileage && (
                          <div className="flex items-center gap-2">
                            <Gauge className="h-4 w-4" />
                            <span>{vehicle.mileage.toLocaleString()} km</span>
                          </div>
                        )}
                        {vehicle.fuel_type && (
                          <div className="flex items-center gap-2">
                            <Fuel className="h-4 w-4" />
                            <span className="capitalize">{vehicle.fuel_type}</span>
                          </div>
                        )}
                        {vehicle.transmission && (
                          <div className="flex items-center gap-2">
                            <Settings className="h-4 w-4" />
                            <span className="capitalize">{vehicle.transmission}</span>
                          </div>
                        )}
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                        <div>
                          <span className="text-gray-500">Chassis:</span>
                          <p className="font-mono">{vehicle.chassis_number}</p>
                        </div>
                        {vehicle.license_plate && (
                          <div>
                            <span className="text-gray-500">License Plate:</span>
                            <p className="font-mono">{vehicle.license_plate}</p>
                          </div>
                        )}
                        {vehicle.purchase_price && (
                          <div>
                            <span className="text-gray-500">Purchase Price:</span>
                            <p className="font-semibold">{formatPrice(vehicle.purchase_price)}</p>
                          </div>
                        )}
                      </div>

                      {vehicle.purchased_from_customer && (
                        <div className="mt-3 text-sm text-gray-600">
                          <span className="text-gray-500">Purchased from:</span>
                          <span className="ml-2 font-medium">{vehicle.purchased_from_customer.name}</span>
                        </div>
                      )}
                    </div>
                    
                    <div className="flex flex-col gap-2 ml-4">
                      {/* Status Change */}
                      {canChangeStatus && (
                        <Select
                          value={vehicle.status}
                          onValueChange={(value) => handleStatusChange(vehicle.id, value)}
                        >
                          <SelectTrigger className="w-40">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="purchased">Purchased</SelectItem>
                            <SelectItem value="in_repair">In Repair</SelectItem>
                            <SelectItem value="ready_to_sell">Ready to Sell</SelectItem>
                            <SelectItem value="reserved">Reserved</SelectItem>
                            <SelectItem value="sold">Sold</SelectItem>
                          </SelectContent>
                        </Select>
                      )}

                      <div className="flex gap-2">
                        {canEdit && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => navigate(`/vehicles/${vehicle.id}/edit`)}
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                        )}
                        {canDelete && (
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDelete(vehicle.id)}
                            className="text-red-600 hover:text-red-700"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        )}
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center gap-2 mt-8">
            <Button
              variant="outline"
              onClick={() => setPage(page - 1)}
              disabled={page === 1}
            >
              Previous
            </Button>
            <span className="flex items-center px-4 text-sm text-gray-600">
              Page {page} of {totalPages}
            </span>
            <Button
              variant="outline"
              onClick={() => setPage(page + 1)}
              disabled={page === totalPages}
            >
              Next
            </Button>
          </div>
        )}
      </main>
    </div>
  );
}
