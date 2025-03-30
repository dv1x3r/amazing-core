import { w2form, w2grid, w2popup, w2utils, query } from '../lib/w2ui.es6.min.js'

export function createRandomNamesGrid() {
  return new w2grid({
    name: 'randomNamesGrid',
    url: {
      get: '/api/v1/randomnames/records',
      remove: '/api/v1/randomnames/remove',
    },
    httpHeaders: { 'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content },
    recid: 'id',
    multiSearch: true,
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
      { field: 'id', text: 'ID', size: '60px', sortable: true, hidden: true },
      { field: 'part_type', text: 'Part Type', size: '250px', render: 'safe', sortable: true },
      { field: 'name', text: 'Name', size: '250px', render: 'safe', sortable: true },
    ],
    searches: [
      { field: 'part_type', label: 'Part Type', type: 'text' },
      { field: 'name', label: 'Name', type: 'text' },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'part_type', direction: 'asc' },
      { field: 'name', direction: 'asc' },
    ],
    onAdd: function(event) { openRandomNamePopup(event) },
    onEdit: function(event) { openRandomNamePopup(event) },
    onDblClick: function(event) { openRandomNamePopup(event) },
  })
}

function openRandomNamePopup(event) {
  const randomNameForm = new w2form({
    name: `randomNameForm`,
    url: '/api/v1/randomnames/form',
    httpHeaders: { 'X-CSRF-Token': document.querySelector('meta[name="csrf-token"]').content },
    recid: event.detail.recid,
    fields: [
      { field: 'id', type: 'text', html: { label: 'ID', attr: 'size="10" readonly', span: 4, column: 0 } },
      { field: 'part_type', type: 'text', required: true, html: { label: 'Part Type', span: 4, column: 0 } },
      { field: 'name', type: 'text', required: true, html: { label: 'Name', span: 4, column: 0 } },
    ],
    actions: {
      async Save() {
        const res = await this.save()
        if (res.status == 'success') {
          this.recid = res.recid
          this.reload()
          event.owner.reload()
          query('.w2ui-popup-title').text('Edit Random Name')
          w2utils.notify('Data has been successfully saved!', { timeout: 4000 })
        }
      },
      Cancel() { w2popup.close() },
    },
  })

  w2popup.open({
    title: event.type == 'add' ? 'New Random Name' : 'Edit Random Name',
    body: '<div id="random-name-form" style="width: 100%; height: 100%;"></div>',
    width: 500, height: 300, showMax: true, resizable: true,
  })
    .then(() => randomNameForm.render('#random-name-form'))
    .close(() => randomNameForm.destroy())
}

