import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token to headers
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for handling errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, clear localStorage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  signup: async (username, email, password) => {
    const response = await api.post('/auth/signup', {
      username,
      email,
      password,
    });
    return response.data;
  },

  login: async (email, password) => {
    const response = await api.post('/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  getProfile: async () => {
    const response = await api.get('/auth/profile');
    return response.data;
  },
};

// Product API
export const productAPI = {
  getAll: async () => {
    const response = await axios.get('http://localhost:8080/products');
    return response.data;
  },

  getOne: async (id) => {
    const response = await axios.get(`http://localhost:8080/products/${id}`);
    return response.data;
  },

  create: async (product) => {
    const response = await axios.post('http://localhost:8080/products', product);
    return response.data;
  },

  update: async (id, product) => {
    const response = await axios.put(`http://localhost:8080/products/${id}`, product);
    return response.data;
  },

  delete: async (id) => {
    const response = await axios.delete(`http://localhost:8080/products/${id}`);
    return response.data;
  },
};

export default api;
