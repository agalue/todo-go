/* eslint-env node */
module.exports = {
  root: true,
  extends: [
    'plugin:vue/vue3-recommended',
  ],
  parserOptions: {
    parser: '@babel/eslint-parser',
    requireConfigFile: false,
    ecmaVersion: 'latest',
    sourceType: 'module',
  }
}