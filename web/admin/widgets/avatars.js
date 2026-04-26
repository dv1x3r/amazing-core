import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createWidget() {
  const grid = new w2grid({
    name: 'avatarGrid',
    url: {
      get: '/api/v1/avatar/grid',
      save: '/api/v1/avatar/grid',
      remove: '/api/v1/avatar/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: true,
      toolbarSearch: true,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        sortable: true,
        searchable: 'int',
      },
      {
        field: 'name',
        text: 'Avatar Name',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'max_outfits',
        text: 'Max Outfits',
        size: '120px',
        render: 'int',
        sortable: true,
        searchable: 'int',
        editable: { type: 'int' },
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '250px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
        editable: helpers.remoteListOptions('/api/v1/container'),
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function(event) { openAvatarPopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
  })

  return new w2layout({
    name: 'avatarLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/avatars.html')
    },
    onDestroy: function() {
      grid.destroy()
    },
  })
}

function openAvatarPopup(event) {
  const form = new w2form({
    name: 'avatarForm',
    url: '/api/v1/avatar/form',
    fields: [
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Avatar Name',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'max_outfits',
        type: 'int',
        required: true,
        html: {
          label: 'Max Outfits',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'container',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/container'),
        html: {
          label: 'Asset Container',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
    ],
    actions: {
      async Save() {
        const res = await this.save()
        if (res.status == 'success') {
          event.owner.reload()
          w2popup.close()
        }
      },
      Cancel() { w2popup.close() },
    },
  })

  w2popup.open({
    title: 'New Avatar',
    body: '<div id="avatar-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 300, showMax: false, resizable: false,
  })
    .then(() => form.render('#avatar-form'))
    .close(() => form.destroy())
}

