import React from 'react';
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
import { Car, Users, Receipt, Wrench, LogOut, User } from 'lucide-react';

export default function DashboardPage() {
  const { user, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
  };

  const dashboardCards = [
    {
      title: 'Vehicles',
      subtitle: 'Coming in Phase 2',
      icon: Car,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      title: 'Customers',
      subtitle: 'Coming in Phase 2',
      icon: Users,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      title: 'Transactions',
      subtitle: 'Coming in Phase 3',
      icon: Receipt,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
    {
      title: 'Repairs',
      subtitle: 'Coming in Phase 4',
      icon: Wrench,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
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
            Phase 1: Authentication & Foundation - Completed ✅
          </p>
        </div>

        {/* Dashboard Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {dashboardCards.map((card, index) => {
            const IconComponent = card.icon;
            return (
              <Card key={index} className="hover:shadow-lg transition-shadow">
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
                <span className="text-sm text-yellow-600 font-semibold">🚧 Next</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 3: Transaction System</span>
                <span className="text-sm text-gray-500">⏳ Planned</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium">Phase 4: Repair & Parts Management</span>
                <span className="text-sm text-gray-500">⏳ Planned</span>
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
