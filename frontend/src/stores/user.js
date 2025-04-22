import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_BASE_URL 


export const useUserStore = defineStore('user', {
  state: () => ({
    users: [],
    selectedUser: {},
  }),
  actions:{
    async loadUsers() {
      try {
        const response = await axios.get(`${BASE_URL}/api/v1/users` ,{
          withCredentials: true,
        })
        this.users = response.data
      }catch(error){
        console.log('error', error)
        let errorMessage = 'Something went wrong. Please try again.'
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },
    async loadUser(id) {
      try {
        const response = await axios.get(`${BASE_URL}/api/v1/users/${id}`,{
          withCredentials: true,
        })
        this.selectedUser = response.data
      }catch(error){
        console.log('error', error)
        let errorMessage = 'Something went wrong. Please try again.'
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },
    async createUser(name,email, password, role) {
      const bodyData = {
        name: name,
        email: email,
        role: role,
        password : password
      }
      try {
        const response = await axios.post(`${BASE_URL}/api/v1/users`, bodyData, {
          withCredentials: true,
        })
        return true
      } catch (error) {
        console.log('error', error)
        let errorMessage = 'Something went wrong. Please try again.'
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },
    async editUser ( name,email,role, id) {
      try {
        const bodyData = {
          name: name,
          email: email,
          role: role
        }
        const response = await axios.put(`${BASE_URL}/api/v1/users/${id}`, bodyData, {
          withCredentials: true,
        })
        return true
      } catch (error) {
        console.log('error', error)
        let errorMessage = 'Something went wrong. Please try again.'
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },
    async deleteUser (id) {
      try {
        const response = await axios.delete(`${BASE_URL}/api/v1/users/${id}`, {
          withCredentials: true,
        })
        return true
      } catch (error) {
        console.log('error', error)
        let errorMessage = 'Something went wrong. Please try again.'
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage = error.response.data.error
        }
        throw new Error(errorMessage)
      }
    },
  }
})
