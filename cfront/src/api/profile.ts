import type { RegisterPayload, UserProfile } from '@/api/auth'
import { http } from '@/utils/http'

export function getProfile() {
  return http.get<UserProfile>('/api/v1/profile')
}

export function updateProfile(data: RegisterPayload) {
  return http.put<UserProfile>('/api/v1/profile', data)
}
