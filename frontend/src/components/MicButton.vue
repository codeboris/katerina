<template>
  <div class="flex flex-col items-center gap-4">
    <button
      @click="toggle"
      :disabled="voice.isProcessing"
      class="w-24 h-24 rounded-full flex items-center justify-center text-4xl shadow-lg transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-indigo-300"
      :class="voice.isRecording
        ? 'bg-red-500 hover:bg-red-600 animate-pulse text-white'
        : 'bg-indigo-600 hover:bg-indigo-700 text-white disabled:opacity-50'"
    >
      <span v-if="voice.isRecording">⏹</span>
      <span v-else-if="voice.isProcessing" class="text-2xl">⏳</span>
      <span v-else>🎤</span>
    </button>

    <p class="text-sm text-gray-600 font-medium">
      <span v-if="voice.isRecording">Recording… click to stop</span>
      <span v-else-if="voice.isProcessing">Processing with AI…</span>
      <span v-else>Click to speak in English</span>
    </p>
  </div>
</template>

<script setup>
import { useVoiceStore } from '../stores/voice'

const voice = useVoiceStore()

async function toggle() {
  if (voice.isRecording) {
    await voice.stopRecording()
    await voice.processAudio()
  } else {
    await voice.startRecording()
  }
}
</script>
