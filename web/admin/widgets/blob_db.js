import { w2confirm, w2form, w2grid, w2popup, w2utils } from '../lib/w2ui.es6.min.js'

export function createBlobGrid() {
  return new w2grid({
    name: 'blobGrid',
    url: {
      get: '/api/v1/blob/records',
      remove: '/api/v1/blob/remove',
    },
    recid: 'id',
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
              this.owner.lock({ spinner: true, msg: 'Importing cache files...' })
              try {
                const res = await fetch('/api/v1/blob/import', { method: 'POST' })
                if (!res.ok) {
                  const err = await res.json().catch(() => {
                    return { message: res.statusText }
                  })
                  throw new Error(err.message)
                }
                const result = await res.json()
                w2utils.notify(result.message, { timeout: 6000 })
              }
              catch (err) {
                this.owner.message(err.toString())
              }
              finally {
                this.owner.unlock()
              }
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
              this.owner.lock({ spinner: true, msg: 'Exporting cache files...' })
              try {
                const res = await fetch('/api/v1/blob/export', { method: 'POST' })
                if (!res.ok) {
                  const err = await res.json().catch(() => {
                    return { message: res.statusText }
                  })
                  throw new Error(err.message)
                }
                const result = await res.json()
                w2utils.notify(result.message, { timeout: 6000 })
              }
              catch (err) {
                this.owner.message(err.toString())
              }
              finally {
                this.owner.unlock()
              }
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
      },
      {
        field: 'cdnid',
        text: 'CDN ID',
        render: 'safe',
        size: '250px',
        sortable: true,
        clipboardCopy: row => row.cdnid,
      },
      {
        field: 'url',
        text: 'File URL',
        render: 'safe',
        size: '450px',
        sortable: true,
        clipboardCopy: row => row.url,
      },
      {
        field: 'hash',
        text: 'File Hash',
        render: 'safe',
        size: '370px',
        sortable: true,
        clipboardCopy: row => row.hash,
      },
      {
        field: 'size',
        text: 'Size Bytes',
        render: 'safe',
        size: '88px',
        sortable: true,
      },
      {
        field: 'size_str',
        text: 'Size',
        render: 'safe',
        size: '88px',
        sortable: true,
      },
    ],
    searches: [
      { field: 'cdnid', label: 'CDN ID', type: 'text' },
      { field: 'hash', label: 'File Hash', type: 'text' },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function() {
      const input = document.createElement('input')
      input.type = 'file'
      input.multiple = true

      input.onchange = async event => {
        const form = new FormData()
        for (const file of event.target.files) {
          form.append('files[]', file)
        }

        this.lock({ spinner: true, msg: 'Uploading files...' })

        try {
          const res = await fetch('/api/v1/blob/upload', {
            method: 'POST',
            body: form,
          })

          if (!res.ok) {
            const err = await res.json().catch(() => {
              return { message: res.statusText }
            })
            throw new Error(err.message)
          }

          w2utils.notify('Upload completed!', { timeout: 6000 })
          this.reload()
        }
        catch (err) {
          this.message(err.toString())
        }
        finally {
          this.unlock()
        }
      }

      input.click()
    },
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
      async Upload() {
        const errors = this.validate()
        if (errors.length > 0) {
          return
        }

        const record = this.record
        const payload = {
          endpoint: record.endpoint,
          region: record.region,
          bucket: record.bucket,
          path_prefix: record.path_prefix,
          access_key_id: record.access_key_id,
          secret_access_key: record.secret_access_key,
        }

        this.lock({ spinner: true, msg: 'Syncing to S3...' })

        try {
          const res = await fetch('/api/v1/blob/s3sync', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
          })

          if (!res.ok) {
            const err = await res.json().catch(() => {
              return { message: res.statusText }
            })
            throw new Error(err.message)
          }

          const result = await res.json()
          w2utils.notify(result.message, { timeout: 6000 })
        }
        catch (err) {
          this.message(err.toString())
        }
        finally {
          this.unlock()
        }
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
