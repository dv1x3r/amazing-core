import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createStoreThemeLayout() {
  const grid = new w2grid({
    name: 'storeThemeGrid',
    header: '<i class="fa fa-store"></i> Store Themes',
    url: {
      get: '/api/v1/store/theme/grid',
      remove: '/api/v1/store/theme/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      header: true,
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
        hidden: true,
        sortable: true,
        searchable: 'int',
      },
      {
        field: 'icon',
        text: x => x.field ? '' : 'Icon',
        size: '40px',
        render: 'icon-sm',
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '250px',
        render: 'dropdown-tooltip',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'parent',
        text: 'Parent Theme',
        size: '200px',
        render: 'dropdown',
        sortable: true,
      },
      {
        field: 'ordinal',
        text: 'Ordinal',
        size: '70px',
        render: 'int',
        sortable: true,
      },
      {
        field: 'items',
        text: 'Items',
        size: '70px',
        render: 'int',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'parent', direction: 'asc' },
      { field: 'ordinal', direction: 'asc' },
    ],
    onAdd: function(event) { openStoreThemePopup(event) },
    onEdit: function(event) { openStoreThemePopup(event) },
    onDblClick: function(event) { openStoreThemePopup(event) },
  })

  return new w2layout({
    name: 'storeThemeLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/store_themes.html')
    },
    onDestroy: function() {
      grid.destroy()
    },
  })
}

function openStoreThemePopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'storeThemeForm',
    url: '/api/v1/store/theme/form',
    record: record ?? { ordinal: 1 },
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
        options: helpers.remoteListOptions('/api/v1/container/dropdown'),
        html: {
          label: 'Asset Container',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'parent',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/store/theme/dropdown'),
        html: {
          label: 'Parent Theme',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'ordinal',
        type: 'int',
        required: true,
        html: {
          label: 'Ordinal',
          attr: 'size="10"',
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
    title: isEditMode ? 'Edit Store Theme' : 'New Store Theme',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 300, showMax: false, resizable: false,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

