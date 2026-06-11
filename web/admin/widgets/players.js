import { w2form, w2grid, w2layout, w2popup, w2utils } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

let selectedPlayerID = null

export function createPlayerLayout() {
  selectedPlayerID = null
  const playerListGrid = new w2grid({
    name: 'playerListGrid',
    url: {
      get: '/api/v1/player/grid',
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
      toolbarReload: false,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        searchable: 'int',
      },
      {
        field: 'oid',
        text: 'Player OID',
        render: 'text',
        size: '135px',
        searchable: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'active_name',
        text: 'Active Avatar Name',
        render: 'text',
        size: '135px',
        searchable: 'text',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    onSelect: async function(event) {
      await event.complete
      const selection = event.owner.getSelection()
      selectedPlayerID = selection.length == 1 ? selection[0] : null
      renderActiveTab()
    },
    onDelete: function(event) { event.preventDefault() },
  })

  const subLayout = new w2layout({
    name: 'playerSubLayout',
    style: 'margin: 5px;',
    padding: 5,
    panels: [
      {
        type: 'main',
        html: createPlayerDetailsForm(),
        tabs: {
          active: 'details',
          tabs: [
            { id: 'details', text: 'Details' },
            { id: 'avatars', text: 'Avatars' },
            { id: 'outfits', text: 'Outfits' },
            { id: 'items', text: 'Items' },
            { id: 'information', text: 'Information' },
          ],
          onClick(event) {
            renderActiveTab(event.target)
          }
        }
      },
    ],
    onDestroy: function() {
      if (subLayout.get('main').html.destroy) {
        subLayout.get('main').html.destroy()
      }
    },
  })

  function renderActiveTab(tabID = subLayout.get('main').tabs.active) {
    if (subLayout.get('main').html.destroy) {
      subLayout.get('main').html.destroy()
    }
    switch (tabID) {
      case 'details':
        subLayout.html('main', createPlayerDetailsForm())
        break
      case 'avatars':
        subLayout.html('main', createPlayerAvatarsGrid())
        break
      case 'outfits':
        subLayout.html('main', createPlayerOutfitsGrid())
        break
      case 'items':
        subLayout.html('main', createPlayerItemsGrid())
        break
      case 'information':
        subLayout.load('main', '/admin/pages/players.html')
        break
    }
  }

  return new w2layout({
    name: 'playerLayout',
    panels: [
      { type: 'left', size: '340px', resizable: true, html: playerListGrid },
      { type: 'main', size: '100%', html: subLayout },
    ],
    onDestroy: function() {
      playerListGrid.destroy()
      subLayout.destroy()
    },
  })
}

function createPlayerDetailsForm() {
  return new w2form({
    name: `playerDetailsForm`,
    url: '/api/v1/player/form',
    recid: selectedPlayerID,
    focus: -1,
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
        field: 'created_at',
        type: 'datetime',
        html: {
          label: 'Created Date',
          attr: 'readonly',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'oid',
        type: 'text',
        disabled: selectedPlayerID == null,
        html: {
          label: 'Player OID',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'active_avatar',
        type: 'list',
        options: helpers.remoteListOptions(`/api/v1/player/${selectedPlayerID}/avatar`),
        disabled: selectedPlayerID == null,
        html: {
          label: 'Active Avatar',
          attr: 'placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'is_tutorial_completed',
        type: 'checkbox',
        disabled: selectedPlayerID == null,
        html: {
          label: 'Is Tutorial Completed',
          span: 2,
          column: 0,
        },
      },
      {
        field: 'is_qa',
        type: 'checkbox',
        disabled: selectedPlayerID == null,
        html: {
          label: 'Is QA',
          span: 2,
          column: 0,
        },
      },
      {
        field: 'max_outfits',
        type: 'int',
        html: {
          label: 'Max Outfits',
          attr: 'size="15"',
          span: 6,
          column: 0,
        },
      },
    ],
    toolbar: {
      items: [
        {
          id: 'save',
          type: 'button',
          text: 'Save',
          icon: 'fa fa-floppy-disk',
          disabled: selectedPlayerID == null,
          onClick: async event => {
            await event.owner.owner.save()
            w2utils.notify('Your changes have been saved.', { timeout: 5000 })
          },
        },
      ],
    },
  })
}

function createPlayerAvatarsGrid() {
  return new w2grid({
    name: 'playerAvatarGrid',
    url: {
      get: '/api/v1/player/:id/avatar/grid',
      remove: '/api/v1/player/avatar/remove',
    },
    routeData: { id: selectedPlayerID ?? 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: false,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
      },
      {
        field: 'oid',
        text: 'Player Avatar OID',
        size: '135px',
        render: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'name',
        text: 'Player Avatar Name',
        render: 'text',
        size: '200px',
      },
      {
        field: 'avatar',
        text: 'Avatar',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'outfit_no',
        text: 'Outfit No',
        size: '120px',
        render: 'int',
      },
      {
        field: 'is_active',
        text: 'Active',
        render: 'toggle',
        size: '60px',
      },
    ],
    onRender: function(event) {
      if (selectedPlayerID == null) {
        event.owner.toolbar.disable('w2ui-add')
      }
    },
    onAdd: function(event) { openPlayerAvatarPopup(event) },
    onEdit: function(event) { openPlayerAvatarPopup(event) },
    onDblClick: function(event) { openPlayerAvatarPopup(event) },
  })
}

function openPlayerAvatarPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'playerAvatarForm',
    url: '/api/v1/player/:id/avatar/form',
    routeData: { id: selectedPlayerID ?? 0 },
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
        field: 'oid',
        type: 'text',
        required: isEditMode,
        html: {
          label: 'Player Avatar OID',
          attr: isEditMode ? 'size="15"' : 'size="15" readonly',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Player Avatar Name',
          attr: 'style="width:100%;"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'avatar',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/avatar'),
        html: {
          label: 'Avatar',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'outfit_no',
        type: 'int',
        required: true,
        html: {
          label: 'Outfit No',
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
    title: isEditMode ? 'Edit Player Avatar' : 'New Player Avatar',
    body: '<div id="player-avatar-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 340, showMax: false, resizable: false,
  })
    .then(() => form.render('#player-avatar-form'))
    .close(() => form.destroy())
}

function createPlayerOutfitsGrid() {
  return new w2grid({
    name: 'playerOutfitGrid',
    url: {
      get: '/api/v1/player/:id/outfit/grid',
      remove: '/api/v1/player/outfit/remove',
    },
    routeData: { id: selectedPlayerID ?? 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: false,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
      },
      {
        field: 'oid',
        text: 'Player Outfit OID',
        size: '135px',
        render: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'player_avatar',
        text: 'Player Avatar',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'outfit_no',
        text: 'Outfit No',
        size: '120px',
        render: 'int',
      },
    ],
    onRender: function(event) {
      if (selectedPlayerID == null) {
        event.owner.toolbar.disable('w2ui-add')
      }
    },
    onAdd: function(event) { openPlayerOutfitPopup(event) },
    onEdit: function(event) { openPlayerOutfitPopup(event) },
    onDblClick: function(event) { openPlayerOutfitPopup(event) },
  })
}

function openPlayerOutfitPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'playerOutfitForm',
    url: '/api/v1/player/outfit/form',
    focus: -1,
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
        field: 'oid',
        type: 'text',
        required: isEditMode,
        html: {
          label: 'Player Outfit OID',
          attr: isEditMode ? 'size="15"' : 'size="15" readonly',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'player_avatar',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions(`/api/v1/player/${selectedPlayerID}/avatar`),
        html: {
          label: 'Player Avatar',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'outfit_no',
        type: 'int',
        required: true,
        html: {
          label: 'Outfit No',
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
    title: isEditMode ? 'Edit Player Outfit' : 'New Player Outfit',
    body: '<div id="player-outfit-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 300, showMax: false, resizable: false,
  })
    .then(() => form.render('#player-outfit-form'))
    .close(() => form.destroy())
}

function createPlayerItemsGrid() {
  return new w2grid({
    name: 'playerItemGrid',
    url: {
      get: '/api/v1/player/:id/item/grid',
      remove: '/api/v1/player/item/remove',
    },
    routeData: { id: selectedPlayerID ?? 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: true,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: false,
      toolbarReload: true,
      searchSave: false,
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
      },
      {
        field: 'oid',
        text: 'Player Item OID',
        size: '135px',
        render: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'item',
        text: 'Item',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'player_avatar',
        text: 'Player Avatar',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'player_outfit',
        text: 'Player Outfit',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'avatar_slot',
        text: 'Avatar Slot',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'quantity',
        text: 'Quantity',
        size: '120px',
        render: 'int',
      },
    ],
    onRender: function(event) {
      if (selectedPlayerID == null) {
        event.owner.toolbar.disable('w2ui-add')
      }
    },
    onAdd: function(event) { openPlayerItemPopup(event) },
    onEdit: function(event) { openPlayerItemPopup(event) },
    onDblClick: function(event) { openPlayerItemPopup(event) },
  })
}

function openPlayerItemPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'playerItemForm',
    url: '/api/v1/player/:id/item/form',
    routeData: { id: selectedPlayerID ?? 0 },
    focus: -1,
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
        field: 'oid',
        type: 'text',
        required: isEditMode,
        html: {
          label: 'Player Item OID',
          attr: isEditMode ? 'size="15"' : 'size="15" readonly',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'item',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/item'),
        html: {
          label: 'Item',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'player_avatar',
        type: 'list',
        options: helpers.remoteListOptions(`/api/v1/player/${selectedPlayerID}/avatar`),
        html: {
          label: 'Player Avatar',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'player_outfit',
        type: 'list',
        options: helpers.remoteListOptions(`/api/v1/player/${selectedPlayerID}/outfit`),
        html: {
          label: 'Player Outfit',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'avatar_slot',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/avatar/slot'),
        html: {
          label: 'Avatar Slot',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'quantity',
        type: 'int',
        required: true,
        html: {
          label: 'Quantity',
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
    title: isEditMode ? 'Edit Player Item' : 'New Player Item',
    body: '<div id="player-item-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 420, showMax: false, resizable: false,
  })
    .then(() => form.render('#player-item-form'))
    .close(() => form.destroy())
}

