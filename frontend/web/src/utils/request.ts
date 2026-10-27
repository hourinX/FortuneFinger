import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
} from 'axios'

interface ApiEnvelope<T> {
  code: number
  data?: T
  message?: string
  msg?: string
}

interface RequestInstance extends Omit<AxiosInstance, 'delete' | 'get' | 'patch' | 'post' | 'put'> {
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T>
  patch<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig<D>): Promise<T>
  post<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig<D>): Promise<T>
  put<T = unknown, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig<D>): Promise<T>
}

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 15_000,
}) as RequestInstance

request.interceptors.response.use(
  (response) => {
    const body = response.data as ApiEnvelope<unknown>
    if (body?.code !== undefined && Number(body.code) !== 0) {
      const error = new Error(body.message || body.msg || '请求失败')
      return Promise.reject(error)
    }
    return (body?.data !== undefined ? body.data : body) as unknown as AxiosResponse
  },
  (error) => {
    const body = error.response?.data as ApiEnvelope<unknown> | undefined
    error.message = body?.message || body?.msg || error.message || '网络请求失败'
    return Promise.reject(error)
  },
)

export default request
