import { ref } from 'vue'

const STORAGE_KEY = 'theme'
const isDark = ref(false)

function applyTheme(dark) {
  isDark.value = dark
  document.documentElement.classList.toggle('dark', dark)
  localStorage.setItem(STORAGE_KEY, dark ? 'dark' : 'light')
}

function initTheme() {
  const saved = localStorage.getItem(STORAGE_KEY)
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  applyTheme(saved === 'dark' || (saved === null && prefersDark))
}

function toggleTheme() {
  applyTheme(!isDark.value)
}

export function useTheme() {
  return { isDark, initTheme, toggleTheme }
}
