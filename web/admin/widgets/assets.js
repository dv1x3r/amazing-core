import { w2confirm, w2grid } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createAssetGrid() {
  return new w2grid({
    name: 'assetGrid',
    url: {
      get: '/api/v1/asset/grid',
      save: '/api/v1/asset/grid',
      remove: '/api/v1/asset/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: false,
      toolbarEdit: false,
      toolbarDelete: true,
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
          id: 'import-blob',
          text: 'Import assets',
          icon: 'fa fa-file-import',
          onClick: function() {
            w2confirm({
              title: 'Import assets from blob.db',
              msg: `
                Scan blob.db and add missing assets to core.db.<br>
                Existing asset records and blob files are not changed.
              `,
              btn_yes: { text: 'Import', class: 'w2ui-btn-blue' },
              btn_no: { text: 'Cancel' },
            }).yes(async () => {
              await new Promise(r => setTimeout(r, 300));
              await helpers.w2fetch({
                owner: this.owner,
                reload: false,
                lock: 'Importing assets...',
                url: '/api/v1/asset/import',
                method: 'POST',
              })
              this.owner.reload()
            })
          },
        },
        {
          type: 'button',
          id: 'download-index',
          text: 'Download index file',
          icon: 'fa fa-download',
          onClick: async function() {
            await helpers.w2download({
              owner: this.owner,
              lock: 'Generating index file...',
              url: '/api/v1/asset/index',
              name: 'index.json',
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
        size: '135px',
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
        field: 'icon',
        text: x => x.field ? '' : 'Icon',
        size: '40px',
        render: 'icon-sm',
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
        text: 'Res Name',
        size: '200px',
        render: 'hover',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'bundle_version',
        text: 'Bundle Version',
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
        hidden: true,
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
    onSave: function(event) { helpers.reloadOnSuccess(event) },
    onSearch: function(event) { helpers.searchAllFilter(event) },
  })
}

