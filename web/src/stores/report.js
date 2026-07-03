import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useReportStore = defineStore('report', () => {
  const reports = ref(null)
  const statistics = ref(null)
  const loading = ref(false)

  const getReportByDateRange = async (startDate, endDate) => {
    loading.value = true
    try {
      const response = await api.get('/reports/by-date', {
        params: { start_date: startDate, end_date: endDate }
      })
      reports.value = response.data
      return response.data
    } catch (error) {
      throw error.response?.data?.error || 'Error al obtener reporte'
    } finally {
      loading.value = false
    }
  }

  const getStatistics = async () => {
    loading.value = true
    try {
      const response = await api.get('/reports/statistics')
      statistics.value = response.data
      return response.data
    } catch (error) {
      throw error.response?.data?.error || 'Error al obtener estadísticas'
    } finally {
      loading.value = false
    }
  }

  const exportToExcel = async (startDate, endDate) => {
    try {
      const response = await api.get('/reports/export-excel', {
        params: { start_date: startDate, end_date: endDate },
        responseType: 'blob'
      })
      
      // Crear descarga
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `reporte_${startDate}_${endDate}.xlsx`)
      document.body.appendChild(link)
      link.click()
      link.parentURL.removeChild(link)
    } catch (error) {
      throw error.response?.data?.error || 'Error al exportar reporte'
    }
  }

  return {
    reports,
    statistics,
    loading,
    getReportByDateRange,
    getStatistics,
    exportToExcel
  }
})
