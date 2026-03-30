import { w2layout, w2sidebar, w2toolbar } from '/lib/w2ui.es6.min.js'
import { w2init, registerSidebarSearch } from '/lib/w2ui.helpers.js'

w2init()

const dashboardSidebar = new w2sidebar({
  name: 'dashboardSidebar',
  topHTML: `
    <div style="margin: 3px 5px;">
      <div style="margin: 10px 5px;">
        <div style="margin-bottom: 10px; font-size: 16px;">
          Amazing Core Dashboard
        </div>
        <div id="sidebar-host" style="font-size: 12px;">
          Loading...
        </div>
      </div>
      <div>
        <input id="dashboard-sidebar-search" class="w2ui-input" style="width:100%;" placeholder="Jump to...">
      </div>
    </div>
  `,
  bottomHTML: '<div id="logout-toolbar"></div>',
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
      id: 'general',
      text: 'General',
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
          id: 'random-names',
          text: 'Random Names',
          icon: 'fa fa-dice',
          onClick: async function() {
            const module = await import('./widgets/random_names.js')
            setDashboardWidget(module.createRandomNameGrid)
          },
        },
        {
          id: 'dummy-config',
          text: 'Dummy Config',
          icon: 'fa fa-person-digging',
          onClick: async function() {
            const module = await import('./widgets/dummy_config.js')
            setDashboardWidget(module.createDummyForm)
          },
        },
      ]
    },
    {
      id: 'cdn',
      text: 'Databases',
      group: true,
      expanded: true,
      nodes: [
        {
          id: 'core-db',
          text: 'SQL Explorer: core.db',
          icon: 'fa fa-database',
          onClick: async function() {
            const module = await import('/lib/w2ui.widgets.js')
            setDashboardWidget(() => module.createSqlExplorerLayout({ url: '/api/v1/sql' }))
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
        {
          id: 'blob-db',
          text: 'Cache Files: blob.db',
          icon: 'fa fa-database',
          onClick: async function() {
            const module = await import('./widgets/blob_db.js')
            setDashboardWidget(module.createBlobGrid)
          },
        },
      ]
    }
  ],
  onRender: async function(event) {
    await event.complete
    document.getElementById('sidebar-host').innerText = '@' + window.location.host
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
      ],
    })
    const search = registerSidebarSearch(dashboardSidebar)
    const el = document.getElementById('dashboard-sidebar-search')
    el.addEventListener('keyup', e => search(e.target.value))
  },
})

const dashboardLayout = new w2layout({
  name: 'dashboardLayout',
  box: '#dashboard-layout',
  panels: [
    { type: 'left', size: 240, html: dashboardSidebar },
    { type: 'main', style: 'border-left: 1px solid #e0e0e0;', html: '' },
  ],
})

function setDashboardWidget(createWidget) {
  if (dashboardLayout.get('main').html.destroy) {
    dashboardLayout.get('main').html.destroy()
  }
  dashboardLayout.html('main', createWidget())
}

dashboardLayout.load('main', '/admin/pages/welcome.html')

