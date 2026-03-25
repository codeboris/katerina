<template>
  <div class="bg-white rounded-2xl shadow-md p-6 space-y-5 animate-fade-in">
    <div>
      <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-1">You said</p>
      <p class="text-gray-700">{{ result.original }}</p>
    </div>

    <div class="border-t pt-4">
      <p class="text-xs font-semibold text-green-600 uppercase tracking-wide mb-1">Corrected English</p>
      <p class="text-gray-800 font-medium">{{ result.corrected }}</p>
    </div>

    <div class="border-t pt-4">
      <p class="text-xs font-semibold text-blue-500 uppercase tracking-wide mb-1">Russian Translation</p>
      <p class="text-gray-700">{{ result.translation }}</p>
    </div>

    <div class="border-t pt-4">
      <p class="text-xs font-semibold text-indigo-600 uppercase tracking-wide mb-1">AI Teacher Response</p>
      <p class="text-gray-800">{{ result.answer }}</p>
    </div>

    <div v-if="result.audio_url" class="border-t pt-4">
      <AudioPlayer :src="result.audio_url" />
    </div>
  </div>
</template>

<script setup>
import AudioPlayer from './AudioPlayer.vue'
defineProps({ result: { type: Object, required: true } })
</script>
