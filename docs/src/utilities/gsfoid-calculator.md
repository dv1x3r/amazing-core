# GSFOID Calculator

<div style="display:flex; flex-direction: column; gap: 1.5rem; max-width: 420px;">
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      CDN ID
      <input id="oid-cdn" type="text"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      GSF OID
      <input id="oid-gsf" type="text"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      Class <span style="color:#999; font-size:1.25rem;">(byte 0–255)</span>
      <input id="oid-class" type="number" min="0" max="255" placeholder="0–255"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      Type <span style="color:#999; font-size:1.25rem;">(byte 0–255)</span>
      <input id="oid-type" type="number" min="0" max="255" placeholder="0–255"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      Server <span style="color:#999; font-size:1.25rem;">(byte 0–255)</span>
      <input id="oid-server" type="number" min="0" max="255" placeholder="0–255"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
    <label style="display: flex; flex-direction: column; gap: 0.5rem; font-size: 1.75rem;">
      Number <span style="color:#999; font-size:1.25rem;">(40-bit, 0–1099511627775)</span>
      <input id="oid-number" type="number" min="0" max="1099511627775" placeholder="0–1099511627775"
        style="font-family: monospace; padding: 0.4rem 0.6rem; border: 1px solid #bbb; border-radius: 3px; font-size:1.75rem;">
    </label>
  <div id="oid-error" style="margin-top: 0.75rem; color: #c00; font-weight: bold;"></div>
</div>

<script>
  const BYTE_MAX    = 255n
  const NUMBER_MAX  = 0xFFFFFFFFFFn        // 40-bit max
  const OID_MAX     = 0x7FFFFFFFFFFFFFFFn  // int64 max

  const CLASS_SHIFT  = 56n
  const TYPE_SHIFT   = 48n
  const SERVER_SHIFT = 40n

  const $ = id => document.getElementById(id)

  function setError(msg) {
    $('oid-error').textContent = msg || ''
  }

  function parseBigInt(value) {
    return BigInt(typeof value === 'bigint' ? value : String(value).trim())
  }

  function assertInRange(value, max) {
    if (value < 0n || value > max) {
      throw new RangeError(`${value} out of range 0–${max}`)
    }
    return value
  }

  function oidFromLong(long) {
    const n = assertInRange(parseBigInt(long), OID_MAX)
    return {
      objectClass:  (n >> CLASS_SHIFT)  & 0xFFn,
      objectType:   (n >> TYPE_SHIFT)   & 0xFFn,
      server:       (n >> SERVER_SHIFT) & 0xFFn,
      objectNumber: (n & NUMBER_MAX),
    }
  }

  function oidToLong({ objectClass, objectType, server, objectNumber }) {
    const c = assertInRange(parseBigInt(objectClass),  BYTE_MAX)   << CLASS_SHIFT
    const t = assertInRange(parseBigInt(objectType),   BYTE_MAX)   << TYPE_SHIFT
    const s = assertInRange(parseBigInt(server),       BYTE_MAX)   << SERVER_SHIFT
    const n = assertInRange(parseBigInt(objectNumber), NUMBER_MAX)
    return (c | t | s | n)
  }

  $('oid-cdn').addEventListener('input', () => {
    const cdnid = $('oid-cdn').value.trim()
    try {
      const long = atob(cdnid.trim())
      const oid = oidFromLong(long)
      $('oid-gsf').value    = long
      $('oid-class').value  = oid.objectClass
      $('oid-type').value   = oid.objectType
      $('oid-server').value = oid.server
      $('oid-number').value = oid.objectNumber
      setError('')
    } catch (error) {
      setError('Invalid CDN ID: ' + error)
    }
  })

  $('oid-gsf').addEventListener('input', () => {
    const long = $('oid-gsf').value.trim()
    try {
      const oid = oidFromLong(long)
      $('oid-cdn').value    = btoa(parseBigInt(long).toString())
      $('oid-class').value  = oid.objectClass
      $('oid-type').value   = oid.objectType
      $('oid-server').value = oid.server
      $('oid-number').value = oid.objectNumber
      setError('')
    } catch (error) {
      setError('Invalid GSF OID: ' + error)
    }
  })

  ;['oid-class', 'oid-type', 'oid-server', 'oid-number'].forEach(id => {
    $(id).addEventListener('input', () => {
      try {
        const oid = {
          objectClass:  $('oid-class').value  || '0',
          objectType:   $('oid-type').value   || '0',
          server:       $('oid-server').value || '0',
          objectNumber: $('oid-number').value || '0',
        }
        const long = oidToLong(oid)
        $('oid-gsf').value = long
        $('oid-cdn').value = btoa(long.toString())
        setError('')
      } catch (error) {
        setError(error)
      }
    })
  })

  // disable mdBook arrow key navigation for input fields
  document.querySelectorAll("input").forEach(el => {
    el.addEventListener("keydown", e => {
      if (["ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown"].includes(e.key)) {
        e.stopPropagation()
      }
    })
  })
</script>
