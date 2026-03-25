import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const STORAGE_KEY = 'locale'
const SUPPORTED = ['en', 'ru']

const locale = ref('en')

export function initLocale(i18n) {
  const saved = localStorage.getItem(STORAGE_KEY)
  const lang = SUPPORTED.includes(saved) ? saved : 'en'
  i18n.global.locale.value = lang
  locale.value = lang
}

export function useLocale() {
  const { locale: i18nLocale } = useI18n()

  function applyLocale(lang) {
    locale.value = lang
    i18nLocale.value = lang
    localStorage.setItem(STORAGE_KEY, lang)
  }

  function toggleLocale() {
    applyLocale(locale.value === 'en' ? 'ru' : 'en')
  }

  return { locale, applyLocale, toggleLocale }
}
