<template>
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
    @click.self="$emit('close')"
  >
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-sm mx-4 p-6 space-y-5">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold text-gray-800">Settings</h2>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600 text-xl leading-none">&times;</button>
      </div>

      <div class="space-y-2">
        <p class="text-sm font-semibold text-gray-500 uppercase tracking-wide">AI Teacher Voice</p>
        <div v-if="settings.voices.length" class="space-y-2">
          <label
            v-for="v in settings.voices"
            :key="v.id"
            class="flex items-center gap-3 p-3 rounded-xl border cursor-pointer transition"
            :class="selected === v.id
              ? 'border-indigo-500 bg-indigo-50'
              : 'border-gray-200 hover:border-indigo-300'"
          >
            <input
              type="radio"
              :value="v.id"
              v-model="selected"
              class="accent-indigo-600"
            />
            <span class="text-sm text-gray-700">{{ v.label }}</span>
          </label>
        </div>
        <p v-else class="text-sm text-gray-400">Loading voices…</p>
      </div>

      <button
        @click="save"
        :disabled="settings.loading || selected === settings.current"
        class="w-full py-2 rounded-xl bg-indigo-600 text-white text-sm font-medium
               hover:bg-indigo-700 transition disabled:opacity-40 disabled:cursor-not-allowed"
      >
        {{ settings.loading ? 'Saving…' : 'Save' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '../stores/settings'

const emit = defineEmits(['close'])
const settings = useSettingsStore()
const selected = ref(settings.current)

onMounted(async () => {
  await settings.fetchVoices()
  try {
    await settings.fetchSettings()
  } catch {
    // fallback to first available voice
  }
  selected.value = settings.current ?? settings.voices[0]?.id ?? null
})

async function save() {
  await settings.saveVoice(selected.value)
  emit('close')
}
</script>
