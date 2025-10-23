import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { productAPI } from '../services/api'
import Products from './Products'

function Dashboard() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  useEffect(() => {
    const userData = localStorage.getItem('user')
    if (userData) {
      setUser(JSON.parse(userData))
    }
    setLoading(false)
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/login')
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  return (
    <div style={{ width: '100%', maxWidth: '1200px' }}>
      <div className="dashboard-header">
        <div>
          <h2>Welcome, {user?.username}! 👋</h2>
          <p className="user-info">Email: {user?.email}</p>
        </div>
        <button className="logout-btn" onClick={handleLogout}>
          Logout
        </button>
      </div>

      <Products />
    </div>
  )
}

export default Dashboard