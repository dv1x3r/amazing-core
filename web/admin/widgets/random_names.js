import { w2form, w2grid, w2popup } from '/lib/w2ui.es6.min.js'

export function createRandomNameGrid() {
  return new w2grid({
    name: 'randomNameGrid',
    url: {
      get: '/api/v1/randname/records',
      remove: '/api/v1/randname/remove',
    },
    recid: 'id',
    recordHeight: 28,
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
      {
        field: 'id',
        text: 'ID',
        size: '60px',
        sortable: true,
        hidden: true,
      },
      {
        field: 'part_type',
        text: 'Part Type',
        size: '250px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'name',
        text: 'Name',
        size: '250px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
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
    url: '/api/v1/randname/form',
    recid: event.detail.recid,
    fields: [
      {
        field: 'id',
        type: 'text',
        html: {
          label: 'ID',
          attr: 'size="10" readonly',
          span: 4,
          column: 0,
        },
      },
      {
        field: 'part_type',
        type: 'text',
        required: true,
        html: {
          label: 'Part Type',
          span: 4,
          column: 0,
        },
      },
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Name',
          span: 4,
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
    title: event.type == 'add' ? 'New Random Name' : 'Edit Random Name',
    body: '<div id="random-name-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 300, showMax: false, resizable: false,
  })
    .then(() => randomNameForm.render('#random-name-form'))
    .close(() => randomNameForm.destroy())
}

