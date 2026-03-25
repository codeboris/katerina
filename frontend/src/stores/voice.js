import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useVoiceStore = defineStore('voice', () => {
  const isRecording  = ref(false)
  const isProcessing = ref(false)
  const result       = ref(null)
  const error        = ref(null)

  let mediaRecorder = null
  let chunks        = []

  async function startRecording() {
    error.value  = null
    result.value = null
    chunks       = []

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder = new MediaRecorder(stream)
    mediaRecorder.ondataavailable = (e) => { if (e.data.size > 0) chunks.push(e.data) }
    mediaRecorder.start()
    isRecording.value = true
  }

  async function stopRecording() {
    return new Promise((resolve) => {
      mediaRecorder.onstop = resolve
      mediaRecorder.stop()
      mediaRecorder.stream.getTracks().forEach((t) => t.stop())
      isRecording.value = false
    })
  }

  async function processAudio() {
    if (!chunks.length) return
    isProcessing.value = true
    error.value        = null

    try {
      const blob = new Blob(chunks, { type: 'audio/webm' })
      const form = new FormData()
      form.append('audio', blob, 'recording.webm')
      const { data } = await api.post('/voice/process', form)
      result.value = data
    } catch (e) {
      error.value = e.response?.data?.error || 'Processing failed'
    } finally {
      isProcessing.value = false
    }
  }

  return { isRecording, isProcessing, result, error, startRecording, stopRecording, processAudio }
})
