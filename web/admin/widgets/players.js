import { w2form, w2grid, w2layout, w2popup, w2utils } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

let selectedPlayerID = null

export function createWidget() {
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
        hidden: true,
        searchable: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '200px',
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
            { id: 'outfit-items', text: 'Outfit Items' },
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
      case 'information':
        subLayout.load('main', '/admin/pages/players.html')
        break
    }
  }

  return new w2layout({
    name: 'playerLayout',
    panels: [
      { type: 'left', size: '320px', resizable: true, html: playerListGrid },
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
          attr: 'size="10" readonly',
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
        type: 'int',
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
      save: '/api/v1/player/avatar/grid',
      remove: '/api/v1/player/avatar/remove',
    },
    routeData: { id: selectedPlayerID ?? 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: true,
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
        editable: { type: 'int' },
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '200px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'name',
        text: 'Player Avatar Name',
        render: 'text',
        size: '200px',
        editable: { type: 'text' },
      },
      {
        field: 'avatar',
        text: 'Avatar',
        size: '200px',
        render: 'dropdown',
        editable: helpers.remoteListOptions(`/api/v1/avatar`),
      },
      {
        field: 'outfit_no',
        text: 'Outfit No',
        size: '120px',
        render: 'int',
        editable: { type: 'int' },
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
    onSave: function(event) { helpers.reloadOnSuccess(event) },
  })
}

function openPlayerAvatarPopup(event) {
  const form = new w2form({
    name: 'playerAvatarForm',
    url: '/api/v1/player/:id/avatar/form',
    routeData: { id: selectedPlayerID ?? 0 },
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
    title: 'New Player Avatar',
    body: '<div id="player-avatar-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#player-avatar-form'))
    .close(() => form.destroy())
}

function createPlayerOutfitsGrid() {
  return new w2grid({
    name: 'playerOutfitGrid',
    url: {
      get: '/api/v1/player/:id/outfit/grid',
      save: '/api/v1/player/outfit/grid',
      remove: '/api/v1/player/outfit/remove',
    },
    routeData: { id: selectedPlayerID ?? 0 },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: true,
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
        editable: { type: 'int' },
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '200px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'player_avatar',
        text: 'Player Avatar',
        size: '200px',
        render: 'dropdown',
        editable: helpers.remoteListOptions(`/api/v1/player/${selectedPlayerID}/avatar`),
      },
      {
        field: 'outfit_no',
        text: 'Outfit No',
        size: '120px',
        render: 'int',
        editable: { type: 'int' },
      },
    ],
    onRender: function(event) {
      if (selectedPlayerID == null) {
        event.owner.toolbar.disable('w2ui-add')
      }
    },
    onAdd: function(event) { openPlayerOutfitPopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
  })
}

function openPlayerOutfitPopup(event) {
  const form = new w2form({
    name: 'playerOutfitForm',
    url: '/api/v1/player/outfit/form',
    focus: -1,
    fields: [
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
    title: 'New Player Outfit',
    body: '<div id="player-outfit-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#player-outfit-form'))
    .close(() => form.destroy())
}
