<template>
  <div class="flex items-center gap-3">
    <button
      @click="toggle"
      class="w-10 h-10 rounded-full flex items-center justify-center transition text-indigo-600 dark:text-indigo-400 bg-indigo-100 dark:bg-indigo-900/60 hover:bg-indigo-200 dark:hover:bg-indigo-800/60"
      :title="playing ? $t('audio.button.pause') : $t('audio.button.play')"
    >
      <!-- Pause icon -->
      <svg v-if="playing" class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
        <rect x="5" y="4" width="4" height="16" rx="1"/>
        <rect x="15" y="4" width="4" height="16" rx="1"/>
      </svg>
      <!-- Play icon -->
      <svg v-else class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
        <polygon points="5,3 19,12 5,21"/>
      </svg>
    </button>
<audio ref="audioEl" :src="src" @ended="playing = false" @canplaythrough="onCanPlay" />
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  src:      { type: String,  required: true },
  autoplay: { type: Boolean, default: false },
})

const audioEl = ref(null)
const playing  = ref(false)
let   autoPlayed = false

function onCanPlay() {
  if (props.autoplay && !autoPlayed) {
    autoPlayed = true
    audioEl.value.play()
    playing.value = true
  }
}

function toggle() {
  if (playing.value) {
    audioEl.value.pause()
    playing.value = false
  } else {
    audioEl.value.play()
    playing.value = true
  }
}
</script>
