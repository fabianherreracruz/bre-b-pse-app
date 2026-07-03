<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
    <div class="bg-white rounded-lg shadow-xl p-8 w-full max-w-md">
      <h2 class="text-3xl font-bold text-gray-800 mb-6">Iniciar Sesión</h2>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Email</label>
          <input
            v-model="form.email"
            type="email"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="tu@email.com"
            required
          />
        </div>

        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2">Contraseña</label>
          <input
            v-model="form.password"
            type="password"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="••••••••"
            required
          />
        </div>

        <div v-if="error" class="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
          {{ error }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-600 text-white font-bold py-2 px-4 rounded-lg hover:bg-blue-700 disabled:bg-gray-400"
        >
          {{ loading ? 'Cargando...' : 'Iniciar Sesión' }}
        </button>
      </form>

      <p class="text-center text-gray-600 mt-4">
        ¿No tienes cuenta?
        <router-link to="/register" class="text-blue-600 hover:text-blue-700 font-bold">
          Regístrate aquí
        </router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  email: '',
  password: ''
})

const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(form.email, form.password)
    router.push('/dashboard')
  } catch (err) {
    error.value = err
  } finally {
    loading.value = false
  }
}
</script>
