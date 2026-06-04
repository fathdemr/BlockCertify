import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Verify from './pages/Verify';
import AdminLayout from './layouts/AdminLayout';
import Dashboard from './pages/admin/Dashboard';
import Upload from './pages/admin/Upload';
import AdminVerify from './pages/admin/Verify';
import History from './pages/admin/History';

function PublicLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <>
      <Navbar />
      <main>{children}</main>
      <Footer />
    </>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Landing */}
          <Route
            path="/"
            element={
              <PublicLayout>
                <Home />
              </PublicLayout>
            }
          />

          {/* Public verify */}
          <Route path="/dogrula" element={<Verify />} />

          {/* Auth pages (no navbar/footer) */}
          <Route path="/giris" element={<Login />} />
          <Route path="/kaydol" element={<Register />} />

          {/* Protected admin */}
          <Route
            path="/admin"
            element={
              <ProtectedRoute requireAdmin>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Dashboard />} />
            <Route path="yukle" element={<Upload />} />
            <Route path="dogrula" element={<AdminVerify />} />
            <Route path="gecmis" element={<History />} />
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}
