import { w2grid, w2utils, query } from '/lib/w2ui.es6.min.js'
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
        text: '',
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
        field: 'version',
        text: 'Bundle Version',
        size: '200px',
        render: 'text',
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
      { field: 'id', direction: 'desc' },
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

