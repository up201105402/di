// formkit.config.js
import { defaultConfig } from '@formkit/vue'
import { createMultiStepPlugin } from '@formkit/addons'
import '@formkit/addons/css/multistep'

const dirPath = (node) => {
  return node.value?.match('^(.+)\/([^\/]+)$') != null;
}

// override default rule behaviors for your custom rule
dirPath.blocking = false
dirPath.skipEmpty = false
dirPath.debounce = 20 // milliseconds
dirPath.force = true

const floats = (node) => {
  return node.value?.match('^(\\s*-?\\d+(\\.\\d+)?)(\\s*,\\s*-?\\d+(\\.\\d+)?)*$') != null;
}

// override default rule behaviors for your custom rule
floats.blocking = false
floats.skipEmpty = false
floats.debounce = 20 // milliseconds
floats.force = true

const dict = (node) => {
  return node.value?.match('(?<key>[^:]+):(?<value>[^,]+)') != null;
}

// override default rule behaviors for your custom rule
dict.blocking = false
dict.skipEmpty = false
dict.debounce = 20 // milliseconds
dict.force = true

const config = defaultConfig({
  plugins: [createMultiStepPlugin()],
  rules: { dirPath, floats, dict }
})

export default config