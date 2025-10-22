import { Navigate } from 'react-router-dom';

function ProtectedRoute({ children }) {
  const token = localStorage.getItem('token');

  if (!token) {
    // User not authenticated, redirect to login
    return <Navigate to="/login" replace />;
  }

  // User is authenticated, render the component
  return children;
}

export default ProtectedRoute;
