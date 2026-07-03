import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const usePaymentStore = defineStore('payment', () => {
  const payments = ref([])
  const loading = ref(false)

  const createPayment = async (paymentData) => {
    loading.value = true
    try {
      const response = await api.post('/payments/create', paymentData)
      payments.value.push(response.data)
      return response.data
    } catch (error) {
      throw error.response?.data?.error || 'Error al crear recaudo'
    } finally {
      loading.value = false
    }
  }

  const getPaymentStatus = async (id) => {
    try {
      const response = await api.get(`/payments/status/${id}`)
      return response.data
    } catch (error) {
      throw error.response?.data?.error || 'Error al obtener estado'
    }
  }

  const verifyPayment = async (reference) => {
    try {
      const response = await api.get(`/payments/verify/${reference}`)
      return response.data
    } catch (error) {
      throw error.response?.data?.error || 'Error al verificar pago'
    }
  }

  return {
    payments,
    loading,
    createPayment,
    getPaymentStatus,
    verifyPayment
  }
})
