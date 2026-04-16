import { w2grid, w2layout } from '/lib/w2ui.es6.min.js'

export function createWidget() {
  const grid = new w2grid({
    name: 'dummyGrid',
    url: {
      get: '/api/v1/dummy/grid',
      save: '/api/v1/dummy/grid',
    },
    recid: 'rowid',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: false,
      toolbarEdit: false,
      toolbarDelete: false,
      toolbarSave: true,
      toolbarSearch: true,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'key',
        text: 'Key',
        size: '250px',
        render: 'text',
        searchable: 'text',
      },
      {
        field: 'value',
        text: 'Value',
        size: '250px',
        render: 'text',
        searchable: 'text',
        editable: { type: 'text' },
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
  })

  return new w2layout({
    name: 'dummyLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/dummy_params.html')
    },
    onDestroy: function() {
      grid.destroy()
    }
  })
}

