import { w2grid, w2utils, query } from '/lib/w2ui.es6.min.js'
import { w2upload, remoteListOptions, reloadOnSuccess, searchAllFilter } from '/lib/w2ui.helpers.js'

export function createAssetGrid() {
  return new w2grid({
    name: 'assetGrid',
    url: {
      get: '/api/v1/asset/records',
      save: '/api/v1/asset/save',
      remove: '/api/v1/asset/remove',
    },
    recid: 'id',
    recordHeight: 28,
    multiSearch: true,
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
          id: 'import-json',
          text: 'Import cache.json',
          tooltip: 'Import and update assets from cache.json',
          icon: 'fa fa-code',
          onClick: function() {
            w2upload({
              owner: this.owner,
              reload: true,
              lock: 'Uploading json...',
              url: '/api/v1/asset/cache.json',
              method: 'POST',
              accept: '.json,application/json',
            })
          },
        },
        {
          type: 'button',
          id: 'import-sql',
          text: 'Import assets.sql',
          tooltip: 'Replace your current asset table with the selected SQL dump (overwrites existing data)',
          icon: 'fa fa-database',
          onClick: function() {
            w2upload({
              owner: this.owner,
              reload: true,
              lock: 'Uploading sql...',
              url: '/api/v1/asset/assets.sql',
              method: 'POST',
              accept: '.sql',
            })
          },
        },
        {
          type: 'button',
          id: 'export-sql',
          text: 'Export assets.sql',
          tooltip: 'Export your asset table to an SQL file for backup, sharing, or creating a new Git version',
          icon: 'fa fa-database',
          onClick: function() {
            window.location.href = '/api/v1/asset/assets.sql'
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
        hidden: true,
      },
      {
        field: 'cdnid',
        text: 'CDN ID',
        render: 'text',
        size: '200px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: row => row.cdnid,
      },
      {
        field: 'url',
        text: 'File URL',
        render: 'text',
        size: '400px',
        hidden: true,
        sortable: true,
        clipboardCopy: row => row.url,
      },
      {
        field: 'gsfoid',
        text: 'GSF OID',
        render: 'text',
        size: '130px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: row => row.gsfoid,
      },
      {
        field: 'class',
        text: 'Class',
        render: 'text',
        size: '60px',
        sortable: true,
        hidden: true,
      },
      {
        field: 'type',
        text: 'Type',
        render: 'text',
        size: '60px',
        sortable: true,
        hidden: true,
      },
      {
        field: 'server',
        text: 'Server',
        render: 'text',
        size: '60px',
        sortable: true,
        hidden: false,
        searchable: 'int',
      },
      {
        field: 'number',
        text: 'Number',
        render: 'text',
        size: '130px',
        sortable: true,
        hidden: true,
        searchable: 'text',
      },
      {
        field: 'version',
        text: 'Version',
        render: 'text',
        size: '120px',
        sortable: true,
        hidden: true,
      },
      {
        field: 'file_type',
        text: 'File Type',
        render: 'dropdown',
        size: '150px',
        sortable: true,
        searchable: { type: 'enum', options: remoteListOptions('/api/v1/asset/filetypes') },
      },
      {
        field: 'asset_type',
        text: 'Asset Type',
        render: 'dropdown',
        size: '150px',
        sortable: true,
        searchable: { type: 'enum', options: remoteListOptions('/api/v1/asset/assettypes') },
        editable: remoteListOptions('/api/v1/asset/assettypes'),
      },
      {
        field: 'asset_group',
        text: 'Asset Group',
        render: 'dropdown',
        size: '150px',
        sortable: true,
        searchable: { type: 'enum', options: remoteListOptions('/api/v1/asset/assetgroups') },
        editable: remoteListOptions('/api/v1/asset/assetgroups'),
      },
      {
        field: 'res_name',
        text: 'ResName',
        render: 'hover',
        size: '200px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'description',
        text: 'Description',
        render: 'text',
        size: '200px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'hash',
        text: 'File Hash',
        render: 'text',
        size: '350px',
        sortable: true,
        hidden: true,
        searchable: 'text',
        clipboardCopy: row => row.hash,
      },
      {
        field: 'metadata',
        text: 'Metadata',
        render: 'text',
        size: '350px',
        hidden: true,
        searchable: 'text',
      },
      {
        field: 'size',
        text: 'Size Bytes',
        render: 'text',
        size: '80px',
        sortable: true,
        searchable: 'int',
      },
      {
        field: 'size_str',
        text: 'Size',
        render: 'text',
        size: '80px',
        sortable: true,
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'cdnid', direction: 'desc' },
    ],
    onSave: function(event) { reloadOnSuccess(event) },
    onSearch: function(event) { searchAllFilter(event) },
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

