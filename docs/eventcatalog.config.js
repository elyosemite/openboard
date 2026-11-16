/** @type {import('@eventcatalog/core/bin/eventcatalog.config').Config} */
export default {
  title: 'OpenBoard',
  tagline:
    'Architecture of OpenBoard, the open source, self-hosted project management platform: its decisions, domains and API.',
  organizationName: 'OpenBoard',
  homepageLink: 'https://github.com/elyosemite/openboard',
  editUrl: 'https://github.com/elyosemite/openboard/edit/main/docs',
  // EventCatalog sends usage telemetry by default, and it is not anonymous (ADR 0008).
  telemetry: false,
  // 3000 is the web app in docker compose.
  port: 3001,
  output: 'static',
  linkValidation: {
    // Broken internal links fail CI; locally they only warn.
    onBrokenLinks: process.env.CI ? 'error' : 'warn',
    onBrokenAnchors: process.env.CI ? 'error' : 'warn',
    // Absolute GitHub URLs (editUrl, homepageLink, EventCatalog's own license link)
    // are reported by their path as if they were internal pages.
    ignore: ['/elyosemite/openboard', '/elyosemite/openboard/**', '/event-catalog/eventcatalog/**'],
  },
  trailingSlash: false,
  base: '/',
  search: {
    type: 'resource',
  },
  navigation: {
    pages: ['list:all'],
  },
  logo: {
    alt: 'OpenBoard',
    text: 'OpenBoard',
  },
  llmsTxt: {
    enabled: true,
  },
  // Required by EventCatalog: generated once when the catalog was created.
  cId: 'fa469f2e-9ed6-4b45-a385-5c713a3caa61',
  tsd: 1794848689850,
};
