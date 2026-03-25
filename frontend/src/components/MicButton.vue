<template>
  <div class="flex flex-col items-center gap-4">
    <button
      @click="toggle"
      :disabled="voice.isProcessing"
      class="w-20 h-20 sm:w-24 sm:h-24 rounded-full flex items-center justify-center shadow-lg transition-all duration-200 focus:outline-none focus:ring-4"
      :class="voice.isRecording
        ? 'bg-red-500 hover:bg-red-600 focus:ring-red-400 dark:focus:ring-red-700 animate-pulse text-white'
        : 'bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-300 dark:focus:ring-indigo-700 text-white disabled:opacity-50'"
    >
      <!-- Stop square (recording) -->
      <svg v-if="voice.isRecording" class="w-7 h-7 sm:w-8 sm:h-8" viewBox="0 0 24 24" fill="currentColor">
        <rect x="5" y="5" width="14" height="14" rx="2"/>
      </svg>
      <!-- Spinner (processing) -->
      <svg v-else-if="voice.isProcessing" class="w-7 h-7 sm:w-8 sm:h-8 animate-spin" viewBox="0 0 24 24"
           fill="none" stroke="currentColor" stroke-width="2.5">
        <circle cx="12" cy="12" r="10" stroke-opacity="0.3"/>
        <path d="M12 2a10 10 0 0 1 10 10" stroke-linecap="round"/>
      </svg>
      <!-- Microphone (idle) -->
      <svg v-else class="w-7 h-7 sm:w-8 sm:h-8" viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 1a4 4 0 0 1 4 4v6a4 4 0 0 1-8 0V5a4 4 0 0 1 4-4z"/>
        <path d="M19 10a7 7 0 0 1-14 0H3a9 9 0 0 0 18 0h-2z"/>
        <line x1="12" y1="19" x2="12" y2="23" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        <line x1="9" y1="23" x2="15" y2="23" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </button>

    <p class="text-sm text-gray-600 dark:text-gray-300 font-medium">
      <span v-if="voice.isRecording">{{ $t('mic.status.recording') }}</span>
      <span v-else-if="voice.isProcessing">{{ $t('mic.status.processing') }}</span>
      <span v-else>{{ $t('mic.status.idle') }}</span>
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
