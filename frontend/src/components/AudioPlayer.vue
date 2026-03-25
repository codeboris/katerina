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
    <audio ref="audioEl" :src="src" @ended="playing = false" />
  </div>
</template>

<script setup>
import { ref } from 'vue'

defineProps({ src: { type: String, required: true } })

const audioEl = ref(null)
const playing  = ref(false)

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
