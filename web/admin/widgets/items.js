import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createItemLayout() {
  const grid = new w2grid({
    name: 'itemGrid',
    url: {
      get: '/api/v1/item/grid',
      remove: '/api/v1/item/remove',
    },
    recid: 'id',
    recordHeight: 28,
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
        searchable: 'int',
      },
      {
        field: 'name',
        text: 'Item Name',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '200px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
      },
      {
        field: 'categories',
        text: 'Categories',
        size: '200px',
        render: 'dropdown-multi',
      },
      {
        field: 'slots',
        text: 'Acceptable Slots',
        size: '200px',
        render: 'dropdown-multi',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'asc' },
    ],
    onAdd: function(event) { openItemPopup(event) },
    onEdit: function(event) { openItemPopup(event) },
    onDblClick: function(event) { openItemPopup(event) },
  })

  return new w2layout({
    name: 'itemLayout',
    panels: [
      { type: 'left', html: grid, resizable: true, size: -420 },
      { type: 'main' },
    ],
    onRender: async function(event) {
      await event.complete
      event.owner.load('main', '/admin/pages/items.html')
    },
    onDestroy: function() {
      grid.destroy()
    },
  })
}

function openItemPopup(event) {
  const record = event.owner.get(event.detail.recid)
  const isEditMode = record != null
  const form = new w2form({
    name: `itemForm`,
    url: '/api/v1/item/form',
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
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Item Name',
          attr: 'style="width:100%;"',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'container',
        type: 'list',
        required: true,
        options: helpers.remoteListOptions('/api/v1/container'),
        html: {
          label: 'Asset Container',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'categories',
        type: 'enum',
        options: helpers.remoteListOptions('/api/v1/item/category'),
        html: {
          label: 'Categories',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
          column: 0,
        },
      },
      {
        field: 'slots',
        type: 'enum',
        options: helpers.remoteListOptions('/api/v1/avatar/slot'),
        html: {
          label: 'Acceptable Slots',
          attr: 'style="width:100%;" placeholder="Type to search..."',
          span: 6,
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
    title: isEditMode ? 'Edit Item' : 'New Item',
    body: '<div id="item-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 340, showMax: false, resizable: false,
  })
    .then(() => form.render('#item-form'))
    .close(() => form.destroy())
}

export function createCategoryLayout() {
  const grid = new w2grid({
    name: 'itemCategoryGrid',
    url: {
      get: '/api/v1/item/category/grid',
      remove: '/api/v1/item/category/remove',
    },
    recid: 'id',
    recordHeight: 28,
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
        searchable: 'int',
      },
      {
        field: 'oid',
        text: 'Category OID',
        render: 'text',
        size: '135px',
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
      { field: 'id', direction: 'asc' },
    ],
    onAdd: function(event) { openCategoryPopup(event) },
    onEdit: function(event) { openCategoryPopup(event) },
    onDblClick: function(event) { openCategoryPopup(event) },
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

function openCategoryPopup(event) {
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
        required: isEditMode,
        html: {
          label: 'Category OID',
          attr: isEditMode ? 'size="15"' : 'size="15" readonly',
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
        options: helpers.remoteListOptions('/api/v1/item/category'),
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
    body: '<div id="item-category-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 420, showMax: false, resizable: false,
  })
    .then(() => form.render('#item-category-form'))
    .close(() => form.destroy())
}

