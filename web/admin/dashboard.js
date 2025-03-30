import { w2ui, w2layout, w2sidebar, w2toolbar, w2utils } from './lib/w2ui.es6.min.js'

window.w2ui = w2ui

w2utils.settings.dataType = 'JSON'

w2utils.formatters['safe'] = (_, extra) => w2utils.encodeTags(extra.value)

window.dashboardSidebarSearch = function(value) {
  // Normalize the string to ensure consistent comparison
  const normalizeString = str => str.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
  dashboardSidebar.expandAll()
  dashboardSidebar.search(value, (str, node) => {
    const str1 = normalizeString(str.toLowerCase())
    const str2 = normalizeString(node.text.toLowerCase())
    return str2.indexOf(str1) != -1
  })
}

const dashboardSidebar = new w2sidebar({
  name: 'dashboardSidebar',
  topHTML: `
    <div style="margin: 3px 5px;">
      <div style="margin: 10px 5px;">
        <div style="margin-bottom: 10px; font-size: 16px;">
          Amazing Core Dashboard
        </div>
        <div id="sidebar-host" style="font-size: 12px;">
          <script>document.getElementById("sidebar-host").innerText = '@' + window.location.host</script>
        </div>
      </div>
      <div>
        <input id="dashboard-sidebar-search" style="width: 100%" class="w2ui-input" placeholder="Jump to..." onkeyup="dashboardSidebarSearch(this.value)">
      </div>
    </div>
  `,
  bottomHTML: '<div id="logout-toolbar"></div>',
  nodes: [
    {
      id: 'general', text: 'General', group: true, expanded: true, nodes: [
        {
          id: 'blob-db', text: 'blob.db', icon: 'fa fa-database', selected: true,
          onClick: async function() {
            const module = await import('./widgets/blob_grid.js')
            setDashboardWidget(module.createBlobGrid())
          },
        },
        {
          id: 'random-names', text: 'Random Names', icon: 'fa fa-dice',
          onClick: async function() {
            const module = await import('./widgets/random_names_grid.js')
            setDashboardWidget(module.createRandomNamesGrid())
          },
        },
      ]
    }
  ],
  onRender: async function(event) {
    await event.complete
    new w2toolbar({
      name: 'logoutToolbar',
      box: '#logout-toolbar',
      items: [
        {
          type: 'button', id: 'logout', text: 'Log out', icon: 'fa fa-right-from-bracket', onClick: async () => {
            const res = await fetch('/logout', {
              method: 'POST',
              headers: {
                'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content,
                'Content-Type': 'application/json'
              },
            })

            if (res.status != 200) {
              res.json()
                .then(data => w2utils.notify(`${data?.message}, ${res.status}: ${res.statusText}`, { timeout: 4000, error: true }))
                .catch(() => w2utils.notify(`Failed to logout, ${res.status}: ${res.statusText}`, { timeout: 4000, error: true }))
              return
            }

            window.location = '/'
          }
        },
        { type: 'button', id: 'github', icon: 'fa-brands fa-github', onClick: () => window.open('https://github.com/dv1x3r/amazing-core', '_blank') },
        { type: 'button', id: 'discord', icon: 'fa-brands fa-discord', onClick: () => window.open('https://discord.gg/TWfTBbfdA9', '_blank') },
      ],
    })
  },
})

const dashboardLayout = new w2layout({
  name: 'dashboardLayout',
  box: '#dashboard-layout',
  panels: [
    { type: 'left', size: 240, html: dashboardSidebar },
    { type: 'main' },
  ],
})

function setDashboardWidget(widget) {
  if (dashboardLayout.get('main').html.destroy) {
    dashboardLayout.get('main').html.destroy()
  }
  dashboardLayout.html('main', widget)
}

import('./widgets/blob_grid.js').then(module => {
  setDashboardWidget(module.createBlobGrid())
})

