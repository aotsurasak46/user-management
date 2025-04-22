import { defineStore } from 'pinia'
import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_BASE_URL 

export const useUserAccountStore = defineStore('userAccount', {
  state: () => ({
    isAuthenticated: false,
    user: null,
  }),

  actions: {
    async login(email, password) {
      try {
        const response = await axios.post(`${BASE_URL}/api/v1/login`, {
            email: email,
            password: password,
        }, {
          withCredentials: true,
        })
        this.isAuthenticated = true
        this.user = response.data
        return true

      } catch (error) {
        this.isAuthenticated = false
        this.user = null
        let errorMessage = 'Something went wrong. Please try again.'

        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },

    async register(email, password, name) {
      try {
        await axios.post(`${BASE_URL}/api/v1/register`, {
          email: email,
          password: password,
          name: name,
        }, {
          withCredentials: true,
        })
        return true
      } catch (error) {
        let errorMessage = 'Something went wrong. Please try again.'

        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },

    async checkAuth() {
        try {
          const response = await axios.get(`${BASE_URL}/api/v1/check-auth`, {
            withCredentials: true,
          })
          this.isAuthenticated = true
          this.user = response.data.user || null
        } catch (err) {
          this.isAuthenticated = false
          this.user = null
          console.error('Auth check failed:', err.status) 
        }
    },
      
    async logout() {
      try {
        await axios.post(`${BASE_URL}/api/v1/logout`, {}, {
          withCredentials: true,
        })
      } catch (err) {
      }
      this.isAuthenticated = false
      this.user = null
    },
  },
})
