import { createI18n } from 'vue-i18n'
import chinese from './zh'
import spanish from './es'
import english from './en'

const browserLang = navigator.language.toLowerCase()
let locale = 'en'
if (browserLang.startsWith('zh')) {
  locale = 'zh'
} else if (browserLang.startsWith('es')) {
  locale = 'es'
}

const i18n = createI18n({
  legacy: false,
  locale,
  fallbackLocale: 'en',
  messages: {
    zh : chinese,
    es : spanish,
    en : english,
  }
})

export default i18n