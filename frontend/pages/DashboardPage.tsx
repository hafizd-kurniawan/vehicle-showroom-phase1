import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu';
import { Car, Users, Receipt, Wrench, LogOut, Plus, TrendingUp, TrendingDown, DollarSign, BarChart3 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { transactionService, DashboardStats } from '../services/transactionService';

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardStats();
  }, []);

  const loadDashboardStats = async () => {
    try {
      const dashboardStats = await transactionService.getDashboardStats();
      setStats(dashboardStats);
    } catch (error) {
      console.error('Failed to load dashboard stats:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await logout();
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const dashboardCards = [
    {
      title: 'Vehicles',
      subtitle: 'Manage vehicle inventory',
      icon: Car,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      onClick: () => navigate('/vehicles'),
    },
    {
      title: 'Customers',
      subtitle: 'Manage customer database',
      icon: Users,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
      onClick: () => navigate('/customers'),
    },
    {
      title: 'Transactions',
      subtitle: 'Purchase & Sales transactions',
      icon: Receipt,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
      onClick: () => navigate('/transactions'),
    },
    {
      title: 'Repairs',
      subtitle: 'Coming in Phase 4',
      icon: Wrench,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      onClick: () => {},
    },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <Car className="h-8 w-8 text-blue-600 mr-3" />
              <h1 className="text-xl font-semibold text-gray-900">
                Vehicle Showroom Dashboard
              </h1>
            </div>
            
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarFallback className="bg-blue-600 text-white">
                      {user?.full_name?.[0]?.toUpperCase() || 'U'}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end" forceMount>
                <DropdownMenuLabel className="font-normal">
                  <div className="flex flex-col space-y-1">
                    <p className="text-sm font-medium leading-none">{user?.full_name}</p>
                    <p className="text-xs leading-none text-muted-foreground">
                      {user?.role?.toUpperCase()}
                    </p>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={handleLogout}>
                  <LogOut className="mr-2 h-4 w-4" />
                  <span>Log out</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Welcome to Vehicle Showroom Management System
          </h2>
          <p className="text-lg text-green-600 font-semibold">
            Phase 3: Transaction System - Completed ✅
          </p>
        </div>

        {/* Statistics Cards */}
        {loading ? (
          <div className="text-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-2 text-gray-500">Loading dashboard...</p>
          </div>
        ) : stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Vehicles</CardTitle>
                <Car className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_vehicles}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.vehicles_for_sale} ready to sell
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Customers</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_customers}</div>
                <p className="text-xs text-muted-foreground">
                  Active customers
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Today's Revenue</CardTitle>
                <DollarSign className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatPrice(stats.today_revenue)}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.today_sales} sales today
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Profit</CardTitle>
                <BarChart3 className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{formatPrice(stats.total_profit)}</div>
                <p className="text-xs text-muted-foreground">
                  All time profit
                </p>
              </CardContent>
            </Card>
          </div>
        )}

        {/* Quick Actions */}
        <div className="mb-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <div className="flex gap-4">
            <Button onClick={() => navigate('/vehicles/new')} className="flex items-center gap-2">
              <Plus className="h-4 w-4" />
              Add Vehicle
            </Button>
            <Button onClick={() => navigate('/customers/new')} variant="outline" className="flex items-center gap-2">
              <Plus className="h-4 w-4" />
              Add Customer
            </Button>
            <Button onClick={() => navigate('/transactions/purchase/new')} variant="outline" className="flex items-center gap-2">
              <TrendingDown className="h-4 w-4" />
              New Purchase
            </Button>
            <Button onClick={() => navigate('/transactions/sales/new')} className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4" />
              New Sale
            </Button>
          </div>
        </div>

        {/* Dashboard Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {dashboardCards.map((card, index) => {
            const IconComponent = card.icon;
            return (
              <Card 
                key={index} 
                className="hover:shadow-lg transition-shadow cursor-pointer"
                onClick={card.onClick}
              >
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">
                    {card.title}
                  </CardTitle>
                  <div className={`p-2 rounded-lg ${card.bgColor}`}>
                    <IconComponent className={`h-6 w-6 ${card.color}`} />
                  </div>
                </CardHeader>
                <CardContent>
                  <CardDescription className="text-xs text-muted-foreground">
                    {card.subtitle}
                  </CardDescription>
                </CardContent>
              </Card>
            );
          })}
        </div>

        {/* Status Card */}
        <Card className="mt-8">
          <CardHeader>
            <CardTitle>Development Status</CardTitle>
            <CardDescription>Current implementation progress</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 1: Foundation & Authentication</span>
                <span className="text-sm text-green-600 font-semibold">✅ Completed</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 2: Customer & Vehicle Management</span>
                <span className="text-sm text-green-600 font-semibold">✅ Completed</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 3: Transaction System</span>
                <span className="text-sm text-green-600 font-semibold">✅ Completed</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 4: Repair & Parts Management</span>
                <span className="text-sm text-yellow-600 font-semibold">🚧 Next</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 5: Reporting & Dashboard</span>
                <span className="text-sm text-gray-500">⏳ Planned</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  );
}
