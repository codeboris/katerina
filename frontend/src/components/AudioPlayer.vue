<template>
  <div class="flex items-center gap-3">
    <button
      @click="toggle"
      class="w-10 h-10 rounded-full bg-indigo-100 hover:bg-indigo-200 flex items-center justify-center transition"
      :title="playing ? 'Pause' : 'Play'"
    >
      <span>{{ playing ? '⏸' : '▶️' }}</span>
    </button>
    <span class="text-sm text-gray-500">{{ playing ? 'Playing response…' : 'Play AI response' }}</span>
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
