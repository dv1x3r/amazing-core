import { w2grid, w2utils, query } from '/lib/w2ui.es6.min.js'

const maxRecords = 1000

export function createLogWatcherGrid() {
  let source = null
  let pending = []

  const grid = new w2grid({
    name: 'logWatcherGrid',
    recid: 'id',
    recordHeight: 24,
    show: {
      footer: true,
      toolbar: true,
      toolbarAdd: false,
      toolbarEdit: false,
      toolbarDelete: false,
      toolbarSave: false,
      toolbarSearch: true,
      toolbarReload: false,
      searchSave: false,
      expandColumn: true,
    },
    toolbar: {
      items: [
        {
          type: 'button',
          id: 'watch',
          text: 'Start Watch',
          icon: 'fa fa-play',
          onClick: function() {
            if (source) {
              stopWatch()
            } else {
              startWatch()
            }
          },
        },
        {
          type: 'button',
          id: 'clear',
          text: 'Clear',
          icon: 'fa fa-trash',
          onClick: function() {
            grid.clear()
          },
        },
      ],
    },
    columns: [
      {
        field: 'time',
        text: 'Time',
        size: '160px',
        render: 'datetime-local-ms',
      },
      {
        field: 'kind',
        text: 'Kind',
        size: '120px',
        render: 'text',
      },
      {
        field: 'remote_ip',
        text: 'Remote IP',
        size: '100px',
        render: 'text',
      },
      {
        field: 'player_oid',
        text: 'Player OID',
        size: '80px',
        render: 'text',
      },
      {
        field: 'service_class',
        text: 'Service Class',
        size: '120px',
        render: 'text',
      },
      {
        field: 'message_type',
        text: 'Message Type',
        size: '200px',
        render: 'text',
      },
      {
        field: 'result',
        text: 'Result',
        size: '100px',
        render: 'text',
      },
      {
        field: 'app',
        text: 'App',
        size: '100px',
        render: 'text',
      },
      {
        field: 'request_id',
        text: 'Req ID',
        size: '60px',
        render: 'int',
      },
      {
        field: 'request_flags',
        text: 'Req FL',
        size: '60px',
        render: 'int',
      },
      {
        field: 'response_flags',
        text: 'Res FL',
        size: '60px',
        render: 'int',
      },
      {
        field: 'latency',
        text: 'Latency',
        size: '60px',
        render: 'text',
      },
    ],
    defaultOperator: {
      'text': 'contains',
    },
    onDestroy: function() {
      stopWatch()
    },
    onExpand: function(event) {
      const row = event.owner.get(event.detail.recid)
      const box = query('#' + event.detail.box_id)
      const details = JSON.stringify(row.details ?? {}, null, 2)
      const text = w2utils.encodeTags(details)
      box.html(`
        <div style="padding: 5px;">
          <textarea style="width: 100%; height: 300px; resize: none; font-family: monospace;" readonly>${text}</textarea>
        </div>
      `)
    },
  })

  function startWatch() {
    if (source) {
      return
    }

    source = new EventSource('/api/v1/logs/watch')
    source.addEventListener('log', event => {
      const record = JSON.parse(event.data)
      const entry = recordToEntry(record)
      if (!entry) {
        return
      }
      if (pending.length == 0) {
        requestAnimationFrame(flushPending)
      }
      pending.push(entry)
      // requestAnimationFrame is paused in a hidden tab, keep only the latest records
      if (pending.length > maxRecords) {
        pending.splice(0, pending.length - maxRecords)
      }
    })
    source.onopen = () => {
      if (source) {
        grid.status('Connection is established')
      }
    }
    source.onerror = () => {
      if (!source) {
        return
      }
      // EventSource reconnects only after network errors, an error response closes it for good
      if (source.readyState == EventSource.CLOSED) {
        stopWatch()
        grid.status('Connection is closed')
      } else {
        grid.status('Reconnecting...')
      }
    }
    setWatchButton(true)
  }

  function flushPending() {
    const overflow = grid.records.length + pending.length - maxRecords
    if (overflow > 0) {
      grid.records.splice(-overflow)
    }
    grid.add(pending, true)
    grid.refresh()
    pending = []
  }

  function stopWatch() {
    if (!source) {
      return
    }

    source.close()
    source = null
    grid.status('')
    setWatchButton(false)
  }

  function setWatchButton(watching) {
    grid.toolbar.set('watch', watching
      ? { text: 'Stop Watch', icon: 'fa fa-stop' }
      : { text: 'Start Watch', icon: 'fa fa-play' }
    )
  }

  return grid
}

function recordToEntry(record) {
  if (record.message !== 'gsf service' && record.message !== 'gsf notify' && record.message !== 'gsf unhandled') {
    return null
  }

  const attrs = record.attrs ?? {}
  return {
    id: record.id,
    time: record.time,
    kind: attrs['kind'],
    remote_ip: attrs['remote_ip'],
    player_oid: attrs['player_oid'],
    service_class: `${attrs['svc_class'] ?? ''} ${attrs['svc_class_text'] ?? ''}`,
    message_type: `${attrs['msg_type'] ?? ''} ${attrs['msg_type_text'] ?? ''}`,
    result: `${attrs['result_code'] ?? ''} ${attrs['result_code_text'] ?? ''}`,
    app: `${attrs['app_code'] ?? ''} ${attrs['app_code_text'] ?? ''}`,
    request_id: attrs['request_id'],
    request_flags: attrs['flags'],
    response_flags: attrs['response_flags'],
    latency: attrs['latency'],
    details: record,
  }
}

