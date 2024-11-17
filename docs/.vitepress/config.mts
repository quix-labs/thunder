import {type DefaultTheme, defineConfig} from 'vitepress'

export default defineConfig({
    title: "Thunder (Quix Labs)",
    lang: 'en-US',
    description: "Fast, efficient, and optimized synchronization between SQL databases and indexers.",

    lastUpdated: false,
    cleanUrls: true,

    srcExclude: [
        'README.md'
    ],

    head: [
        ['link', {rel: 'icon', type: 'image/svg+xml', href: '/logo.svg'}],
        ['link', {rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png'}],
        ['link', {rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png'}],
        ['link', {rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png'}],
        ['meta', {name: 'theme-color', content: '#5f67ee'}],
        ['meta', {property: 'og:type', content: 'website'}],
        ['meta', {property: 'og:locale', content: 'en'}],
        ['meta', {property: 'og:title', content: 'Thunder | Sync your database with Elasticsearch'}],
        ['meta', {property: 'twitter:title', content: 'Thunder | Sync your database with Elasticsearch'}],
        ['meta', {property: 'og:site_name', content: 'Thunder'}],
        ['meta', {property: 'twitter:card', content: 'summary_large_image'}],
        ['meta', {property: 'twitter:image:src', content: 'https://thunder.quix-labs.com/thunder-og.png'}],
        ['meta', {property: 'og:image', content: 'https://thunder.quix-labs.com/thunder-og.png'}],
        ['meta', {property: 'og:image:type', content: 'image/png'}],
        ['meta', {property: 'og:image:width', content: '1280'}],
        ['meta', {property: 'og:image:height', content: '640'}],
        ['meta', {property: 'og:url', content: 'https://thunder.quix-labs.com'}],
    ],

    sitemap: {
        hostname: 'https://thunder.quix-labs.com'
    },

    themeConfig: {
        outline: [2, 3],
        logo: '/logo.svg',
        siteTitle: "Thunder",
        nav: [
            {text: 'Guide', link: '/guide/what-is-thunder', activeMatch: '/guide/'},
            {text: 'Team', link: '/team', activeMatch: '/team/'},
        ],

        socialLinks: [
            {icon: 'github', link: 'https://github.com/quix-labs/thunder'}
        ],

        sidebar: {
            '/guide/': {base: '/guide/', items: sidebarGuide()},
        },

        editLink: {
            pattern: 'https://github.com/quix-labs/thunder/edit/main/docs/:path',
            text: 'Edit this page on GitHub'
        },

        search: {
            provider: 'local',
        },

        footer: {
            message: 'Released under the <a href="https://github.com/quix-labs/thunder/blob/main/LICENSE.md">MIT License</a>.',
            copyright: `Copyright © ${new Date().getFullYear()} - <a href="https://www.quix-labs.com">Quix Labs</a>`
        }
    }
})


function sidebarGuide(): DefaultTheme.SidebarItem[] {
    return [
        {
            text: 'Introduction',
            collapsed: false,
            items: [
                {text: 'What is Thunder?', link: 'what-is-thunder'},
                {text: 'Installation', link: 'installation'},
            ]
        },

        {
            text: 'Core Concepts',
            collapsed: false,
            items: [
                {text: 'Sources', link: 'sources'},
                {text: 'Targets', link: 'targets'},
                {text: 'Processors', link: 'processors'},
                {text: 'Exporters', link: 'exporters'},
            ]
        },
        {
            text: 'Modules',
            collapsed: false,
            base: '/guide/modules/',
            items: [
                {text: 'HTTP Server', link: 'http-server'},
                {
                    text: 'API',
                    collapsed: false,
                    base: '/guide/modules/api/',
                    link: '',
                    items: [
                        {text: 'Sources', link: 'sources'},
                        {text: 'Targets', link: 'targets'},
                        {text: 'Processors', link: 'processors'},
                    ]
                },
                {text: 'Frontend', link: 'frontend'},
            ]
        },
        {
            text: 'Source Drivers',
            collapsed: false,
            base: '/guide/drivers/',
            items: [
                {text: 'PostgreSQL Flash', link: 'postgresql_flash'},
                {text: 'MySQL/MariaDB', link: 'mysql'},
            ]
        },

    ]
}