module.exports = {
  root: true,

  env: {
    node: true
  },

  parserOptions: {
    ecmaVersion: 2021
  },

  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'vue/multi-word-component-names': 'off',
  },

  extends: [
    'plugin:vue/vue3-recommended',
    '@vue/typescript/recommended'
  ]
}
