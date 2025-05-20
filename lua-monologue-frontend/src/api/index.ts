import axios from 'axios';
import { getConfigFileParsingDiagnostics } from 'typescript';

const api = axios.create({
    baseURL: 'http://localhost:8080',
    headers: { 'Content-Type': 'application/json' },
    withCredentials: true,
});


api.interceptors.response.use(
    res => res,
    async err => {
        const originalRequest = err.config

        if (originalRequest.url === '/refresh') {
            localStorage.removeItem('access_token')
            return Promise.reject(err)
        }
        if (err.response?.status ===401) {
            try {
                const res = await api.post('/refresh')
                const newToken = res.data.access_token

                localStorage.setItem('access_token', newToken)
                
                err.config.headers['Authorization'] = `Bearer ${newToken}`
                return api.request(err.config)
            } catch (refreshErr) {
                window.location.href = '/login'
                return Promise.reject(refreshErr)
            }
        }

        return Promise.reject(err)
    }
)


api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})


export default api;