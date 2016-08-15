import babel from 'rollup-plugin-babel'

export default {
  entry: 'static/js/warble.js',
  format: 'iife',
  dest: 'build/app.js',
  plugins: [ babel() ]
}
