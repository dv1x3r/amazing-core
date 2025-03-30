import { w2grid, w2utils } from '../lib/w2ui.es6.min.js'

export function createBlobGrid() {
  return new w2grid({
    name: 'blobGrid',
    url: {
      get: '/api/v1/blob/records',
      remove: '/api/v1/blob/remove',
    },
    httpHeaders: { 'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content },
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
    columns: [
      { field: 'id', text: 'ID', size: '60px', sortable: true, hidden: false },
      { field: 'cdnid', text: 'CDN ID', size: '250px', render: 'safe', sortable: true, clipboardCopy: row => row.cdnid },
      { field: 'url', text: 'File URL', size: '450px', render: 'safe', sortable: true, clipboardCopy: row => row.url },
      { field: 'hash', text: 'File Hash', size: '370px', render: 'safe', sortable: true, clipboardCopy: row => row.hash },
      { field: 'size', text: 'Size Bytes', size: '88px', render: 'safe', sortable: true, hidden: true },
      { field: 'size_str', text: 'Size', size: '88px', render: 'safe', sortable: true },
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

        this.lock('Uploading...')

        const res = await fetch('/api/v1/blob/upload', {
          method: 'POST',
          headers: { 'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content },
          body: form,
        })

        this.unlock()

        if (res.status != 200) {
          res.json()
            .then(data => w2utils.notify(`${data?.message}, ${res.status}: ${res.statusText}`, { timeout: 4000, error: true }))
            .catch(() => w2utils.notify(`Failed to upload the files, ${res.status}: ${res.statusText}`, { timeout: 4000, error: true }))
          return
        }

        w2utils.notify('Files have been successfully uploaded!', { timeout: 4000 })
        await this.reload()
      }

      input.click()
    },
  })
}

