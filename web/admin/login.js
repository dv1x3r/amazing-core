import { w2form, w2toolbar } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

helpers.w2init()

new w2form({
  name: 'loginForm',
  box: '#login-form',
  header: 'Amazing Core Dashboard',
  url: '/login',
  fields: [
    { field: 'username', type: 'text', html: { label: '', span: 0, attr: 'placeholder="Username" style="width: 100%"' } },
    { field: 'password', type: 'password', html: { label: '', span: 0, attr: 'placeholder="Password" style="width: 100%"' } },
  ],
  actions: {
    login: {
      text: 'Login',
      class: 'w2ui-btn w2ui-btn-blue',
      onClick: async function() {
        await this.save()
        window.location.href = '/'
      },
    },
  },
  onRender: async function(event) {
    await event.complete
    event.owner.fields.forEach(x => {
      x.$el.on('keydown', async e => {
        if (e.keyCode == 13) {
          await this.save()
          window.location.href = '/'
        }
      })
    })
  },
})

new w2toolbar({
  name: 'loginToolbar',
  box: '#login-toolbar',
  items: [
    {
      id: 'website',
      type: 'button',
      icon: 'fa fa-globe',
      tooltip: 'Website',
      onClick: () => window.open('https://amazingcore.org', '_blank'),
    },
    {
      id: 'github',
      type: 'button',
      icon: 'fa-brands fa-github',
      tooltip: 'GitHub',
      onClick: () => window.open('https://github.com/dv1x3r/amazing-core', '_blank'),
    },
    {
      id: 'discord',
      type: 'button',
      icon: 'fa-brands fa-discord',
      tooltip: 'Discord',
      onClick: () => window.open('https://discord.gg/TWfTBbfdA9', '_blank'),
    },
    { type: 'spacer' },
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
