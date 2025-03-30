import { w2ui, w2form, w2toolbar, w2utils } from './lib/w2ui.es6.min.js'

window.w2ui = w2ui

w2utils.settings.dataType = 'JSON'

new w2form({
  name: 'loginForm',
  box: '#login-form',
  header: 'Amazing Core Dashboard',
  url: '/login',
  httpHeaders: { 'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content },
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
    { type: 'button', id: 'github', icon: 'fa-brands fa-github', onClick: () => window.open('https://github.com/dv1x3r/amazing-core', '_blank') },
    { type: 'button', id: 'discord', icon: 'fa-brands fa-discord', onClick: () => window.open('https://discord.gg/TWfTBbfdA9', '_blank') },
  ],
})
