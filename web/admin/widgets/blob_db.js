import { w2confirm, w2grid } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createWidget() {
  return new w2grid({
    name: 'blobGrid',
    url: {
      get: '/api/v1/blob/grid',
      remove: '/api/v1/blob/remove',
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
          text: 'Import',
          tooltip: 'Import cache files from the ./cache folder',
          icon: 'fa fa-file-arrow-down',
          onClick: function() {
            w2confirm({
              title: 'Import Cache Files',
              msg: 'This will import all files from the `cache` folder into blob.db.',
              btn_yes: { text: 'Import', class: 'w2ui-btn-blue' },
              btn_no: { text: 'Cancel' },
            }).yes(async () => {
              await new Promise(r => setTimeout(r, 300));
              await helpers.w2fetch({
                owner: this.owner,
                reload: true,
                lock: 'Importing cache files...',
                url: '/api/v1/blob/import',
                method: 'POST',
              })
            })
          },
        },
        {
          type: 'button',
          id: 'export',
          text: 'Export',
          tooltip: 'Export cache files to the ./cache folder',
          icon: 'fa fa-file-arrow-up',
          onClick: function() {
            w2confirm({
              title: 'Export Cache Files',
              msg: 'This will create a `cache` folder containing all assets.',
              btn_yes: { text: 'Export', class: 'w2ui-btn-blue' },
              btn_no: { text: 'Cancel' },
            }).yes(async () => {
              await new Promise(r => setTimeout(r, 300));
              await helpers.w2fetch({
                owner: this.owner,
                reload: false,
                lock: 'Exporting cache files...',
                url: '/api/v1/blob/export',
                method: 'POST',
              })
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
        field: 'url',
        text: 'File URL',
        size: '400px',
        render: 'text',
        sortable: true,
        clipboardCopy: true,
      },
      {
        field: 'hash',
        text: 'File Hash',
        size: '350px',
        render: 'text',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: true,
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
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function() {
      helpers.w2upload({
        owner: this,
        reload: true,
        lock: 'Uploading files...',
        url: '/api/v1/blob/upload',
        method: 'POST',
        multiple: true,
      })
    },
    onSearch: function(event) { helpers.searchAllFilter(event) },
  })
}

