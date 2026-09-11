import { http } from '@/utils/http'

export interface LoginPayload {
  username: string
  password: string
}

export interface RegisterPayload extends LoginPayload {
  phone: string
  email: string
  nickname: string
  avatar: string
  gender: number
  birthday: string
  weight: number
  disease_ids: number[]
  remarks: string
}

export interface UserProfile {
  ID: number
  CreatedAt?: string
  UpdatedAt?: string
  DeletedAt?: string | null
  username: string
  phone: string
  email: string
  nickname: string
  avatar: string
  gender: number
  birthday?: string | null
  weight: number
  disease_ids?: number[]
  remarks: string
}

export interface AuthResult {
  token: string
  data: UserProfile
}

export function login(data: LoginPayload) {
  return http.post<AuthResult>('/api/v1/login', data)
}

export function register(data: RegisterPayload) {
  return http.post<AuthResult>('/api/v1/register', data)
}
