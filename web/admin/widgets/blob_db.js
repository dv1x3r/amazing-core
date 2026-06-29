import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createBlobLayout() {
  const grid = new w2grid({
    name: 'blobGrid',
    url: {
      get: '/api/v1/blob/grid',
      remove: '/api/v1/blob/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: false,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: false,
      toolbarSearch: true,
      toolbarReload: true,
      searchSave: false,
    },
    toolbar: {
      items: [
        { type: 'break' },
        {
          type: 'button',
          id: 'import',
          text: 'Import files',
          icon: 'fa fa-file-import',
          onClick: function() {
            openImportFilesPopup()
          },
        },
        {
          type: 'button',
          id: 'extract',
          text: 'Extract files',
          icon: 'fa fa-box-open',
          onClick: function() {
            openExtractFilesPopup()
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
        field: 'cdnid',
        text: 'CDN ID',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'hash',
        text: 'File Hash',
        size: '350px',
        render: 'text',
        searchable: 'text',
        clipboardCopy: true,
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
      { field: 'id', direction: 'desc' },
    ],
    onSelect: async function(event) {
      await event.complete
      const selection = event.owner.getSelection()
      const selectedBlobID = selection.length == 1 ? selection[0] : null
      if (selectedBlobID) {
        form.setValue('metadata', event.owner.get(selectedBlobID).metadata)
      } else {
        form.clear()
      }
    },
  })

  const form = new w2form({
    name: `blobDetailsForm`,
    focus: -1,
    fields: [
      {
        field: 'metadata',
        type: 'textarea',
        html: {
          label: '',
          group: 'Metadata',
          attr: 'style="width: 100%; height: calc(100vh - 85px); resize: none;" readonly',
          span: 0,
          column: 0,
        },
      },
    ],
  })

  return new w2layout({
    name: 'blobLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main', html: form },
    ],
    onDestroy: function() {
      grid.destroy()
      form.destroy()
    },
  })
}

function openImportFilesPopup() {
  const userAgent = navigator.userAgent.toLowerCase()
  const isWindows = userAgent.includes('win')
  const form = new w2form({
    name: 'importFilesForm',
    url: '/api/v1/blob/import',
    fields: [
      {
        field: 'import_path',
        type: 'text',
        required: true,
        html: {
          label: 'Import Path',
          attr: isWindows
            ? 'style="width:100%;" placeholder="C:\\Users\\Bloom\\amazing-world\\cache"'
            : 'style="width:100%;" placeholder="/home/bloom/amazing-world/cache"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'generate_metadata',
        type: 'checkbox',
        html: {
          text: 'Generate metadata',
          label: '<span style="font-size:11px;">&nbsp;&nbsp;(python3 is required)</span>',
          span: 6,
          column: 0,
        },
      },
    ],
    actions: {
      import: {
        text: 'Import',
        class: 'w2ui-btn-blue',
        async onClick() {
          const res = await this.save()
          if (res) {
            this.message(res.message)
          }
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: 'Import files',
    body: '<div id="import-files-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#import-files-form'))
    .close(() => form.destroy())
}

function openExtractFilesPopup() {
  const userAgent = navigator.userAgent.toLowerCase()
  const isWindows = userAgent.includes('win')
  const form = new w2form({
    name: 'extractFilesForm',
    url: '/api/v1/blob/extract',
    fields: [
      {
        field: 'extract_path',
        type: 'text',
        required: true,
        html: {
          label: 'Extract Path',
          attr: isWindows
            ? 'style="width:100%;" placeholder="C:\\Users\\Bloom\\amazing-world\\cache"'
            : 'style="width:100%;" placeholder="/home/bloom/amazing-world/cache"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'extract_metadata',
        type: 'checkbox',
        html: {
          text: 'Extract metadata',
          label: '<span style="font-size:11px;">&nbsp;&nbsp;(.meta.json files)</span>',
          span: 6,
          column: 0,
        },
      },
    ],
    actions: {
      extract: {
        text: 'Extract',
        class: 'w2ui-btn-blue',
        async onClick() {
          const res = await this.save()
          if (res) {
            this.message(res.message)
          }
        }
      },
      Cancel() { w2popup.close() },
    },
  })
  w2popup.open({
    title: 'Extract files',
    body: '<div id="extract-files-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 220, showMax: false, resizable: false,
  })
    .then(() => form.render('#extract-files-form'))
    .close(() => form.destroy())
}

