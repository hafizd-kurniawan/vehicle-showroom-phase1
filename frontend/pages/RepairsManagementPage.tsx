import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useToast } from '@/components/ui/use-toast';
import { Search, Plus, ArrowLeft, Wrench, Calendar, User, DollarSign, Package } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { repairService, Repair } from '../services/repairService';
import { sparePartService, SparePart } from '../services/sparePartService';
import { useAuth } from '../contexts/AuthContext';

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  in_progress: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  cancelled: 'bg-red-100 text-red-800',
};

const statusLabels = {
  pending: 'Pending',
  in_progress: 'In Progress',
  completed: 'Completed',
  cancelled: 'Cancelled',
};

export default function RepairsManagementPage() {
  const [repairs, setRepairs] = useState<Repair[]>([]);
  const [spareParts, setSpareParts] = useState<SparePart[]>([]);
  const [repairLoading, setRepairLoading] = useState(true);
  const [sparePartLoading, setSparePartLoading] = useState(true);
  const [repairSearch, setRepairSearch] = useState('');
  const [sparePartSearch, setSparePartSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [repairPage, setRepairPage] = useState(1);
  const [sparePartPage, setSparePartPage] = useState(1);
  const [repairTotal, setRepairTotal] = useState(0);
  const [sparePartTotal, setSparePartTotal] = useState(0);
  const navigate = useNavigate();
  const { toast } = useToast();
  const { user } = useAuth();

  const canCreateRepair = user?.role === 'admin' || user?.role === 'mechanic';
  const canCreateSparePart = user?.role === 'admin' || user?.role === 'mechanic';
  const canManageRepairs = user?.role === 'admin' || user?.role === 'mechanic';

  const limit = 10;

  useEffect(() => {
    loadRepairs();
  }, [repairPage, repairSearch, statusFilter]);

  useEffect(() => {
    loadSpareParts();
  }, [sparePartPage, sparePartSearch]);

  const loadRepairs = async () => {
    try {
      setRepairLoading(true);
      const response = await repairService.list(repairPage, limit, repairSearch, statusFilter);
      setRepairs(response.repairs);
      setRepairTotal(response.total);
    } catch (error) {
      console.error('Failed to load repairs:', error);
      toast({
        title: "Error",
        description: "Failed to load repairs",
        variant: "destructive",
      });
    } finally {
      setRepairLoading(false);
    }
  };

  const loadSpareParts = async () => {
    try {
      setSparePartLoading(true);
      const response = await sparePartService.list(sparePartPage, limit, sparePartSearch);
      setSpareParts(response.spare_parts);
      setSparePartTotal(response.total);
    } catch (error) {
      console.error('Failed to load spare parts:', error);
      toast({
        title: "Error",
        description: "Failed to load spare parts",
        variant: "destructive",
      });
    } finally {
      setSparePartLoading(false);
    }
  };

  const handleRepairSearch = (value: string) => {
    setRepairSearch(value);
    setRepairPage(1);
  };

  const handleSparePartSearch = (value: string) => {
    setSparePartSearch(value);
    setSparePartPage(1);
  };

  const handleStatusFilter = (value: string) => {
    setStatusFilter(value);
    setRepairPage(1);
  };

  const handleStatusChange = async (repairId: number, newStatus: string) => {
    try {
      await repairService.updateStatus(repairId, newStatus);
      toast({
        title: "Success",
        description: "Repair status updated successfully",
      });
      loadRepairs();
    } catch (error) {
      console.error('Failed to update repair status:', error);
      toast({
        title: "Error",
        description: "Failed to update repair status",
        variant: "destructive",
      });
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const repairTotalPages = Math.ceil(repairTotal / limit);
  const sparePartTotalPages = Math.ceil(sparePartTotal / limit);

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
                Repair & Parts Management
              </h1>
            </div>
            <div className="flex gap-2">
              {canCreateSparePart && (
                <Button onClick={() => navigate('/spare-parts/new')} variant="outline">
                  <Package className="h-4 w-4 mr-2" />
                  Add Spare Part
                </Button>
              )}
              {canCreateRepair && (
                <Button onClick={() => navigate('/repairs/new')}>
                  <Wrench className="h-4 w-4 mr-2" />
                  New Repair
                </Button>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Tabs defaultValue="repairs" className="space-y-6">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="repairs">Repairs</TabsTrigger>
            <TabsTrigger value="spare-parts">Spare Parts</TabsTrigger>
          </TabsList>

          {/* Repairs Tab */}
          <TabsContent value="repairs" className="space-y-6">
            {/* Search and Filters */}
            <Card>
              <CardHeader>
                <CardTitle>Search Repairs</CardTitle>
                <CardDescription>Find repairs by number, title, vehicle, or mechanic</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex gap-4">
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                    <Input
                      placeholder="Search repairs..."
                      value={repairSearch}
                      onChange={(e) => handleRepairSearch(e.target.value)}
                      className="pl-10"
                    />
                  </div>
                  <Select value={statusFilter} onValueChange={handleStatusFilter}>
                    <SelectTrigger className="w-48">
                      <SelectValue placeholder="Filter by status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">All Status</SelectItem>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="in_progress">In Progress</SelectItem>
                      <SelectItem value="completed">Completed</SelectItem>
                      <SelectItem value="cancelled">Cancelled</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </Card>

            {/* Repair List */}
            <div className="space-y-4">
              {repairLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="mt-2 text-gray-500">Loading repairs...</p>
                </div>
              ) : repairs.length === 0 ? (
                <Card>
                  <CardContent className="text-center py-8">
                    <p className="text-gray-500">No repairs found</p>
                    {canCreateRepair && (
                      <Button 
                        onClick={() => navigate('/repairs/new')}
                        className="mt-4"
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        Create First Repair
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : (
                repairs.map((repair) => (
                  <Card key={repair.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <Wrench className="h-5 w-5 text-purple-600" />
                            <h3 className="text-lg font-semibold">{repair.repair_number}</h3>
                            <Badge className={statusColors[repair.status]}>
                              {statusLabels[repair.status]}
                            </Badge>
                          </div>
                          
                          <h4 className="font-medium text-gray-900 mb-2">{repair.title}</h4>
                          
                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-600 mb-4">
                            <div className="flex items-center gap-2">
                              <Calendar className="h-4 w-4" />
                              <span>{formatDate(repair.created_at)}</span>
                            </div>
                            {repair.mechanic && (
                              <div className="flex items-center gap-2">
                                <User className="h-4 w-4" />
                                <span>{repair.mechanic.full_name}</span>
                              </div>
                            )}
                            <div className="flex items-center gap-2">
                              <DollarSign className="h-4 w-4" />
                              <span>{formatPrice(repair.total_cost)}</span>
                            </div>
                            {repair.vehicle && (
                              <div>
                                <span className="text-gray-500">Vehicle:</span>
                                <p className="font-medium">
                                  {repair.vehicle.brand} {repair.vehicle.model}
                                </p>
                              </div>
                            )}
                          </div>

                          {repair.description && (
                            <p className="text-sm text-gray-600 mb-3">{repair.description}</p>
                          )}

                          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                            <div>
                              <span className="text-gray-500">Labor Cost:</span>
                              <p className="font-semibold">{formatPrice(repair.labor_cost)}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Parts Cost:</span>
                              <p className="font-semibold">{formatPrice(repair.total_parts_cost)}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Total Cost:</span>
                              <p className="font-bold text-purple-600">{formatPrice(repair.total_cost)}</p>
                            </div>
                          </div>
                        </div>
                        
                        <div className="flex flex-col gap-2 ml-4">
                          {/* Status Change */}
                          {canManageRepairs && (
                            <Select
                              value={repair.status}
                              onValueChange={(value) => handleStatusChange(repair.id, value)}
                            >
                              <SelectTrigger className="w-40">
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="pending">Pending</SelectItem>
                                <SelectItem value="in_progress">In Progress</SelectItem>
                                <SelectItem value="completed">Completed</SelectItem>
                                <SelectItem value="cancelled">Cancelled</SelectItem>
                              </SelectContent>
                            </Select>
                          )}

                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => navigate(`/repairs/${repair.id}`)}
                          >
                            View Details
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* Repair Pagination */}
            {repairTotalPages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setRepairPage(repairPage - 1)}
                  disabled={repairPage === 1}
                >
                  Previous
                </Button>
                <span className="flex items-center px-4 text-sm text-gray-600">
                  Page {repairPage} of {repairTotalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => setRepairPage(repairPage + 1)}
                  disabled={repairPage === repairTotalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </TabsContent>

          {/* Spare Parts Tab */}
          <TabsContent value="spare-parts" className="space-y-6">
            {/* Search */}
            <Card>
              <CardHeader>
                <CardTitle>Search Spare Parts</CardTitle>
                <CardDescription>Find spare parts by name, code, or brand</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                  <Input
                    placeholder="Search spare parts..."
                    value={sparePartSearch}
                    onChange={(e) => handleSparePartSearch(e.target.value)}
                    className="pl-10"
                  />
                </div>
              </CardContent>
            </Card>

            {/* Spare Parts List */}
            <div className="space-y-4">
              {sparePartLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="mt-2 text-gray-500">Loading spare parts...</p>
                </div>
              ) : spareParts.length === 0 ? (
                <Card>
                  <CardContent className="text-center py-8">
                    <p className="text-gray-500">No spare parts found</p>
                    {canCreateSparePart && (
                      <Button 
                        onClick={() => navigate('/spare-parts/new')}
                        className="mt-4"
                      >
                        <Plus className="h-4 w-4 mr-2" />
                        Add First Spare Part
                      </Button>
                    )}
                  </CardContent>
                </Card>
              ) : (
                spareParts.map((part) => (
                  <Card key={part.id} className="hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-3">
                            <Package className="h-5 w-5 text-blue-600" />
                            <h3 className="text-lg font-semibold">{part.name}</h3>
                            <span className="text-sm text-gray-500">#{part.part_code}</span>
                            {part.stock_quantity <= part.min_stock_level && (
                              <Badge variant="destructive">Low Stock</Badge>
                            )}
                          </div>
                          
                          {part.description && (
                            <p className="text-sm text-gray-600 mb-3">{part.description}</p>
                          )}
                          
                          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm text-gray-600 mb-4">
                            {part.brand && (
                              <div>
                                <span className="text-gray-500">Brand:</span>
                                <p className="font-medium">{part.brand}</p>
                              </div>
                            )}
                            <div>
                              <span className="text-gray-500">Stock:</span>
                              <p className="font-medium">
                                {part.stock_quantity} {part.unit_measure || 'pcs'}
                              </p>
                            </div>
                            <div>
                              <span className="text-gray-500">Cost Price:</span>
                              <p className="font-medium">{formatPrice(part.cost_price)}</p>
                            </div>
                            <div>
                              <span className="text-gray-500">Selling Price:</span>
                              <p className="font-medium">{formatPrice(part.selling_price)}</p>
                            </div>
                          </div>

                          <div className="text-sm">
                            <span className="text-gray-500">Min Stock Level:</span>
                            <span className="ml-2 font-medium">{part.min_stock_level}</span>
                          </div>
                        </div>
                        
                        <div className="flex gap-2 ml-4">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => navigate(`/spare-parts/${part.id}/edit`)}
                          >
                            Edit
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            {/* Spare Parts Pagination */}
            {sparePartTotalPages > 1 && (
              <div className="flex justify-center gap-2 mt-8">
                <Button
                  variant="outline"
                  onClick={() => setSparePartPage(sparePartPage - 1)}
                  disabled={sparePartPage === 1}
                >
                  Previous
                </Button>
                <span className="flex items-center px-4 text-sm text-gray-600">
                  Page {sparePartPage} of {sparePartTotalPages}
                </span>
                <Button
                  variant="outline"
                  onClick={() => setSparePartPage(sparePartPage + 1)}
                  disabled={sparePartPage === sparePartTotalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
