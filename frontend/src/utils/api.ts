import axios from 'axios';

// Create an Axios instance with the base URL from the environment variable
const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL, // Use REACT_APP_API_URL
  headers: {
    'Content-Type': 'application/json',
  },
  // withCredentials: true, // Enable credentials (cookies, authorization headers)
});

// Optional: Add request/response interceptors for global error handling or token injection
apiClient.interceptors.request.use(
  (config) => {
    // Add authentication token to headers if available
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle global errors (e.g., redirect to login on 401)
    if (error.response?.status === 401) {
      window.location.href = '/SignIn'; // Redirect to the SignIn page
    }

    // Handle CORS errors
    if (error.response?.status === 403 && error.response?.data?.message === 'CORS error') {
      console.error('CORS error: Request blocked by browser');
      alert('CORS error: Request blocked by browser. Please check your backend CORS configuration.');
    }

    return Promise.reject(error);
  }
);

export default apiClient;