import { w2form, w2grid, w2layout, w2popup } from '/lib/w2ui.es6.min.js'
import * as helpers from '/lib/w2ui.helpers.js'

export function createCategoryLayout() {
  const grid = new w2grid({
    name: 'itemCategoryGrid',
    url: {
      get: '/api/v1/item/category/grid',
      save: '/api/v1/item/category/grid',
      remove: '/api/v1/item/category/remove',
    },
    recid: 'id',
    recordHeight: 28,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: true,
      toolbarEdit: false,
      toolbarDelete: true,
      toolbarSave: true,
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
        text: 'Category Name',
        size: '200px',
        render: 'text',
        sortable: true,
        searchable: 'text',
        editable: { type: 'text' },
      },
      {
        field: 'container',
        text: 'Asset Container',
        size: '250px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
        editable: helpers.remoteListOptions('/api/v1/container'),
      },
      {
        field: 'parent',
        text: 'Parent Category',
        size: '250px',
        render: 'dropdown',
        sortable: true,
        searchable: 'text',
        editable: helpers.remoteListOptions('/api/v1/item/category'),
      },
      {
        field: 'is_outdoor',
        text: 'Outdoor',
        size: '70px',
        editable: { type: 'checkbox' },
      },
      {
        field: 'is_walkover',
        text: 'Walkover',
        size: '70px',
        editable: { type: 'checkbox' },
      },
      {
        field: 'show_in_dock',
        text: 'Dock',
        size: '70px',
        editable: { type: 'checkbox' },
      },
      {
        field: 'is_public',
        text: 'Public',
        size: '70px',
        editable: { type: 'checkbox' },
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    sortData: [
      { field: 'id', direction: 'desc' },
    ],
    onAdd: function(event) { openCategoryPopup(event) },
    onSave: function(event) { helpers.reloadOnSuccess(event) },
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
  const form = new w2form({
    name: 'itemCategoryForm',
    url: '/api/v1/item/category/form',
    fields: [
      {
        field: 'name',
        type: 'text',
        required: true,
        html: {
          label: 'Category Name',
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
        field: 'parent',
        type: 'list',
        required: false,
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
    title: 'New Category',
    body: '<div id="category-form" style="width: 100%; height: 100%;"></div>',
    width: 600, height: 400, showMax: false, resizable: false,
  })
    .then(() => form.render('#category-form'))
    .close(() => form.destroy())
}

