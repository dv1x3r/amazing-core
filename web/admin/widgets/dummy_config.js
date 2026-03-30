import { w2form } from '/lib/w2ui.es6.min.js'

export function createDummyForm() {
  return new w2form({
    name: `dummyForm`,
    url: '/api/v1/dummy/form',
    header: 'Dummy Config',
    recid: 1,
    toolbar: {
      items: [
        {
          id: 'save',
          type: 'button',
          text: 'Save',
          icon: 'fa fa-floppy-disk',
          onClick: event => {
            event.owner.owner.save()
          },
        },
      ],
    },
    fields: [
      {
        field: 'map',
        type: 'text',
        required: true,
        html: {
          label: 'Map',
          text: `
            <br>
            <br>
            <ul>
              <li><b>OTYwOTUyODk5OTk1MA</b> for Springbay</li>
              <li><b>OTYxMTQ4NDU5NDE5MA</b> for HomeLotSmall</li>
              <li><b>OTQ1MDc3NTY0MjEyNg</b> for HomeLot_Winter</li>
            </ul>
          `,
          span: 4,
          column: 0,
        },
      },
      {
        field: 'avatar',
        type: 'text',
        required: true,
        html: {
          label: 'Avatar',
          text: `
            <br>
            <br>
            You can search on <a href="https://amazingcore.org/cache-archive/">our cache archive.</a>
          `,
          span: 4,
          column: 0,
        },
      },
    ],
  })
}

