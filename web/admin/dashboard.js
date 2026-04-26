import { w2layout, w2sidebar, w2toolbar } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'
import * as widgets from '/lib/w2ui.widgets.js'

helpers.w2init()
helpers.w2initLocale()

const dashboardSidebar = new w2sidebar({
  name: 'dashboardSidebar',
  topHTML: `
    <div style="margin: 3px 5px;">
      <div style="margin: 10px 5px;">
        <div style="margin-bottom: 10px; font-size: 16px;">
          Amazing Core Dashboard
        </div>
        <div style="display: flex; justify-content: space-between; font-size: 12px;">
          <span id="sidebar-host">Loading...</span>
          <span id="sidebar-version"></span>
        </div>
      </div>
      <div>
        <input id="dashboard-sidebar-search" class="w2ui-input" style="width:100%;" placeholder="Jump to...">
      </div>
    </div>
  `,
  bottomHTML: '<div id="logout-toolbar" style="height: 40px;"></div>',
  nodes: [
    {
      id: 'welcome',
      text: 'Welcome',
      icon: 'fa fa-house',
      selected: true,
      onClick: async function() {
        if (dashboardLayout.get('main').html.destroy) {
          dashboardLayout.get('main').html.destroy()
        }
        dashboardLayout.load('main', '/admin/pages/welcome.html')
      },
    },
    {
      id: 'content',
      text: 'Content',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'assets',
          text: 'Assets',
          icon: 'fa-brands fa-unity',
          onClick: async function() {
            const module = await import('./widgets/assets.js')
            setDashboardWidget(module.createAssetGrid)
          },
        },
        {
          id: 'asset-containers',
          text: 'Containers',
          icon: 'fa fa-folder-open',
          onClick: async function() {
            const module = await import('./widgets/assets.js')
            setDashboardWidget(module.createContainerLayout)
          },
        },
      ],
    },
    {
      id: 'collections',
      text: 'Collections',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'avatars',
          text: 'Avatars',
          icon: 'fa fa-user-astronaut',
          onClick: async function() {
            const module = await import('./widgets/avatars.js')
            setDashboardWidget(module.createWidget)
          },
        },
        {
          id: 'site-frame',
          text: 'Site Frame',
          icon: 'fa fa-layer-group',
          onClick: async function() {
            const module = await import('./widgets/site_frame.js')
            setDashboardWidget(module.createWidget)
          },
        },
        {
          id: 'random-names',
          text: 'Random Names',
          icon: 'fa fa-dice',
          onClick: async function() {
            const module = await import('./widgets/random_names.js')
            setDashboardWidget(module.createWidget)
          },
        },
      ]
    },
    {
      id: 'general',
      text: 'Parameters',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'dummy-params',
          text: 'Dummy Parameters',
          icon: 'fa fa-person-digging',
          onClick: async function() {
            const module = await import('./widgets/dummy_params.js')
            setDashboardWidget(module.createWidget)
          },
        },
      ]
    },
    {
      id: 'core-db',
      text: 'core.db',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'sql-explorer',
          text: 'SQL Explorer',
          icon: 'fa fa-database',
          expanded: false,
          onClick: async function() {
            const module = await import('/lib/w2ui.widgets.js')
            setDashboardWidget(() => module.createSqlExplorerLayout({
              url: '/api/v1/sql',
            }))
          },
          nodes: await (async () => {
            const res = await fetch('/queries')
            const text = await res.text()
            const filenames = [...text.matchAll(/href="([^"]+)"/g)].map(m => m[1]).sort()
            return filenames.map(filename => ({
              id: `query-${filename}`,
              text: filename,
              icon: 'fa fa-file-code',
              onClick: async function() {
                const res = await fetch(`/queries/${filename}`)
                const initialQuery = await res.text()
                const module = await import('/lib/w2ui.widgets.js')
                setDashboardWidget(() => module.createSqlExplorerLayout({
                  url: '/api/v1/sql',
                  initialQuery,
                }))
              },
            }))
          })(),
        },
      ],
    },
    {
      id: 'blob-db',
      text: 'blob.db',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'asset-files',
          text: 'Asset Files',
          icon: 'fa fa-file',
          onClick: async function() {
            const module = await import('./widgets/blob_db.js')
            setDashboardWidget(module.createWidget)
          },
        },
      ],
    },
  ],
  onRender: async function(event) {
    await event.complete
    document.getElementById('sidebar-host').innerText = '@' + window.location.host
    document.getElementById('sidebar-version').innerText = document.getElementById('app-config').dataset.version
    new w2toolbar({
      name: 'logoutToolbar',
      box: '#logout-toolbar',
      items: [
        {
          type: 'button',
          id: 'logout',
          text: 'Log out',
          icon: 'fa fa-right-from-bracket',
          onClick: async () => {
            await fetch('/logout', { method: 'POST' })
            window.location = '/'
          }
        },
        {
          type: 'button',
          id: 'website',
          icon: 'fa fa-globe',
          tooltip: 'Website',
          onClick: () => window.open('https://amazingcore.org', '_blank')
        },
        {
          type: 'button',
          id: 'github',
          icon: 'fa-brands fa-github',
          tooltip: 'GitHub',
          onClick: () => window.open('https://github.com/dv1x3r/amazing-core', '_blank'),
        },
        {
          type: 'button',
          id: 'discord',
          icon: 'fa-brands fa-discord',
          tooltip: 'Discord',
          onClick: () => window.open('https://discord.gg/TWfTBbfdA9', '_blank'),
        },
        { type: 'spacer' },
        {
          type: 'button',
          id: 'locale',
          icon: 'fa fa-language',
          tooltip: 'Locale Settings',
          onClick: () => widgets.openLocalePopup(),
        },
        {
          id: 'dark',
          type: 'check',
          icon: 'fa fa-moon',
          tooltip: 'Dark Theme',
          checked: helpers.isDarkTheme(),
          onClick: async function(event) {
            await event.complete
            helpers.setDarkTheme(event.detail.item.checked)
          },
        },
      ],
    })
    const search = helpers.registerSidebarSearch(dashboardSidebar)
    const el = document.getElementById('dashboard-sidebar-search')
    el.addEventListener('keyup', e => search(e.target.value))
  },
})

const dashboardLayout = new w2layout({
  name: 'dashboardLayout',
  box: '#dashboard-layout',
  panels: [
    { type: 'left', size: 250, html: dashboardSidebar },
    { type: 'main', style: 'border-left: 1px solid #e0e0e0;' },
  ],
  onRender: async function(event) {
    await event.complete
    event.owner.load('main', '/admin/pages/welcome.html')
  },
})

function setDashboardWidget(createWidget) {
  if (dashboardLayout.get('main').html.destroy) {
    dashboardLayout.get('main').html.destroy()
  }
  dashboardLayout.html('main', createWidget())
}

