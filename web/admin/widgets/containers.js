import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createContainerLayout() {
  const containerGrid = new w2grid({
    name: 'containerGrid',
    url: {
      get: '/api/v1/container/grid',
      remove: '/api/v1/container/remove',
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
        field: 'oid',
        text: 'Container OID',
        render: 'text',
        size: '135px',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
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
      },
      {
        field: 'ptag',
        text: 'PTag',
        size: '120px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'icon',
        text: '',
        size: '40px',
        render: 'icon-sm',
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
    onEdit: function(event) { openContainerPopup(event) },
    onDblClick: function(event) { openContainerPopup(event) },
    onSelect: function(event) { reloadSubGrids(event) },
    onDelete: function() { clearSubGrids() },
  })

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
        tooltip: 'macOS asset overrides the Windows asset on macOS; otherwise, the Windows asset is used on all platforms',
      },
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
    onAdd: function(event) { openContainerPackagePopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onReorderRow: function(event) { helpers.w2reorder(event, { url: '/api/v1/container/package/reorder' }) },
  })

  const subLayout = new w2layout({
    name: 'assetContainerSubLayout',
    style: 'margin: 5px;',
    padding: 5,
    panels: [
      { type: 'top', size: '50%', resizable: true, html: containerAssetGrid },
      { type: 'main', size: '50%', html: containerPackageGrid },
    ],
    onDestroy: function() {
      containerAssetGrid.destroy()
      containerPackageGrid.destroy()
    },
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

  return new w2layout({
    name: 'assetContainerLayout',
    panels: [
      { type: 'left', size: '50%', resizable: true, html: containerGrid },
      { type: 'main', size: '50%', html: subLayout },
    ],
    onRender: function() { clearSubGrids() },
    onDestroy: function() {
      containerGrid.destroy()
      subLayout.destroy()
    },
  })
}

function openContainerPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'containerForm',
    url: '/api/v1/container/form',
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
        required: true,
        html: {
          label: 'Container OID',
          attr: 'size="15"',
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
    title: isEditMode ? 'Edit Container' : 'New Container',
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

