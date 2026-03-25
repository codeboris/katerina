<template>
  <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-md dark:shadow-gray-900/50 p-4 sm:p-6 space-y-5 animate-fade-in">
    <div>
      <p class="text-xs font-semibold text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-1">You said</p>
      <p class="text-gray-700 dark:text-gray-300">{{ result.original }}</p>
    </div>

    <div class="border-t dark:border-gray-700 pt-4">
      <p class="text-xs font-semibold text-green-600 dark:text-green-400 uppercase tracking-wide mb-1">Corrected English</p>
      <p class="text-gray-800 dark:text-gray-100 font-medium">{{ result.corrected }}</p>
      <p v-if="result.translation" class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ result.translation }}</p>
    </div>

    <div class="border-t dark:border-gray-700 pt-4">
      <p class="text-xs font-semibold text-indigo-600 dark:text-indigo-400 uppercase tracking-wide mb-1">AI Teacher Response</p>
      <p class="text-gray-800 dark:text-gray-100">{{ result.answer }}</p>
      <p v-if="result.answer_translation" class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ result.answer_translation }}</p>
    </div>

    <div v-if="result.audio_url" class="border-t dark:border-gray-700 pt-4">
      <AudioPlayer :src="result.audio_url" :autoplay="true" />
    </div>
  </div>
</template>

<script setup>
import AudioPlayer from './AudioPlayer.vue'
defineProps({ result: { type: Object, required: true } })
</script>
