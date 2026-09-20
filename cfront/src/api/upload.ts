import { http } from '@/utils/http'

export interface UploadFileResult {
  fileURL: string
}

export function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  return http.post<UploadFileResult>('/api/v1/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}
