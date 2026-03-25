import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useSettingsStore = defineStore('settings', () => {
  const voices  = ref([])
  const current = ref(null)
  const loading = ref(false)

  async function fetchVoices() {
    const { data } = await api.get('/voices')
    voices.value = data
  }

  async function fetchSettings() {
    const { data } = await api.get('/user/settings')
    current.value = data.voice
  }

  async function saveVoice(voiceId) {
    loading.value = true
    try {
      await api.put('/user/settings', { voice: voiceId })
      current.value = voiceId
    } finally {
      loading.value = false
    }
  }

  return { voices, current, loading, fetchVoices, fetchSettings, saveVoice }
})
