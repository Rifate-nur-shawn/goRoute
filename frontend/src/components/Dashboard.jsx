import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Products from './Products';

function Dashboard() {
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    // Check if user is authenticated
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    // Get user from localStorage
    const userData = localStorage.getItem('user');
    if (userData) {
      try {
        setUser(JSON.parse(userData));
      } catch (err) {
        console.error('Error parsing user data:', err);
      }
    }
  }, [navigate]);

  const handleLogout = () => {
    // Clear localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    // Redirect to login
    navigate('/login');
  };

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <div className="header-content">
          <h1>Welcome, {user?.username || 'User'}!</h1>
          <button onClick={handleLogout} className="btn btn-logout">
            Logout
          </button>
        </div>
      </header>

      <main className="dashboard-main">
        <div className="products-section">
          <h2>Our Products</h2>
          <Products />
        </div>
      </main>
    </div>
  );
}

export default Dashboard;
