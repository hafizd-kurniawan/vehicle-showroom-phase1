import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/components/ui/use-toast';
import { ArrowLeft, Save } from 'lucide-react';
import { useNavigate, useParams } from 'react-router-dom';
import { vehicleService, CreateVehicleRequest, UpdateVehicleRequest } from '../services/vehicleService';
import { customerService, Customer } from '../services/customerService';

export default function VehicleFormPage() {
  const [formData, setFormData] = useState<CreateVehicleRequest>({
    chassis_number: '',
    license_plate: '',
    brand: '',
    model: '',
    variant: '',
    year: new Date().getFullYear(),
    color: '',
    mileage: 0,
    fuel_type: '',
    transmission: '',
    purchase_price: 0,
    purchased_from_customer_id: undefined,
    purchase_notes: '',
    condition_notes: '',
  });
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [loading, setLoading] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const navigate = useNavigate();
  const { id } = useParams();
  const { toast } = useToast();

  useEffect(() => {
    loadCustomers();
    if (id) {
      setIsEdit(true);
      loadVehicle(parseInt(id));
    }
  }, [id]);

  const loadCustomers = async () => {
    try {
      const response = await customerService.list(1, 100);
      setCustomers(response.customers);
    } catch (error) {
      console.error('Failed to load customers:', error);
    }
  };

  const loadVehicle = async (vehicleId: number) => {
    try {
      setLoading(true);
      const vehicle = await vehicleService.getById(vehicleId);
      setFormData({
        chassis_number: vehicle.chassis_number,
        license_plate: vehicle.license_plate || '',
        brand: vehicle.brand,
        model: vehicle.model,
        variant: vehicle.variant || '',
        year: vehicle.year,
        color: vehicle.color || '',
        mileage: vehicle.mileage || 0,
        fuel_type: vehicle.fuel_type || '',
        transmission: vehicle.transmission || '',
        purchase_price: vehicle.purchase_price || 0,
        purchased_from_customer_id: vehicle.purchased_from_customer_id,
        purchase_notes: vehicle.purchase_notes || '',
        condition_notes: vehicle.condition_notes || '',
      });
    } catch (error) {
      console.error('Failed to load vehicle:', error);
      toast({
        title: "Error",
        description: "Failed to load vehicle",
        variant: "destructive",
      });
      navigate('/vehicles');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.chassis_number.trim() || !formData.brand.trim() || !formData.model.trim()) {
      toast({
        title: "Error",
        description: "Chassis number, brand, and model are required",
        variant: "destructive",
      });
      return;
    }

    try {
      setLoading(true);
      
      const submitData = {
        ...formData,
        license_plate: formData.license_plate || undefined,
        variant: formData.variant || undefined,
        color: formData.color || undefined,
        mileage: formData.mileage || undefined,
        fuel_type: formData.fuel_type || undefined,
        transmission: formData.transmission || undefined,
        purchase_price: formData.purchase_price || undefined,
        purchased_from_customer_id: formData.purchased_from_customer_id || undefined,
        purchase_notes: formData.purchase_notes || undefined,
        condition_notes: formData.condition_notes || undefined,
      };

      if (isEdit && id) {
        await vehicleService.update(parseInt(id), submitData as UpdateVehicleRequest);
        toast({
          title: "Success",
          description: "Vehicle updated successfully",
        });
      } else {
        await vehicleService.create(submitData);
        toast({
          title: "Success",
          description: "Vehicle created successfully",
        });
      }
      
      navigate('/vehicles');
    } catch (error) {
      console.error('Failed to save vehicle:', error);
      toast({
        title: "Error",
        description: `Failed to ${isEdit ? 'update' : 'create'} vehicle`,
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: keyof CreateVehicleRequest, value: string | number) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Button
                variant="ghost"
                onClick={() => navigate('/vehicles')}
                className="mr-4"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <h1 className="text-xl font-semibold text-gray-900">
                {isEdit ? 'Edit Vehicle' : 'Add New Vehicle'}
              </h1>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Card>
          <CardHeader>
            <CardTitle>{isEdit ? 'Edit Vehicle' : 'Add New Vehicle'}</CardTitle>
            <CardDescription>
              {isEdit ? 'Update vehicle information' : 'Enter vehicle details to add them to the inventory'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Basic Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="chassis_number">Chassis Number *</Label>
                  <Input
                    id="chassis_number"
                    type="text"
                    placeholder="Enter chassis number"
                    value={formData.chassis_number}
                    onChange={(e) => handleInputChange('chassis_number', e.target.value)}
                    required
                    disabled={isEdit}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="license_plate">License Plate</Label>
                  <Input
                    id="license_plate"
                    type="text"
                    placeholder="Enter license plate"
                    value={formData.license_plate}
                    onChange={(e) => handleInputChange('license_plate', e.target.value)}
                  />
                </div>
              </div>

              {/* Vehicle Details */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="brand">Brand *</Label>
                  <Input
                    id="brand"
                    type="text"
                    placeholder="Enter brand"
                    value={formData.brand}
                    onChange={(e) => handleInputChange('brand', e.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="model">Model *</Label>
                  <Input
                    id="model"
                    type="text"
                    placeholder="Enter model"
                    value={formData.model}
                    onChange={(e) => handleInputChange('model', e.target.value)}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="variant">Variant</Label>
                  <Input
                    id="variant"
                    type="text"
                    placeholder="Enter variant"
                    value={formData.variant}
                    onChange={(e) => handleInputChange('variant', e.target.value)}
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="year">Year *</Label>
                  <Input
                    id="year"
                    type="number"
                    min="1900"
                    max="2030"
                    value={formData.year}
                    onChange={(e) => handleInputChange('year', parseInt(e.target.value))}
                    required
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="color">Color</Label>
                  <Input
                    id="color"
                    type="text"
                    placeholder="Enter color"
                    value={formData.color}
                    onChange={(e) => handleInputChange('color', e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="mileage">Mileage (km)</Label>
                  <Input
                    id="mileage"
                    type="number"
                    min="0"
                    placeholder="Enter mileage"
                    value={formData.mileage}
                    onChange={(e) => handleInputChange('mileage', parseInt(e.target.value) || 0)}
                  />
                </div>
              </div>

              {/* Technical Specifications */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="fuel_type">Fuel Type</Label>
                  <Select
                    value={formData.fuel_type}
                    onValueChange={(value) => handleInputChange('fuel_type', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select fuel type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="gasoline">Gasoline</SelectItem>
                      <SelectItem value="diesel">Diesel</SelectItem>
                      <SelectItem value="electric">Electric</SelectItem>
                      <SelectItem value="hybrid">Hybrid</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="transmission">Transmission</Label>
                  <Select
                    value={formData.transmission}
                    onValueChange={(value) => handleInputChange('transmission', value)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select transmission" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="manual">Manual</SelectItem>
                      <SelectItem value="automatic">Automatic</SelectItem>
                      <SelectItem value="cvt">CVT</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Purchase Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="purchase_price">Purchase Price (IDR)</Label>
                  <Input
                    id="purchase_price"
                    type="number"
                    min="0"
                    placeholder="Enter purchase price"
                    value={formData.purchase_price}
                    onChange={(e) => handleInputChange('purchase_price', parseFloat(e.target.value) || 0)}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="purchased_from_customer_id">Purchased From Customer</Label>
                  <Select
                    value={formData.purchased_from_customer_id?.toString() || ''}
                    onValueChange={(value) => handleInputChange('purchased_from_customer_id', value ? parseInt(value) : undefined)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select customer" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">No customer selected</SelectItem>
                      {customers.map((customer) => (
                        <SelectItem key={customer.id} value={customer.id.toString()}>
                          {customer.name} ({customer.customer_code})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>

              {/* Notes */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="purchase_notes">Purchase Notes</Label>
                  <Textarea
                    id="purchase_notes"
                    placeholder="Enter purchase notes"
                    value={formData.purchase_notes}
                    onChange={(e) => handleInputChange('purchase_notes', e.target.value)}
                    rows={3}
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="condition_notes">Condition Notes</Label>
                  <Textarea
                    id="condition_notes"
                    placeholder="Enter condition notes"
                    value={formData.condition_notes}
                    onChange={(e) => handleInputChange('condition_notes', e.target.value)}
                    rows={3}
                  />
                </div>
              </div>

              {/* Submit Button */}
              <div className="flex gap-4 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate('/vehicles')}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={loading}
                  className="flex-1"
                >
                  {loading ? (
                    "Saving..."
                  ) : (
                    <>
                      <Save className="h-4 w-4 mr-2" />
                      {isEdit ? 'Update Vehicle' : 'Create Vehicle'}
                    </>
                  )}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
