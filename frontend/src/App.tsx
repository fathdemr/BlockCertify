import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Navbar from './components/Navbar';
import Footer from './components/Footer';
import Home from './pages/Home';
import Verify from './pages/Verify';
import About from './pages/About';
import Login from './pages/Login';
import AdminLayout from './layouts/AdminLayout';

// Admin Pages
import Dashboard from './pages/admin/Dashboard';
import Upload from './pages/admin/Upload';
import AdminVerify from './pages/admin/Verify';
import History from './pages/admin/History';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="flex flex-col min-h-screen">
          <Navbar />
          <main className="flex-grow">
            <Routes>
              {/* Public Routes */}
              <Route path="/" element={<Home />} />
              <Route path="/verify" element={<Verify />} />
              <Route path="/about" element={<About />} />
              <Route path="/login" element={<Login />} />

              {/* Protected Admin Routes */}
              <Route
                path="/admin"
                element={
                  <ProtectedRoute requireAdmin>
                    <AdminLayout />
                  </ProtectedRoute>
                }
              >
                <Route index element={<Dashboard />} />
                <Route path="upload" element={<Upload />} />
                <Route path="verify" element={<AdminVerify />} />
                <Route path="history" element={<History />} />
              </Route>
            </Routes>
          </main>
          {/* Hide Footer on Admin Pages */}
          <Routes>
            <Route path="/admin/*" element={null} />
            <Route path="*" element={<Footer />} />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
