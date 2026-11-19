import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createItemCategoryLayout() {
  const grid = new w2grid({
    name: 'itemCategoryGrid',
    header: '<i class="fa fa-tags"></i> Item Categories',
    url: {
      get: '/api/v1/item/category/grid',
      remove: '/api/v1/item/category/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      header: true,
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
        hidden: true,
        sortable: true,
        searchable: 'int',
      },
      {
        field: 'oid',
        text: 'Category OID',
        render: 'text',
        size: '135px',
        hidden: true,
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'oid_str',
        text: 'OID Details',
        size: '135px',
        render: 'text',
        hidden: true,
      },
      {
        field: 'name',
        text: 'Category Name',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'parent',
        text: 'Parent Category',
        size: '200px',
        render: 'dropdown',
      },
      {
        field: 'is_public',
        text: 'Public',
        size: '70px',
        render: 'toggle',
      },
      {
        field: 'is_outdoor',
        text: 'Outdoor',
        size: '70px',
        render: 'toggle',
      },
      {
        field: 'is_walkover',
        text: 'Walkover',
        size: '70px',
        render: 'toggle',
      },
      {
        field: 'show_in_dock',
        text: 'Dock',
        size: '70px',
        render: 'toggle',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'name', direction: 'asc' },
    ],
    onAdd: function(event) { openItemCategoryPopup(event) },
    onEdit: function(event) { openItemCategoryPopup(event) },
    onDblClick: function(event) { openItemCategoryPopup(event) },
  })

  return new w2layout({
    name: 'itemCategoryLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/item_categories.html')
    },
    onDestroy: function() {
      grid.destroy()
    },
  })
}

function openItemCategoryPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: 'itemCategoryForm',
    url: '/api/v1/item/category/form',
    record: record,
    fields: [
      {
        field: 'id',
        type: 'text',
        html: {
          label: 'ID',
          attr: 'size="15" readonly',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'oid',
        type: 'text',
        html: {
          label: 'Category OID',
          attr: 'size="20" placeholder="Automatic"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Category Name',
          attr: 'style="width:100%;"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'parent',
        type: 'list',
        options: helpers.remoteListOptions('/api/v1/item/category/dropdown'),
        html: {
          label: 'Parent Category',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'is_public',
        type: 'checkbox',
        html: {
          label: 'Is Public',
          span: 2,
          column: 0,
        },
      },
      {
        field: 'is_outdoor',
        type: 'checkbox',
        html: {
          label: 'Is Outdoor',
          span: 2,
          column: 0,
        },
      },
      {
        field: 'is_walkover',
        type: 'checkbox',
        html: {
          label: 'Is Walkover',
          span: 2,
          column: 0,
        },
      },
      {
        field: 'show_in_dock',
        type: 'checkbox',
        html: {
          label: 'Show in Dock',
          span: 2,
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
    title: isEditMode ? 'Edit Item Category' : 'New Item Category',
    body: '<div id="popup-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 420, showMax: false, resizable: false,
  })
    .then(() => form.render('#popup-form'))
    .close(() => form.destroy())
}

