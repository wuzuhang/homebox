import type { UserProfile } from '@/api/auth'
import { clearAuthToken, getAuthToken, setAuthToken } from '@/utils/auth-token'

const USER_STORAGE_KEY = 'homebox_user'

function readStoredUser(): UserProfile | null {
  try {
    const value = localStorage.getItem(USER_STORAGE_KEY)
    return value ? JSON.parse(value) : null
  }
  catch {
    localStorage.removeItem(USER_STORAGE_KEY)
    return null
  }
}

export const useUserStore = defineStore('user', () => {
  const token = ref(getAuthToken() || '')
  const user = ref<UserProfile | null>(readStoredUser())

  const isLoggedIn = computed(() => Boolean(token.value))
  const displayName = computed(() => user.value?.nickname || user.value?.username || '家庭成员')

  function setUser(profile: UserProfile) {
    user.value = profile
    localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(profile))
  }

  function setSession(nextToken: string, profile: UserProfile) {
    token.value = nextToken
    setAuthToken(nextToken)
    setUser(profile)
  }

  function clearSession() {
    token.value = ''
    user.value = null
    clearAuthToken()
    localStorage.removeItem(USER_STORAGE_KEY)
  }

  return {
    token,
    user,
    isLoggedIn,
    displayName,
    setUser,
    setSession,
    clearSession,
  }
})
