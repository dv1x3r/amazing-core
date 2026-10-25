import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createQuestLayout() {
  const npcListGrid = new w2grid({
    name: 'questNPCListGrid',
    url: {
      get: '/api/v1/npc/grid',
    },
    recid: 'id',
    recordHeight: 28,
    multiSelect: false,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: false,
      toolbarEdit: false,
      toolbarDelete: false,
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
        field: 'quests',
        text: '# Quests',
        size: '70px',
        render: 'int',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'name', direction: 'asc' },
    ],
    onSelect: function(event) { reloadQuestGrid(event) },
    onDelete: function(event) { event.preventDefault() },
  })

  const questGrid = new w2grid({
    name: 'questGrid',
    header: '<i class="fa fa-scroll"></i> Quests',
    url: {
      get: '/api/v1/npc/:id/quest/grid',
      remove: '/api/v1/npc/quest/remove',
    },
    routeData: { id: 0 },
    recid: 'id',
    recordHeight: 28,
    reorderRows: true,
    show: {
      header: true,
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: false,
      toolbarReload: true,
      searchSave: false,
      lineNumbers: true,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        hidden: true,
      },
      {
        field: 'name',
        text: 'Quest Name',
        size: '200px',
        render: 'text',
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '200px',
        render: 'dropdown-tooltip',
      },
      {
        field: 'xp',
        text: 'XP',
        size: '60px',
        render: 'int',
      },
      {
        field: 'coins',
        text: 'Coins',
        size: '60px',
        render: 'int',
      },
      {
        field: 'npc_points',
        text: 'NPC Points',
        size: '80px',
        render: 'int',
      },
      {
        field: 'cooldown',
        text: 'Cooldown',
        size: '80px',
        render: 'int',
      },
    ],
    onAdd: function(event) { openQuestPopup(event) },
    onEdit: function(event) { openQuestPopup(event) },
    onDblClick: function(event) { openQuestPopup(event) },
    onReorderRow: function(event) { helpers.w2reorder(event, { url: '/api/v1/npc/quest/reorder' }) },
  })

  const subLayout = new w2layout({
    name: 'questSubLayout',
    padding: 5,
    panels: [
      { type: 'main', size: '60%', html: questGrid },
      { type: 'bottom', size: '40%', resizable: true },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('bottom', '/admin/pages/quests.html')
    },
    onDestroy: function() {
      questGrid.destroy()
    },
  })

  function clearQuestGrid() {
    questGrid.routeData.id = 0
    questGrid.clear()
    questGrid.toolbar.disable('w2ui-add')
  }

  async function reloadQuestGrid(event) {
    await event.complete
    const selection = event.owner.getSelection()
    if (selection.length != 1) {
      clearQuestGrid()
    } else {
      questGrid.routeData.id = selection[0]
      await questGrid.reload()
      questGrid.toolbar.enable('w2ui-add')
    }
  }

  return new w2layout({
    name: 'questLayout',
    panels: [
      { type: 'left', size: '350px', resizable: true, html: npcListGrid },
      { type: 'main', size: '100%', html: subLayout },
    ],
    onRender: function() { clearQuestGrid() },
    onDestroy: function() {
      npcListGrid.destroy()
      subLayout.destroy()
    },
  })
}

function openQuestPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'questForm',
    url: '/api/v1/npc/:id/quest/form',
    routeData: { id: event.owner.routeData.id },
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
          label: 'Quest Name',
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
        field: 'xp',
        type: 'int',
        html: {
          label: 'XP',
          attr: 'size="10" placeholder="0"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'coins',
        type: 'int',
        html: {
          label: 'Coins',
          attr: 'size="10" placeholder="0"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'npc_points',
        type: 'int',
        html: {
          label: 'NPC Points',
          attr: 'size="10" placeholder="0"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'cooldown',
        type: 'int',
        html: {
          label: 'Cooldown',
          attr: 'size="10" placeholder="0"',
          text: ' seconds',
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
    title: isEditMode ? 'Edit Quest' : 'New Quest',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 400, showMax: false, resizable: false,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

