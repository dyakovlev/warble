import babel from 'rollup-plugin-babel'
import commonjs from 'rollup-plugin-commonjs'
import resolve from 'rollup-plugin-node-resolve'
import replace from 'rollup-plugin-replace'

export default {
  dest: 'build/app.js',
  entry: 'client/js/warble.js',
  format: 'iife',
  plugins: [
    babel({
      babelrc: false,
      exclude: 'node_modules/**',
      presets: [ [ 'es2015', { modules: false } ], 'react' ],
      plugins: [ 'external-helpers' ]
    }),
    commonjs({
      exclude: 'node_modules/process-es6/**',
      include: [
        'node_modules/fbjs/**',
        'node_modules/object-assign/**',
        'node_modules/react/**',
        'node_modules/react-dom/**',
        'node_modules/react-redux/**',
        'node_modules/hoist-non-react-statics/**',
        'node_modules/invariant/**',
      ],
      namedExports: {
          'node_modules/react/react.js': ['PropTypes', 'createElement', 'Children', 'Component']
      }
    }),
    replace({ 'process.env.NODE_ENV': JSON.stringify('development') }),
    resolve({ jsnext: true, browser: true, main: true })
  ]
}
