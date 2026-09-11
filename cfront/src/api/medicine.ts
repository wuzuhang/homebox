import { http } from '@/utils/http'

export interface Medicine {
  ID: number
  CreatedAt?: string
  UpdatedAt?: string
  user_id: number
  name: string
  manufacturer: string
  disease_ids?: number[]
  stock: number
  unit: string
  min_stock_warn: number
  daily_dose: number
  usage: string
  photo: string
  remark: string
}

export interface MedicinePayload {
  id?: number
  user_id: number
  name: string
  manufacturer: string
  disease_ids: number[]
  stock: number
  unit: string
  min_stock_warn: number
  daily_dose: number
  usage: string
  photo: string
  remark: string
}

export interface Disease {
  id: number
  name: string
  code: string
  description: string
}

export function getMedicines() {
  return http.get<Medicine[]>('/api/v1/medicines')
}

export function addMedicine(data: MedicinePayload) {
  return http.post<void>('/api/v1/medicine', data)
}

export function updateMedicine(data: MedicinePayload) {
  return http.put<void>('/api/v1/medicine', data)
}

export function deleteMedicine(id: number) {
  return http.delete<void>('/api/v1/medicine', { data: { id } })
}

export function getDiseases() {
  return http.get<Disease[]>('/api/v1/diseases')
}
