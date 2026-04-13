import { w2confirm, w2form, w2grid, w2popup } from '/lib/w2ui.es6.min.js'
import { w2fetch, w2upload, searchAllFilter } from '/lib/w2ui.helpers.js'

export function createBlobGrid() {
  return new w2grid({
    name: 'blobGrid',
    url: {
      get: '/api/v1/blob/records',
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
              await w2fetch({
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
              await w2fetch({
                owner: this.owner,
                reload: false,
                lock: 'Exporting cache files...',
                url: '/api/v1/blob/export',
                method: 'POST',
              })
            })
          },
        },
        {
          type: 'button',
          id: 's3sync',
          text: 'Upload to S3',
          icon: 'fa fa-cloud-arrow-up',
          onClick: () => openS3SyncPopup(),
        },
      ],
    },
    columns: [
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        sortable: true,
        hidden: false,
        searchable: 'int',
      },
      {
        field: 'cdnid',
        text: 'CDN ID',
        render: 'text',
        size: '200px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'url',
        text: 'File URL',
        render: 'text',
        size: '400px',
        sortable: true,
        clipboardCopy: true,
      },
      {
        field: 'hash',
        text: 'File Hash',
        render: 'text',
        size: '350px',
        sortable: true,
        searchAll: true,
        searchable: 'text',
        clipboardCopy: true,
      },
      {
        field: 'size',
        text: 'Bytes',
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
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function() {
      w2upload({
        owner: this,
        reload: true,
        lock: 'Uploading files...',
        url: '/api/v1/blob/upload',
        method: 'POST',
        multiple: true,
      })
    },
    onSearch: function(event) { searchAllFilter(event) },
  })
}

function openS3SyncPopup() {
  const s3SyncForm = new w2form({
    name: 's3SyncForm',
    fields: [
      {
        field: 'endpoint',
        type: 'text',
        html: {
          label: 'Endpoint',
          attr: 'style="width:100%"; placeholder="https://s3.amazonaws.com";',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'region',
        type: 'text',
        required: true,
        html: {
          label: 'Region',
          attr: 'style="width:100%";',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'bucket',
        type: 'text',
        required: true,
        html: {
          label: 'Bucket',
          attr: 'style="width:100%";',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'path_prefix',
        type: 'text',
        html: {
          label: 'Path Prefix',
          attr: 'style="width:100%";',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'access_key_id',
        type: 'text',
        required: true,
        html: {
          label: 'Access Key ID',
          attr: 'style="width:100%";',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'secret_access_key',
        type: 'password',
        required: true,
        html: {
          label: 'Secret Access Key',
          attr: 'style="width:100%";',
          span: 6,
          column: 0,
        },
      },
    ],
    actions: {
      async Sync() {
        const errors = this.validate()
        if (errors.length > 0) {
          return
        }
        await w2fetch({
          owner: this,
          reload: false,
          lock: 'Syncing to S3...',
          url: '/api/v1/blob/s3sync',
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.record),
        })
      },
      Cancel() { w2popup.close() },
    },
  })

  w2popup.open({
    title: 'Upload Cache Files to S3',
    body: '<div id="s3-sync-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 400, showMax: false, resizable: false,
  })
    .then(() => s3SyncForm.render('#s3-sync-form'))
    .close(() => s3SyncForm.destroy())
}
