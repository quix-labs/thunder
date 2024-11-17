// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: true},
    vue: {
        propsDestructure: true,
    },
    vite: {
        optimizeDeps: {
            include: ["fast-deep-equal",] // FIX NUXT UI
        }
    },
    nitro: {
        output: {
            publicDir: process.env.NUXT_OUTPUT_DIR || '.output',
        },
        routeRules: {
            '/go-api/**': {proxy: 'http://localhost:3000/go-api/**'},
            '/processors/**': {ssr: false},
        }
    },
    modules: ['@nuxt/ui'],
})