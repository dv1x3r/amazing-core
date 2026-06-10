import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createSiteFrameLayout() {
  const grid = new w2grid({
    name: 'siteFrameGrid',
    url: {
      get: '/api/v1/siteframe/grid',
      remove: '/api/v1/siteframe/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
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
        field: 'container',
        text: 'Asset Container',
        size: '200px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'type_value',
        text: 'Type Value',
        size: '120px',
        render: 'text',
        sortable: true,
        searchable: 'int',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'asc' },
    ],
    onAdd: function(event) { openSiteFramePopup(event) },
    onEdit: function(event) { openSiteFramePopup(event) },
    onDblClick: function(event) { openSiteFramePopup(event) },
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
    },
  })
}

function openSiteFramePopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'siteFrameForm',
    url: '/api/v1/siteframe/form',
    record: record,
    fields: [
      {
        field: 'id',
        type: 'text',
        html: {
          label: 'ID',
          attr: 'size="15" readonly',
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
      {
        field: 'type_value',
        type: 'int',
        required: true,
        html: {
          label: 'Type Value',
          attr: 'size="15"',
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
    title: isEditMode ? 'Edit Site Frame' : 'New Site Frame',
    body: '<div id="site-frame-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 270, showMax: false, resizable: false,
  })
    .then(() => form.render('#site-frame-form'))
    .close(() => form.destroy())
}

