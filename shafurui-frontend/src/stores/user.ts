import { ref } from 'vue'
import { defineStore } from 'pinia'
import { request } from '@/utils/fetch'

interface AuthTokenResponse {
  accessToken: string
  refreshToken: string
}

const ACCESS_TOKEN_KEY = 'accessToken'
const REFRESH_TOKEN_KEY = 'refreshToken'

export const useUserStore = defineStore('user', () => {
  const accessToken = ref(localStorage.getItem(ACCESS_TOKEN_KEY) || '')
  const refreshToken = ref(localStorage.getItem(REFRESH_TOKEN_KEY) || '')

  function setTokens(tokens: AuthTokenResponse) {
    accessToken.value = tokens.accessToken
    refreshToken.value = tokens.refreshToken
    localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) {
      return Promise.reject(new Error('Refresh token missing'))
    }

    const response = await request<AuthTokenResponse>({
      method: 'POST',
      url: '/auth/refresh',
      data: {
        refreshToken: refreshToken.value,
      },
    })

    setTokens(response)
    return response.accessToken
  }

  function logout() {
    accessToken.value = ''
    refreshToken.value = ''
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  }

  return {
    accessToken,
    refreshToken,
    setTokens,
    refreshAccessToken,
    logout,
  }
})
