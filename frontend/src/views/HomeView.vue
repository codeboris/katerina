<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <nav class="bg-white shadow-sm px-6 py-4 flex justify-between items-center">
      <h1 class="text-xl font-bold text-indigo-700">Katerina · English Club</h1>
      <button
        @click="handleLogout"
        class="text-sm text-gray-500 hover:text-gray-800 transition"
      >
        Logout
      </button>
    </nav>

    <main class="max-w-2xl mx-auto py-12 px-4 space-y-8">
      <MicButton />
      <ResponseCard v-if="voice.result" :result="voice.result" />
      <p v-if="voice.error" class="text-center text-red-500 text-sm">{{ voice.error }}</p>
    </main>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useVoiceStore } from '../stores/voice'
import MicButton from '../components/MicButton.vue'
import ResponseCard from '../components/ResponseCard.vue'

const router = useRouter()
const auth   = useAuthStore()
const voice  = useVoiceStore()

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
