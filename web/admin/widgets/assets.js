import { w2form, w2grid, w2layout, w2popup, w2utils, query } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createAssetGrid() {
  return new w2grid({
    name: 'assetGrid',
    url: {
      get: '/api/v1/asset/grid',
      save: '/api/v1/asset/grid',
    },
    recid: 'id',
    recordHeight: 28,
    multiSearch: true,
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
      expandColumn: true,
    },
    toolbar: {
      items: [
        { type: 'break' },
        {
          type: 'button',
          id: 'import-json',
          text: 'Import cache.json',
          tooltip: 'Import and merge base asset data from cache.json',
          icon: 'fa fa-code',
          onClick: function() {
            helpers.w2upload({
              owner: this.owner,
              reload: true,
              lock: 'Uploading json...',
              url: '/api/v1/asset/cache.json',
              method: 'POST',
              accept: '.json,application/json',
            })
          },
        },
      ],
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        sortable: true,
        searchAll: true,
        searchable: 'int',
      },
      {
        field: 'oid',
        text: 'Asset OID',
        size: '135px',
        render: 'text',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '200px',
        render: 'text',
        sortable: true,
        hidden: true,
      },
      {
        field: 'cdnid',
        text: 'CDN ID',
        size: '200px',
        render: 'text',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'file_type',
        text: 'File Type',
        size: '150px',
        render: 'dropdown',
        sortable: true,
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/asset/filetype') },
      },
      {
        field: 'asset_type',
        text: 'Asset Type',
        size: '150px',
        render: 'dropdown',
        sortable: true,
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/asset/assettype') },
        editable: helpers.remoteListOptions('/api/v1/asset/assettype'),
      },
      {
        field: 'asset_group',
        text: 'Asset Group',
        size: '150px',
        render: 'dropdown',
        sortable: true,
        searchable: { type: 'enum', options: helpers.remoteListOptions('/api/v1/asset/assetgroup') },
        editable: helpers.remoteListOptions('/api/v1/asset/assetgroup'),
      },
      {
        field: 'res_name',
        text: 'ResName',
        size: '200px',
        render: 'hover',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'description',
        text: 'Description',
        size: '200px',
        render: 'text',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'version',
        text: 'Version',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'hash',
        text: 'File Hash',
        size: '350px',
        render: 'text',
        sortable: true,
        hidden: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'metadata',
        text: 'Metadata',
        size: '350px',
        render: 'text',
        hidden: true,
        searchable: 'text',
      },
      {
        field: 'size',
        text: 'Bytes',
        size: '80px',
        render: 'text',
        sortable: true,
        searchable: 'int',
      },
      {
        field: 'size_str',
        text: 'Size',
        size: '80px',
        render: 'text',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'oid', direction: 'desc' },
    ],
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onSearch: function(event) { helpers.searchAllFilter(event) },
    onDelete: function(event) { event.preventDefault() },
    onExpand: function(event) {
      const row = event.owner.get(event.detail.recid)
      const box = query('#' + event.detail.box_id)
      box.html(`
        <div style="padding: 5px">
          <div style="height: 300px">
            Loading...
          </div>
        </div>
      `)
      if (row.file_type.text.startsWith('image/')) {
        box.html(`
          <div style="padding: 5px; height: 150px">
            <img style="height:100%" src="${row.url}"/>
          </div>
        `)
      } else if (row.file_type.text.startsWith('audio/')) {
        box.html(`
          <div style="padding: 5px;">
            <audio controls type="${row.file_type.text}" src="${row.url}"></audio>
          </div>
        `)
      } else if (
        row.file_type.text.startsWith('TreeNode/') ||
        row.file_type.text.includes('json')
      ) {
        fetch(row.url).then(async res => {
          const raw = await res.text()
          const text = w2utils.encodeTags(raw)
          box.html(`
            <div style="padding: 5px;">
              <textarea style="width: 100%; height: 300px; resize: none; font-family: monospace;" readonly>${text}</textarea>
            </div>
          `)
        })
      } else if (row.metadata) {
        const prettyJson = JSON.stringify(JSON.parse(row.metadata), null, 2)
        const text = w2utils.encodeTags(prettyJson)
        box.html(`
          <div style="padding: 5px;">
            <textarea style="width: 100%; height: 300px; resize: none; font-family: monospace;" readonly>${text}</textarea>
          </div>
        `)
      } else {
        box.html(`
          <div style="padding: 5px;">
            preview not available
          </div>
        `)
      }
    },
  })
}

export function createContainerLayout() {
  const containerGrid = new w2grid({
    name: 'containerGrid',
    url: {
      get: '/api/v1/container/grid',
      save: '/api/v1/container/grid',
      remove: '/api/v1/container/remove',
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
        field: 'oid',
        text: 'Container OID',
        render: 'text',
        size: '135px',
        sortable: true,
        searchable: 'text',
        editable: { type: 'int' },
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '200px',
        render: 'text',
        sortable: true,
        hidden: true,
      },
      {
        field: 'name',
        text: 'Container Name',
        render: 'text',
        size: '200px',
        sortable: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'ptag',
        text: 'PTag',
        size: '120px',
        render: 'text',
        sortable: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'assets',
        text: 'Assets',
        size: '60px',
        sortable: true,
      },
      {
        field: 'packages',
        text: 'Pkgs',
        size: '60px',
        sortable: true,
      },
      {
        field: 'created_at',
        text: 'Created at',
        size: '135px',
        render: 'datetime',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'oid', direction: 'desc' },
    ],
    onAdd: function(event) { openContainerPopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onSelect: function(event) { reloadSubGrids(event) },
    onDelete: function() { clearSubGrids() },
  })

  function clearSubGrids() {
    containerAssetGrid.routeData.id = 0
    containerPackageGrid.routeData.id = 0
    containerAssetGrid.clear()
    containerPackageGrid.clear()
    containerAssetGrid.toolbar.disable('w2ui-add')
    containerPackageGrid.toolbar.disable('w2ui-add')
  }

  async function reloadSubGrids(event) {
    await event.complete
    const selection = event.owner.getSelection()
    if (selection.length != 1) {
      clearSubGrids()
    } else {
      containerAssetGrid.routeData.id = selection[0]
      containerPackageGrid.routeData.id = selection[0]
      await Promise.all([containerAssetGrid.reload(), containerPackageGrid.reload()])
      containerAssetGrid.toolbar.enable('w2ui-add')
      containerPackageGrid.toolbar.enable('w2ui-add')
    }
  }

  const containerAssetGrid = new w2grid({
    name: 'containerAssetGrid',
    header: '<i class="fa fa-list"></i> Asset Map',
    url: {
      get: '/api/v1/container/:id/asset/grid',
      save: '/api/v1/container/asset/grid',
      remove: '/api/v1/container/asset/remove',
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
        hidden: true,
      },
      {
        field: 'win_asset',
        text: 'Windows Asset',
        render: 'dropdown-tooltip',
        size: '50%',
        editable: helpers.remoteListOptions('/api/v1/asset'),
      },
      {
        field: 'osx_asset',
        text: 'OSX Asset',
        render: 'dropdown-tooltip',
        size: '50%',
        editable: helpers.remoteListOptions('/api/v1/asset'),
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'position', direction: 'asc' },
    ],
    onAdd: function(event) { openContainerAssetPopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onReorderRow: function(event) { helpers.w2reorder(event, { url: '/api/v1/container/asset/reorder' }) },
  })

  const containerPackageGrid = new w2grid({
    name: 'containerPackageGrid',
    header: '<i class="fa fa-box-open"></i> Asset Packages',
    url: {
      get: '/api/v1/container/:id/package/grid',
      save: '/api/v1/container/package/grid',
      remove: '/api/v1/container/package/remove',
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
        hidden: true,
      },
      {
        field: 'pkg_container',
        text: 'Package Container',
        render: 'dropdown-tooltip',
        size: '100%',
        editable: helpers.remoteListOptions('/api/v1/container'),
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'position', direction: 'asc' },
    ],
    onAdd: function(event) { openContainerPackagePopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onReorderRow: function(event) { helpers.w2reorder(event, { url: '/api/v1/container/package/reorder' }) },
  })

  const subLayout = new w2layout({
    name: 'assetContainerSubLayout',
    style: 'margin: 5px;',
    padding: 5,
    panels: [
      { type: 'top', size: '50%', html: containerAssetGrid },
      { type: 'main', size: '50%', html: containerPackageGrid },
    ],
  })

  return new w2layout({
    name: 'assetContainerLayout',
    panels: [
      { type: 'left', size: '50%', resizable: true, html: containerGrid },
      { type: 'main', size: '50%', html: subLayout },
    ],
    onRender: function() { clearSubGrids() },
    onDestroy: function() {
      containerGrid.destroy()
      containerAssetGrid.destroy()
      containerPackageGrid.destroy()
      subLayout.destroy()
    }
  })
}

function openContainerPopup(event) {
  const form = new w2form({
    name: 'containerForm',
    url: '/api/v1/container/form',
    fields: [
      {
        field: 'oid',
        type: 'int',
        required: true,
        html: {
          label: 'Container OID',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Container Name',
          attr: 'style="width:100%;"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'ptag',
        type: 'text',
        html: {
          label: 'PTag',
          attr: 'style="width:100%;"',
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
    title: 'New Container',
    body: '<div id="container-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 300, showMax: false, resizable: false,
  })
    .then(() => form.render('#container-form'))
    .close(() => form.destroy())
}

function openContainerAssetPopup(event) {
  const form = new w2form({
    name: 'containerAssetForm',
    url: '/api/v1/container/:id/asset/form',
    routeData: { id: event.owner.routeData.id },
    focus: -1,
    fields: [
      {
        field: 'win_asset',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/asset'),
        html: {
          label: 'Windows Asset',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'osx_asset',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/asset'),
        html: {
          label: 'OSX Asset',
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
    title: 'Add Asset',
    body: '<div id="container-asset-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#container-asset-form'))
    .close(() => form.destroy())
}

function openContainerPackagePopup(event) {
  const form = new w2form({
    name: 'containerPackageForm',
    url: '/api/v1/container/:id/package/form',
    routeData: { id: event.owner.routeData.id },
    focus: -1,
    fields: [
      {
        field: 'pkg_container',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/container'),
        html: {
          label: 'Package Container',
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
    title: 'Add Package',
    body: '<div id="container-package-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 170, showMax: false, resizable: false,
  })
    .then(() => form.render('#container-package-form'))
    .close(() => form.destroy())
}

