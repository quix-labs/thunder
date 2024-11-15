import {type DefaultTheme, defineConfig} from 'vitepress'

export default defineConfig({
    title: "Thunder (Quix Labs)",
    lang: 'en-US',
    description: "Fast, efficient, and optimized synchronization between SQL databases and indexers.",

    lastUpdated: false,
    cleanUrls: true,

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
            copyright: `Copyright © ${new Date().getFullYear()} - <a href="https://github.com/quix-labs">Quix Labs</a>`
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
    ]
}