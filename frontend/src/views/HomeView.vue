<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
    <nav class="bg-white dark:bg-gray-800 shadow-sm px-6 py-4 flex justify-between items-center">
      <h1 class="text-lg sm:text-xl font-bold text-indigo-700 dark:text-indigo-400">{{ $t('nav.title') }}</h1>
      <div class="flex items-center gap-3 sm:gap-4">
        <!-- Theme toggle -->
        <button
          @click="toggleTheme()"
          class="text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition"
          :title="isDark ? $t('nav.theme.switchToLight') : $t('nav.theme.switchToDark')"
        >
          <!-- Sun icon (shown in dark mode) -->
          <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
          <!-- Moon icon (shown in light mode) -->
          <svg v-else xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </button>

        <!-- Locale toggle -->
        <button
          @click="toggleLocale()"
          class="text-sm font-medium text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition"
        >
          {{ locale === 'en' ? 'RU' : 'EN' }}
        </button>

        <!-- Settings gear -->
        <button
          @click="showSettings = true"
          class="text-gray-400 dark:text-gray-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition"
          :title="$t('nav.settings')"
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
          class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 transition"
        >
          {{ $t('nav.logout') }}
        </button>
      </div>
    </nav>

    <main class="max-w-2xl mx-auto py-8 sm:py-12 px-4 space-y-8">
      <MicButton />
      <ResponseCard v-if="voice.result" :result="voice.result" />
      <p v-if="voice.error" class="text-center text-red-500 dark:text-red-400 text-sm">{{ voice.error }}</p>
    </main>

    <SettingsModal v-if="showSettings" @close="showSettings = false" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useVoiceStore } from '../stores/voice'
import { useTheme } from '../composables/useTheme'
import { useLocale } from '../composables/useLocale'
import MicButton from '../components/MicButton.vue'
import ResponseCard from '../components/ResponseCard.vue'
import SettingsModal from '../components/SettingsModal.vue'

const router      = useRouter()
const auth        = useAuthStore()
const voice       = useVoiceStore()
const { isDark, toggleTheme } = useTheme()
const { locale, toggleLocale } = useLocale()
const showSettings = ref(false)

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
