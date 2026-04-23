// @ts-check
// Docusaurus config. Docs live in ../docs (shared source-of-truth with the
// repo-level recipes) so nothing is duplicated between the site and agent
// consumption.

import { themes as prismThemes } from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'openrouter-go',
  tagline: 'Zero-dependency Go client for the OpenRouter API',

  url: 'https://hra42.github.io',
  baseUrl: '/openrouter-go/',

  organizationName: 'hra42',
  projectName: 'openrouter-go',
  trailingSlash: false,

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  i18n: { defaultLocale: 'en', locales: ['en'] },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          path: '../docs',
          routeBasePath: 'docs',
          sidebarPath: './sidebars.js',
          editUrl:
            'https://github.com/hra42/openrouter-go/edit/main/docs/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'openrouter-go',
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'recipes',
            position: 'left',
            label: 'Recipes',
          },
          {
            href: 'https://pkg.go.dev/github.com/hra42/openrouter-go',
            label: 'Godoc',
            position: 'right',
          },
          {
            href: 'https://github.com/hra42/openrouter-go',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              { label: 'Recipes', to: '/docs/recipes/' },
              { label: 'Godoc', href: 'https://pkg.go.dev/github.com/hra42/openrouter-go' },
            ],
          },
          {
            title: 'Community',
            items: [
              { label: 'GitHub', href: 'https://github.com/hra42/openrouter-go' },
              { label: 'Issues', href: 'https://github.com/hra42/openrouter-go/issues' },
            ],
          },
          {
            title: 'OpenRouter',
            items: [
              { label: 'openrouter.ai', href: 'https://openrouter.ai' },
              { label: 'API docs', href: 'https://openrouter.ai/docs' },
            ],
          },
        ],
        copyright: `Built with Docusaurus. Library is Unlicensed — public domain.`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
        additionalLanguages: ['go', 'bash', 'json'],
      },
      colorMode: {
        defaultMode: 'dark',
        respectPrefersColorScheme: true,
      },
    }),
};

export default config;
