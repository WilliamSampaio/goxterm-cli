/**
 * plugins/index.js
 *
 * Automatically included in `./src/main.js`
 */

// Plugins
import pinia from './pinia'
import vuetify from './vuetify'

export function registerPlugins(app) {
  app.use(vuetify)
  app.use(pinia)
}
