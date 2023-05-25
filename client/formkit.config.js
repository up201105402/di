// formkit.config.js
import { defaultConfig } from '@formkit/vue'
import { createMultiStepPlugin } from '@formkit/addons'
import '@formkit/addons/css/multistep'

const isDirectoryPath = (node) => {
  return node.value?.match('^(.+)\/([^\/]+)$') != null;
}

// override default rule behaviors for your custom rule
isDirectoryPath.blocking = false
isDirectoryPath.skipEmpty = false
isDirectoryPath.debounce = 20 // milliseconds
isDirectoryPath.force = true

const config = defaultConfig({
  plugins: [createMultiStepPlugin()],
  rules: { isDirectoryPath }
})

export default config