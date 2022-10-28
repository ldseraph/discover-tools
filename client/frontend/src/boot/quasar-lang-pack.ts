import { boot } from 'quasar/wrappers'
import { Quasar } from 'quasar'
const langList = import.meta.glob('../../node_modules/quasar/lang/*.mjs')

function loadLang(lang: string) {
  langList[`../../node_modules/quasar/lang/${lang}.mjs`]().then(lang => {
    Quasar.lang.set(lang.default)
  })
}

export default boot(async () => {
  try {
    loadLang('zh-CN')
  }
  catch (err) {
    console.error(err)
  }
})

