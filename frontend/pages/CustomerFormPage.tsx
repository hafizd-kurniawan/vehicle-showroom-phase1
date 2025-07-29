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
import { customerService, CreateCustomerRequest, UpdateCustomerRequest, Customer } from '../services/customerService';

export default function CustomerFormPage() {
  const [formData, setFormData] = useState<CreateCustomerRequest>({
    name: '',
    phone: '',
    email: '',
    address: '',
    id_card_number: '',
    type: 'individual',
  });
  const [loading, setLoading] = useState(false);
  const [isEdit, setIsEdit] = useState(false);
  const navigate = useNavigate();
  const { id } = useParams();
  const { toast } = useToast();

  useEffect(() => {
    if (id) {
      setIsEdit(true);
      loadCustomer(parseInt(id));
    }
  }, [id]);

  const loadCustomer = async (customerId: number) => {
    try {
      setLoading(true);
      const customer = await customerService.getById(customerId);
      setFormData({
        name: customer.name,
        phone: customer.phone || '',
        email: customer.email || '',
        address: customer.address || '',
        id_card_number: customer.id_card_number || '',
        type: customer.type,
      });
    } catch (error) {
      console.error('Failed to load customer:', error);
      toast({
        title: "Error",
        description: "Failed to load customer",
        variant: "destructive",
      });
      navigate('/customers');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      toast({
        title: "Error",
        description: "Customer name is required",
        variant: "destructive",
      });
      return;
    }

    try {
      setLoading(true);
      
      const submitData = {
        ...formData,
        phone: formData.phone || undefined,
        email: formData.email || undefined,
        address: formData.address || undefined,
        id_card_number: formData.id_card_number || undefined,
      };

      if (isEdit && id) {
        await customerService.update(parseInt(id), submitData as UpdateCustomerRequest);
        toast({
          title: "Success",
          description: "Customer updated successfully",
        });
      } else {
        await customerService.create(submitData);
        toast({
          title: "Success",
          description: "Customer created successfully",
        });
      }
      
      navigate('/customers');
    } catch (error) {
      console.error('Failed to save customer:', error);
      toast({
        title: "Error",
        description: `Failed to ${isEdit ? 'update' : 'create'} customer`,
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: keyof CreateCustomerRequest, value: string) => {
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
                onClick={() => navigate('/customers')}
                className="mr-4"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <h1 className="text-xl font-semibold text-gray-900">
                {isEdit ? 'Edit Customer' : 'Add New Customer'}
              </h1>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Card>
          <CardHeader>
            <CardTitle>{isEdit ? 'Edit Customer' : 'Add New Customer'}</CardTitle>
            <CardDescription>
              {isEdit ? 'Update customer information' : 'Enter customer details to add them to the system'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Customer Type */}
              <div className="space-y-2">
                <Label htmlFor="type">Customer Type *</Label>
                <Select
                  value={formData.type}
                  onValueChange={(value: 'individual' | 'corporate') => handleInputChange('type', value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select customer type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="individual">Individual</SelectItem>
                    <SelectItem value="corporate">Corporate</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Name */}
              <div className="space-y-2">
                <Label htmlFor="name">
                  {formData.type === 'corporate' ? 'Company Name' : 'Full Name'} *
                </Label>
                <Input
                  id="name"
                  type="text"
                  placeholder={formData.type === 'corporate' ? 'Enter company name' : 'Enter full name'}
                  value={formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  required
                />
              </div>

              {/* Phone */}
              <div className="space-y-2">
                <Label htmlFor="phone">Phone Number</Label>
                <Input
                  id="phone"
                  type="tel"
                  placeholder="Enter phone number"
                  value={formData.phone}
                  onChange={(e) => handleInputChange('phone', e.target.value)}
                />
              </div>

              {/* Email */}
              <div className="space-y-2">
                <Label htmlFor="email">Email Address</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter email address"
                  value={formData.email}
                  onChange={(e) => handleInputChange('email', e.target.value)}
                />
              </div>

              {/* ID Card Number */}
              <div className="space-y-2">
                <Label htmlFor="id_card_number">
                  {formData.type === 'corporate' ? 'Tax ID / Business Registration' : 'ID Card Number'}
                </Label>
                <Input
                  id="id_card_number"
                  type="text"
                  placeholder={formData.type === 'corporate' ? 'Enter tax ID or business registration' : 'Enter ID card number'}
                  value={formData.id_card_number}
                  onChange={(e) => handleInputChange('id_card_number', e.target.value)}
                />
              </div>

              {/* Address */}
              <div className="space-y-2">
                <Label htmlFor="address">Address</Label>
                <Textarea
                  id="address"
                  placeholder="Enter complete address"
                  value={formData.address}
                  onChange={(e) => handleInputChange('address', e.target.value)}
                  rows={3}
                />
              </div>

              {/* Submit Button */}
              <div className="flex gap-4 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate('/customers')}
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
                      {isEdit ? 'Update Customer' : 'Create Customer'}
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
