// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: true},
    vue: {
        propsDestructure: true
    },
    nitro: {
        output: {
            publicDir: './build',
        },
        routeRules: {
            '/go-api/**': {proxy: 'http://localhost:3000/go-api/**'}
        }
    },
    modules: ['@nuxt/ui'],
})