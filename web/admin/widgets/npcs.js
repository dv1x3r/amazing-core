import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createNPCLayout() {
  const grid = new w2grid({
    name: 'npcGrid',
    url: {
      get: '/api/v1/npc/grid',
      remove: '/api/v1/npc/remove',
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
        field: 'name',
        text: 'NPC Name',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '200px',
        render: 'dropdown-tooltip',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'zone',
        text: 'Zone',
        size: '200px',
        render: 'dropdown-tooltip',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'is_locked',
        text: 'Locked',
        size: '70px',
        render: 'toggle',
      },
      {
        field: 'has_win_3d',
        text: 'Win 3D',
        size: '60px',
        render: 'toggle',
      },
      {
        field: 'has_osx_3d',
        text: 'OSX 3D',
        size: '60px',
        render: 'toggle',
      },
      {
        field: 'has_location',
        text: 'Loc',
        size: '60px',
        render: 'toggle',
      },
      {
        field: 'has_text',
        text: 'Text',
        size: '60px',
        render: 'toggle',
      },
      {
        field: 'has_relationship',
        text: 'Rel',
        size: '60px',
        render: 'toggle',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'name', direction: 'asc' },
    ],
    onAdd: function(event) { openNPCPopup(event) },
    onEdit: function(event) { openNPCPopup(event) },
    onDblClick: function(event) { openNPCPopup(event) },
  })

  return new w2layout({
    name: 'npcLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/npcs.html')
    },
    onDestroy: function() {
      grid.destroy()
    },
  })
}

function openNPCPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'npcForm',
    url: '/api/v1/npc/form',
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
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'NPC Name',
          attr: 'style="width:100%;"',
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
        field: 'zone',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/zone/dropdown'),
        html: {
          label: 'Zone',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'is_locked',
        type: 'checkbox',
        html: {
          label: 'Is Locked',
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
    title: isEditMode ? 'Edit NPC' : 'New NPC',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 320, showMax: false, resizable: false,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

