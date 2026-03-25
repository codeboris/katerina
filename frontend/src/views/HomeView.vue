<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <nav class="bg-white shadow-sm px-6 py-4 flex justify-between items-center">
      <h1 class="text-xl font-bold text-indigo-700">Katerina · English Club</h1>
      <div class="flex items-center gap-4">
        <button
          @click="showSettings = true"
          class="text-gray-400 hover:text-indigo-600 transition"
          title="Settings"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06
                     a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09
                     A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83
                     l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09
                     A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83
                     l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09
                     a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83
                     l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09
                     a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </button>
        <button
          @click="handleLogout"
          class="text-sm text-gray-500 hover:text-gray-800 transition"
        >
          Logout
        </button>
      </div>
    </nav>

    <main class="max-w-2xl mx-auto py-12 px-4 space-y-8">
      <MicButton />
      <ResponseCard v-if="voice.result" :result="voice.result" />
      <p v-if="voice.error" class="text-center text-red-500 text-sm">{{ voice.error }}</p>
    </main>

    <SettingsModal v-if="showSettings" @close="showSettings = false" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useVoiceStore } from '../stores/voice'
import MicButton from '../components/MicButton.vue'
import ResponseCard from '../components/ResponseCard.vue'
import SettingsModal from '../components/SettingsModal.vue'

const router      = useRouter()
const auth        = useAuthStore()
const voice       = useVoiceStore()
const showSettings = ref(false)

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
