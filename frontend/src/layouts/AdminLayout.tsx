import React from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from '../components/admin/Sidebar';

const AdminLayout: React.FC = () => {
    return (
        <div className="flex min-h-screen bg-brand-dark">
            <Sidebar />
            <main className="flex-grow p-8 overflow-y-auto">
                <div className="max-w-6xl mx-auto">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};

export default AdminLayout;
