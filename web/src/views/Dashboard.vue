<template>
  <div class="min-h-screen bg-gray-100 p-8">
    <div class="container mx-auto">
      <h1 class="text-4xl font-bold text-gray-800 mb-8">Dashboard</h1>

      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-gray-600 text-sm font-bold uppercase">Total Recaudos</h3>
          <p class="text-3xl font-bold text-blue-600">{{ stats.totalRecaudos }}</p>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-gray-600 text-sm font-bold uppercase">Monto Total</h3>
          <p class="text-3xl font-bold text-green-600">${{ formatCurrency(stats.totalAmount) }}</p>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-gray-600 text-sm font-bold uppercase">Exitosos</h3>
          <p class="text-3xl font-bold text-emerald-600">{{ stats.successfulCount }}</p>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-gray-600 text-sm font-bold uppercase">Fallidos</h3>
          <p class="text-3xl font-bold text-red-600">{{ stats.failedCount }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-2xl font-bold text-gray-800 mb-4">Crear Recaudo</h2>
          <form @submit.prevent="createRecaudo">
            <div class="mb-4">
              <label class="block text-gray-700 text-sm font-bold mb-2">Monto</label>
              <input
                v-model.number="recaudoForm.amount"
                type="number"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg"
                required
              />
            </div>

            <div class="mb-4">
              <label class="block text-gray-700 text-sm font-bold mb-2">Descripción</label>
              <input
                v-model="recaudoForm.description"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg"
                required
              />
            </div>

            <div class="mb-4">
              <label class="block text-gray-700 text-sm font-bold mb-2">Email del Cliente</label>
              <input
                v-model="recaudoForm.customer_email"
                type="email"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg"
                required
              />
            </div>

            <button
              type="submit"
              :disabled="loading"
              class="w-full bg-blue-600 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-700 disabled:bg-gray-400"
            >
              {{ loading ? 'Procesando...' : 'Crear Recaudo' }}
            </button>
          </form>
        </div>

        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-2xl font-bold text-gray-800 mb-4">Últimos Recaudos</h2>
          <div class="space-y-3">
            <div v-for="recaudo in recentRecaudos" :key="recaudo.id" class="border-l-4 border-blue-500 p-3 bg-gray-50">
              <p class="text-sm text-gray-600">{{ recaudo.reference_code }}</p>
              <p class="font-bold">${{ formatCurrency(recaudo.amount) }}</p>
              <p class="text-sm" :class="recaudo.status === 'completado' ? 'text-green-600' : 'text-yellow-600'">
                {{ recaudo.status }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import api from '../services/api'

const authStore = useAuthStore()

const stats = ref({
  totalRecaudos: 0,
  totalAmount: 0,
  successfulCount: 0,
  failedCount: 0
})

const recentRecaudos = ref([])

const recaudoForm = ref({
  amount: '',
  description: '',
  customer_email: ''
})

const loading = ref(false)

const formatCurrency = (value) => {
  return new Intl.NumberFormat('es-CO', {
    style: 'currency',
    currency: 'COP'
  }).format(value).replace('COP', '').trim()
}

const fetchStats = async () => {
  try {
    const response = await api.get('/reports/statistics')
    if (response.data.monthly) {
      stats.value = {
        totalRecaudos: response.data.monthly.TotalRecaudos,
        totalAmount: response.data.monthly.TotalAmount,
        successfulCount: response.data.monthly.SuccessfulCount,
        failedCount: response.data.monthly.FailedCount
      }
    }
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
}

const createRecaudo = async () => {
  loading.value = true
  try {
    await api.post('/payments/create', {
      ...recaudoForm.value,
      splits: []
    })
    recaudoForm.value = {
      amount: '',
      description: '',
      customer_email: ''
    }
    await fetchStats()
  } catch (error) {
    console.error('Error creating recaudo:', error)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await fetchStats()
})
</script>
