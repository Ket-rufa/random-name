import axios from 'axios';

const httpClient = axios.create({
  baseURL: (import.meta.env.VITE_API_BASE_URL as string) || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

httpClient.interceptors.request.use(
  (config) => {
    try {
      const raw = localStorage.getItem('random_name_wheel_data');
      if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed && parsed.editToken) {
          config.headers['X-Edit-Token'] = parsed.editToken;
        }
      }
    } catch (e) {
      // Ignore parsing errors
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default httpClient;
