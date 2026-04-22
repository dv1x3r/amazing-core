import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createWidget() {
  const grid = new w2grid({
    name: 'siteFrameGrid',
    url: {
      get: '/api/v1/siteframe/grid',
      save: '/api/v1/siteframe/grid',
      remove: '/api/v1/siteframe/remove',
    },
    recid: 'id',
    recordHeight: 28,
    multiSearch: true,
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
        field: 'type_value',
        text: 'Type Value',
        size: '120px',
        render: 'text',
        sortable: true,
        searchable: 'int',
        editable: { type: 'int' },
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '400px',
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
    onAdd: function(event) { openSiteFramePopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
  })

  return new w2layout({
    name: 'siteFrameLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/site_frame.html')
    },
    onDestroy: function() {
      grid.destroy()
    }
  })
}

function openSiteFramePopup(event) {
  const form = new w2form({
    name: 'siteFrameForm',
    url: '/api/v1/siteframe/form',
    fields: [
      {
        field: 'type_value',
        type: 'int',
        required: true,
        html: {
          label: 'Type Value',
          attr: 'size="10"',
          span: 5,
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
          span: 5,
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
    title: 'New Site Frame',
    body: '<div id="site-frame-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#site-frame-form'))
    .close(() => form.destroy())
}

