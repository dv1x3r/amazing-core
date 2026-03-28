# Cache Archive

Amazing World used Unity Streaming Assets, meaning things like images, audio, and asset bundles were stored on official servers and loaded by the game whenever they were needed.
These files were also cached locally. Below is a list of known assets that have been found and shared by members of the community.

<br>

<script src="//unpkg.com/alpinejs" defer></script>
<script src="//unpkg.com/jszip@3.10.1/dist/jszip.min.js"></script>

<script type="importmap">
{
  "imports": {
    "three": "https://cdn.jsdelivr.net/npm/three@0.183.2/build/three.module.js",
    "three/addons/": "https://cdn.jsdelivr.net/npm/three@0.183.2/examples/jsm/"
  }
}
</script>

<script type="module">
  import * as THREE from 'three';
  import { OBJLoader } from 'three/addons/loaders/OBJLoader.js';
  import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

  window.THREE = THREE;
  window.OBJLoader = OBJLoader;
  window.OrbitControls = OrbitControls;
</script>

<script>
  function cacheList() {
    return {
      base: '{{#include ../vars/cache-url.md}}',
      loading: true,
      error: null,

      _items: [],      // all items from cache.json
      _filtered: [],   // cached filtered+sorted list of items
      _selected: null, // index into items[]

      search: '',
      view: 'list',      // 'list' or 'detail'
      sortCol: 'name',   // 'name', 'type', 'asset', 'size'
      sortDir: 'desc',
      page: 1,
      perPage: 25,
      detailTab: 'info',

      async init() {
        try {
          const res = await fetch(`${this.base}/cache.json`)
          if (!res.ok) throw new Error('HTTP ' + res.status)
          this._items = await res.json()
          this._refilter()
        } catch (e) {
          this.error = 'Failed to load cache.json: ' + e.toString()
        } finally {
          this.loading = false
        }
      },

      _refilter() {
        const q = this.search.toLowerCase()
        let list = this._items.map((item, idx) => ({item, idx}))
        if (q) {
          list = list.filter(({item}) => {
            const f = item.file
            if (f.name.toLowerCase().includes(q)) return true
            if (f.type.toLowerCase().includes(q)) return true
            if (String(f.oid).toLowerCase().includes(q)) return true
            if (this.shortestAsset(item).name.toLowerCase().includes(q)) return true
            return false
          })
        }
        const col = this.sortCol
        const dir = this.sortDir === 'asc' ? 1 : -1
        list.sort((a, b) => {
          let va, vb
          if (col === 'name') { va = a.item.file.name.toLowerCase(); vb = b.item.file.name.toLowerCase(); return va < vb ? -dir : va > vb ? dir : 0 }
          if (col === 'type') { va = a.item.file.type.toLowerCase(); vb = b.item.file.type.toLowerCase(); return va < vb ? -dir : va > vb ? dir : 0 }
          if (col === 'asset') { va = this.shortestAsset(a.item).name.toLowerCase(); vb = this.shortestAsset(b.item).name.toLowerCase(); return va < vb ? -dir : va > vb ? dir : 0 }
          if (col === 'size') { return (a.item.file.size - b.item.file.size) * dir }
          return 0
        })
        this._filtered = list
      },

      formatSize(bytes) {
        if (bytes >= 1048576) return (bytes / 1048576).toFixed(1) + ' MB'
        if (bytes >= 1024) return (bytes / 1024).toFixed(1) + ' KB'
        return bytes + ' B'
      },

      shortestAsset(item) {
        const assets = (item.bundle && item.bundle.assets) || []
        if (assets.length === 0) {
          return {name: '', platform: ''}
        }
        let shortest = assets[0]
        for (let i = 1; i < assets.length; i++) {
          if (assets[i].name.length < shortest.name.length) {
            shortest = assets[i]
          }
        }
        return {name: shortest.name, platform: shortest.target_platform || ''}
      },

      selectItem(idx) {
        this._selected = idx
        this.view = 'detail'
        this.detailTab = 'info'
      },

      goBack() {
        this.view = 'list'
        this._selected = null
      },

      toggleListSort(col) {
        if (this.sortCol === col) {
          this.sortDir = this.sortDir === 'asc' ? 'desc' : 'asc'
        } else {
          this.sortCol = col
          this.sortDir = 'asc'
        }
        this.page = 1
        this._refilter()
      },

      listSortIndicator(col) {
        if (this.sortCol !== col) return ''
        return this.sortDir === 'asc' ? ' ↑' : ' ↓'
      },

      get listPaged() {
        const start = (this.page - 1) * this.perPage
        return this._filtered.slice(start, start + this.perPage)
      },

      get listTotalPages() {
        return Math.max(1, Math.ceil(this._filtered.length / this.perPage))
      },

      get listPageRange() {
        const range = []
        for (let i = Math.max(1, this.page - 4); i <= Math.min(this.listTotalPages, this.page + 4); i++) {
          range.push(i)
        }
        return range
      },

      get current() {
        return this._selected !== null ? this._items[this._selected] : null
      },

      get detailRows() {
        const item = this.current
        if (!item) return []
        const b = item.bundle
        const f = item.file

        if (this.detailTab === 'info') {
          const rows = [
            {key: 'File Name', value: f.name},
            {key: 'File Size', value: this.formatSize(f.size)},
            {key: 'File Type', value: f.type},
            {key: 'Hash', value: f.hash},
            {key: 'OID', value: f.oid},
          ]
          if (b && b.info) {
            rows.push(
              {key: 'Signature', value: b.info.signature},
              {key: 'Bundle Version', value: b.info.version},
              {key: 'Player Version', value: b.info.version_player},
              {key: 'Engine Version', value: b.info.version_engine},
            )
          }
          if (b && b.counts) {
            rows.push(
              {key: 'Total Assets', value: b.counts.assets},
              {key: 'Total Objects', value: b.counts.objects},
              {key: 'Total Containers', value: b.counts.container},
            )
          }
          if (b && b.assets) {
            b.assets.forEach((a, idx) => {
              rows.push({key: 'Asset ' + (idx + 1), value: a.name + ' (' + a.target_platform + ')'})
            })
          }
          return rows
        }

        if (this.detailTab === 'counts') {
          if (!b || !b.counts || !b.counts.types) return []
          return Object.entries(b.counts.types)
            .map(([type, count]) => ({ type, count }))
            .sort((a, b) => a.type.localeCompare(b.type))
        }

        if (this.detailTab === 'containers') {
          if (!b || !b.containers) return []
          return Object.entries(b.containers).map(([id, path]) => ({id, path}))
        }

        return []
      },
    }
  }

  function zipLoader({url} = {}) {
    return {
      loading: true,
      error: null,

      zip: null,
      files: {
        audio: [],
        images: [],
        models: [],
      },

      async init() {
        try {
          const res = await fetch(url)
          if (!res.ok) throw new Error('HTTP ' + res.status)
          const blob = await res.blob()
          this.zip = await JSZip.loadAsync(blob)
          const resolve = async (filePath) => {
            const file = this.zip.file(filePath)
            const b = await file.async('blob')
            return {
              path: filePath,
              name: filePath.split('/').pop(),
              url: URL.createObjectURL(b),
            }
          }
          const keys = Object.keys(this.zip.files).filter(f => !f.endsWith('/'))
          const [audio, images, models] = await Promise.all([
            Promise.all(keys.filter(f => f.includes('/audio/')).map(resolve)),
            Promise.all(keys.filter(f => f.includes('/images/')).map(resolve)),
            Promise.all(keys.filter(f => f.includes('/models/')).map(resolve)),
          ])
          this.files = {
            audio: audio.sort((a, b) => a.name.localeCompare(b.name)),
            images: images.sort((a, b) => a.name.localeCompare(b.name)),
            models: models.sort((a, b) => a.name.localeCompare(b.name)),
          }
        } catch (e) {
          this.error = 'Failed to load zip: ' + e.toString()
        } finally {
          this.loading = false
        }
      },

      destroy() {
        const all = [...this.files.audio, ...this.files.images, ...this.files.models]
        all.forEach(f => URL.revokeObjectURL(f.url))
      },
    }
  }

  function fileLoader({url, type} = {}) {
    return {
      loading: true,
      error: null,
      content: null,

      async init() {
        try {
          const res = await fetch(url)
          if (!res.ok) throw new Error('HTTP ' + res.status)
          if (type === 'json') {
            this.content = await res.json()
          } else {
            this.content = await res.text()
          }
        } catch (e) {
          this.error = 'Failed to load file: ' + e.toString()
        } finally {
          this.loading = false
        }
      },
    }
  }

  function audioPlayer({url} = {}) {
    return {
      loading: true,
      playing: false,
      error: null,
      duration: '-',

      _buf: null,
      _src: null,
      _actx: null,

      async init() {
        try {
          const res = await fetch(url)
          if (!res.ok) throw new Error('HTTP ' + res.status)
          this._actx = new (window.AudioContext || window.webkitAudioContext)()
          this._buf = await this._actx.decodeAudioData(await res.arrayBuffer())
          const d = this._buf.duration
          this.duration = Math.floor(d / 60) + ':' + String(Math.floor(d % 60)).padStart(2, '0')
        } catch (e) {
          this.error = 'Failed to load audio: ' + e.toString()
        } finally {
          this.loading = false
        }
      },

      play() {
        if (this.playing) {
          this._src.stop()
          this.playing = false
        } else {
          if (this._actx.state === 'suspended') this._actx.resume()
          this._src = this._actx.createBufferSource()
          this._src.buffer = this._buf
          this._src.connect(this._actx.destination)
          this._src.onended = () => { this.playing = false }
          this._src.start(0)
          this.playing = true
        }
      },

      destroy() {
        if (this._src) this._src.stop()
        if (this._actx) this._actx.close()
      },
    }
  }

  function imageCanvas({url} = {}) {
    return {
      loading: true,
      error: null,

      draw(canvas) {
        const ctx = canvas.getContext('2d')
        const img = new Image()
        img.crossOrigin = 'anonymous'
        img.onload = () => {
          canvas.width = img.naturalWidth
          canvas.height = img.naturalHeight
          canvas.style.maxWidth = '100%'
          canvas.style.height = 'auto'
          ctx.drawImage(img, 0, 0)
          this.loading = false
        };
        img.onerror = (e) => {
          this.error = 'Failed to load image'
          this.loading = false
        };
        img.src = url
      }
    }
  }

  function threeViewer() {
    return {
      _scene: null,
      _renderer: null,
      _controls: null,
      _resizeObserver: null,

      init() {
        const THREE = window.THREE
        const container = this.$el

        const w = container.clientWidth || 400
        const h = container.clientHeight || 400

        const scene = new THREE.Scene()
        scene.background = new THREE.Color(0x333333)
        this._scene = scene

        var grid = new THREE.GridHelper(100, 100);
        scene.add(grid);

        const camera = new THREE.PerspectiveCamera(75, w / h, 0.1, 1000)
        camera.position.y = 0.8
        camera.position.z = 1

        const renderer = new THREE.WebGLRenderer({ antialias: true })
        renderer.setPixelRatio(window.devicePixelRatio)
        renderer.setSize(w, h)
        this._renderer = renderer

        container.appendChild(renderer.domElement)

        const controls = new OrbitControls(camera, renderer.domElement)
        controls.enableDamping = true
        //controls.autoRotate = true
        //controls.autoRotateSpeed = 1
        //controls.listenToKeyEvents(window)
        this._controls = controls

        //const geometry = new THREE.BoxGeometry(1, 1, 1)
        //const material = new THREE.MeshBasicMaterial({ color: 0xffffff, wireframe: true })
        //const cube = new THREE.Mesh(geometry, material)
        //scene.add(cube)

        const light = new THREE.AmbientLight(0xffffff, 0.8)
        scene.add(light)

        const resizeObserver = new ResizeObserver(() => {
          const w = container.clientWidth
          const h = container.clientHeight
          if (w !== 0 && h !== 0) {
            renderer.setSize(w, h)
            camera.aspect = w / h
            camera.updateProjectionMatrix()
          }
        })

        resizeObserver.observe(container);
        this._resizeObserver = resizeObserver

        function render(time) {
          controls.update()
          renderer.render(scene, camera)
        }

        renderer.setAnimationLoop(render)
        this.$watch('files', async () => await this.load())
      },

      destroy() {
        this._resizeObserver?.disconnect()
        this._controls?.dispose()
        this._renderer?.setAnimationLoop(null)
        this._renderer?.domElement.remove()
        this._renderer?.dispose()
      },

      async load() {
        const meshMap = {}
        for (const model of this.files.models) {
          const id = parseInt(model.name.split('.')[0])
          meshMap[id] = model.url
        }

        const textureMap = {}
        for (const image of this.files.images) {
          const id = parseInt(image.name.split('.')[0])
          textureMap[id] = image.url
        }

        const textureNameMap = {}
        for (const image of this.files.images) {
          const name = image.name.toLowerCase()
          textureNameMap[name] = image.url
        }

        //console.log('models', meshMap)
        //console.log('images', textureMap)
        //console.log('scene', this.content.bundle.scene)

        const results = []

        const walkNode = (node, parentMatrix, rootsLength) => {
          //console.log('walkNode', node.name)
          //console.log('rootsLength', rootsLength)

          let worldMatrix = parentMatrix
          if (node.transform) {
            const tr = node.transform
            // some root objects are not centered, so we reset the position if scene has only one object
            const isSingleRootNode = rootsLength !== undefined && rootsLength === 1
            const position = isSingleRootNode
              ? new THREE.Vector3(0, 0, 0)
              : new THREE.Vector3(tr.position.x, tr.position.y, -tr.position.z)
            const rotation = new THREE.Quaternion(tr.rotation.x, -tr.rotation.y, -tr.rotation.z, tr.rotation.w)
            const scale = new THREE.Vector3(tr.scale.x, tr.scale.y, tr.scale.z)
            const localMatrix = new THREE.Matrix4().compose(position, rotation, scale)
            worldMatrix = parentMatrix ? parentMatrix.clone().multiply(localMatrix) : localMatrix
          }

          for (const comp of node.components ?? []) {
            if (comp.type === 'MeshFilter' || comp.type === 'SkinnedMeshRenderer') {
              const mesh_url = meshMap[comp.mesh?.id]
              if (!mesh_url) {
                continue
              }
              const materials = comp.materials ?? node.components.find(c => c.type === 'MeshRenderer')?.materials ?? []
              const texture_urls = materials.map(mat => {
                const textureId = mat.textures?.[0]?.id
                return textureId ? textureMap[textureId] : null
              })
              results.push({ name: node.name, mesh_url, texture_urls, worldMatrix })
            }
          }

          for (const child of node.children ?? []) {
            walkNode(child, worldMatrix)
          }
        }

        for (const root of this.content.bundle.scene) {
          walkNode(root, new THREE.Matrix4(), this.content.bundle.scene.length)
        }

        for (const result of results) {
          const loader = new OBJLoader()
          const object = await loader.loadAsync(result.mesh_url)
          const textures = await Promise.all(
            result.texture_urls.map(url => url ? new THREE.TextureLoader().loadAsync(url) : Promise.resolve(null))
          )

          let meshIndex = 0
          object.traverse(child => {
            if (!child.isMesh) {
              return
            }
            const texture = textures[meshIndex++]
            child.material = new THREE.MeshBasicMaterial({
              map: texture ?? null,
              transparent: true,
              side: THREE.DoubleSide,
            })
            if (texture) {
              texture.wrapS = THREE.RepeatWrapping
              texture.wrapT = THREE.RepeatWrapping
            }
          })

          if (result.worldMatrix) {
            object.applyMatrix4(result.worldMatrix)
          }

          this._scene.add(object)
        }

      },
    }
  }
</script>

<div id="cache-archive" x-data="cacheList()" @keydown.stop @keyup.stop @keypress.stop>
  <p x-show="loading">Loading cache list...</p>
  <p x-show="error" x-text="error" style="color:red"></p>
  <div x-show="!loading && !error">
    <!-- List View -->
    <template x-if="view === 'list'">
      <div>
        <div style="margin-bottom:8px;">
          <input type="text" x-model="search"
            @input="page=1, _refilter()"
            placeholder="Search by name, type, oid or asset…"
            style="font-family:monospace; padding:3px 6px; width:300px;">
          <span style="margin-left:10px; font-size:1.25rem;" x-text="_filtered.length + ' items'"></span>
        </div>
        <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
          <thead>
            <tr>
              <th style="width:25%; text-align:left; padding:4px 8px; cursor:pointer;"
                @click="toggleListSort('name')">Name<span x-text="listSortIndicator('name')"></span>
              </th>
              <th style="width:20%; text-align:left; padding:4px 8px; cursor:pointer;"
                @click="toggleListSort('type')">Type<span x-text="listSortIndicator('type')"></span>
              </th>
              <th style="width:45%; text-align:left; padding:4px 8px; cursor:pointer;"
                @click="toggleListSort('asset')">Asset<span x-text="listSortIndicator('asset')"></span>
              </th>
              <th style="width:10%; text-align:right; padding:4px 8px; cursor:pointer;"
                @click="toggleListSort('size')">Size<span x-text="listSortIndicator('size')"></span>
              </th>
            </tr>
          </thead>
          <tbody>
            <template x-for="entry in listPaged" :key="entry.idx">
              <tr @click="selectItem(entry.idx)" style="cursor:pointer;"
                  @mouseenter="$el.style.background='var(--table-header-bg)'" @mouseleave="$el.style.background=''">
                <td style="padding:3px 8px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;" x-text="entry.item.file.name"></td>
                <td style="padding:3px 8px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;" x-text="entry.item.file.type"></td>
                <td style="padding:3px 8px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis;">
                  <span x-text="shortestAsset(entry.item).name"></span>
                  <span style="color:#888; font-size:0.85em;"
                    x-show="shortestAsset(entry.item).platform"
                    x-text="' (' + shortestAsset(entry.item).platform + ')'"></span>
                </td>
                <td style="padding:3px 8px; text-align:right; white-space:nowrap;" x-text="formatSize(entry.item.file.size)"></td>
              </tr>
            </template>
            <tr x-show="listPaged.length === 0">
              <td colspan="4" style="padding:8px; color:#888;">No results.</td>
            </tr>
          </tbody>
        </table>
        <!-- List Pagination -->
        <div x-show="listTotalPages > 1" style="font-size:1.25rem; margin-top:10px;">
          <button @click="page--" :disabled="page === 1"
            :style="page > 1 ? { cursor: 'pointer' } : {}"
            style="font-family:monospace; margin-right:4px;">&#8249; Prev</button>
          <template x-for="p in listPageRange" :key="p">
            <button @click="page = p" x-text="p"
              :style="p === page ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
              style="font-family:monospace; margin-right:4px; cursor:pointer;"></button>
          </template>
          <button @click="page++" :disabled="page === listTotalPages"
            :style="page < listTotalPages ? { cursor: 'pointer' } : {}"
            style="font-family:monospace;">Next &#8250;</button>
          <span style="margin-left:10px;" x-text="'Page ' + page + ' of ' + listTotalPages"></span>
        </div>
      </div>
    </template>
    <!-- Detail View -->
    <template x-if="view === 'detail' && current">
      <div>
        <!-- Tabs -->
        <div style="font-size:1.25rem; margin-bottom:10px;">
          <button @click="goBack()"
            style="font-family:monospace; margin-right:4px; cursor:pointer;">&#8592; Back</button>
          <button @click="detailTab = 'info'"
            x-show="current.file.type.startsWith('AssetBundle/')"
            :style="detailTab === 'info' ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
            style="font-family:monospace; margin-right:4px; cursor:pointer;">Info</button>
          <button @click="detailTab = 'counts'"
            x-show="current.file.type.startsWith('AssetBundle/')"
            :style="detailTab === 'counts' ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
            style="font-family:monospace; margin-right:4px; cursor:pointer;">Object Counts</button>
          <button @click="detailTab = 'containers'"
            x-show="current.file.type.startsWith('AssetBundle/')"
            :style="detailTab === 'containers' ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
            style="font-family:monospace; margin-right:4px; cursor:pointer;">Containers</button>
        </div>
        <!-- Info tab -->
        <template x-if="detailTab === 'info'">
          <div>
            <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
              <thead>
                <tr>
                  <th style="width:25%; text-align:left;">Field</th>
                  <th style="width:75%; text-align:left;">Value</th>
                </tr>
              </thead>
              <tbody>
                <template x-for="row in detailRows">
                  <tr>
                    <td x-text="row.key"></td>
                    <td x-text="row.value"></td>
                  </tr>
                </template>
                <tr x-show="detailRows.length === 0">
                  <td colspan="2" style="padding:8px; color:#888;">No results.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
        <!-- Object Counts tab -->
        <template x-if="detailTab === 'counts'">
          <div>
            <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
              <thead>
                <tr>
                  <th style="width:25%; text-align:left;">Type</th>
                  <th style="width:75%; text-align:left;">Count</th>
                </tr>
              </thead>
              <tbody>
                <template x-for="row in detailRows">
                  <tr>
                    <td x-text="row.type"></td>
                    <td x-text="row.count"></td>
                  </tr>
                </template>
                <tr x-show="detailRows.length === 0">
                  <td colspan="2" style="padding:8px; color:#888;">No results.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
        <!-- Containers tab -->
        <template x-if="detailTab === 'containers'">
          <div>
            <table style="width:100%; font-size:1.25rem;">
              <thead>
                <tr>
                  <th style="text-align:left;">Path</th>
                </tr>
              </thead>
              <tbody>
                <template x-for="row in detailRows">
                  <tr>
                    <td x-text="row.path"></td>
                  </tr>
                </template>
                <tr x-show="detailRows.length === 0">
                  <td colspan="2" style="padding:8px; color:#888;">No results.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
        <div>
          <h2>Preview</h2>
          <!-- Preview images -->
          <template x-if="current.file.type.startsWith('image/')">
            <div x-data="imageCanvas({url: base + '/cache/' + current.file.name})">
              <p x-show="loading">Loading image...</p>
              <p x-show="error" x-text="error" style="color:red"></p>
              <canvas x-show="!loading && !error" x-init="draw($el)"></canvas>
            </div>
          </template>
          <!-- Preview audio -->
          <template x-if="current.file.type.startsWith('audio/')">
            <div x-data="audioPlayer({url: base + '/cache/' + current.file.name})">
              <p x-show="loading">Loading audio...</p>
              <p x-show="error" x-text="error" style="color:red"></p>
              <div x-show="!loading && !error">
                <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
                  <thead>
                    <tr>
                      <th style="width:50%; text-align:left;">File</th>
                      <th style="width:25%; text-align:left;">Duration</th>
                      <th style="width:25%; text-align:left;">Controls</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td x-text="current.file.name"></td>
                      <td x-text="duration"></td>
                      <td>
                        <button @click="play()"
                          x-text="playing ? 'Stop' : 'Play'"
                          style="font-family:monospace; cursor:pointer;"></button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </template>
          <!-- Preview text -->
          <template x-if="current.file.type.startsWith('TreeNode/') || current.file.type === 'json'">
            <div x-data="fileLoader({url: base + '/cache/' + current.file.name})">
              <p x-show="loading">Loading file...</p>
              <p x-show="error" x-text="error" style="color:red"></p>
              <div x-show="content && !loading">
                <pre><code class="hljs" style="font-size:1.25rem; max-height:400px;" x-text="content"></code></pre>
              </div>
            </div>
          </template>
          <!-- Preview asset bundle -->
          <template x-if="current.file.type.startsWith('AssetBundle/')">
            <div x-data="fileLoader({url: base + '/unpacked/' + current.file.name + '.json', type: 'json'})">
              <p x-show="loading">Loading summary...</p>
              <p x-show="error" x-text="error" style="color:red"></p>
              <div x-data="{ bundleTab: 'summary' }" x-show="!loading && !error">
                <!-- Tabs -->
                <div style="margin-bottom:10px;">
                  <button @click="bundleTab = 'summary'"
                    x-show="current.file.type.startsWith('AssetBundle/')"
                    :style="bundleTab === 'summary' ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
                    style="font-family:monospace; font-size:1.25rem; margin-right:4px; cursor:pointer;">Summary</button>
                  <button @click="bundleTab = 'files'"
                    x-show="content?.bundle.counts.types.Mesh > 0 || content?.bundle.counts.types.AudioClip > 0 || content?.bundle.counts.types.Texture2D > 0"
                    :style="bundleTab === 'files' ? { fontWeight: 'bold', textDecoration: 'underline' } : {}"
                    style="font-family:monospace; font-size:1.25rem; margin-right:4px; cursor:pointer;">Files</button>
                </div>
                <!-- Summary -->
                <template x-if="bundleTab === 'summary'">
                  <pre><code class="hljs" style="font-size:1.25rem; max-height:400px;" x-text="JSON.stringify(content, null, 2)"></code></pre>
                </template>
                <!-- Files -->
                <template x-if="bundleTab === 'files'">
                  <div x-data="zipLoader({url: base + '/unpacked/' + current.file.name + '.zip'})">
                    <p x-show="loading">Loading archive...</p>
                    <p x-show="error" x-text="error" style="color:red"></p>
                    <div x-show="!loading && !error">
                      <!-- Scene -->
                      <div x-show="content?.bundle.scene.length > 0 && content?.bundle.counts.types.Mesh > 0">
                        <h3>Scene</h3>
                        <div x-data="threeViewer()"></div>
                        <p style="font-size:1.25rem;">
                          Drag to rotate, right-click drag to pan, scroll to zoom. Some assets may not be centered.
                        </p>
                      </div>
                      <!-- Audio files -->
                      <div x-show="files.audio.length > 0" style="margin-top:16px;">
                        <h3>Audio</h3>
                        <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
                          <thead>
                            <tr>
                              <th style="width:50%; text-align:left;">File</th>
                              <th style="width:25%; text-align:left;">Duration</th>
                              <th style="width:25%; text-align:left;">Controls</th>
                            </tr>
                          </thead>
                          <tbody>
                            <template x-for="aud in files.audio">
                              <tr x-data="audioPlayer({url: aud.url})">
                                <td x-text="aud.name"></td>
                                <td x-text="duration"></td>
                                <td>
                                  <span x-show="loading">Loading...</span>
                                  <span x-show="error" x-text="error" style="color:red;"></span>
                                  <template x-if="!loading && !error">
                                    <button @click="play()"
                                      x-text="playing ? 'Stop' : 'Play'"
                                      style="font-family:monospace; cursor:pointer;"></button>
                                  </template>
                                </td>
                              </tr>
                            </template>
                          </tbody>
                        </table>
                      </div>
                      <!-- Image files -->
                      <div x-show="files.images.length > 0" style="margin-top:16px;">
                        <h3>Images</h3>
                        <table style="width:100%; font-size:1.25rem; table-layout:fixed;">
                          <thead>
                            <tr>
                              <th style="width:50%; text-align:left;">File</th>
                              <th style="width:50%; text-align:left;">Image</th>
                            </tr>
                          </thead>
                          <tbody>
                            <template x-for="img in files.images">
                              <tr x-data="imageCanvas({url: img.url})">
                                <td x-text="img.name"></td>
                                <td>
                                  <span x-show="loading">Loading...</span>
                                  <span x-show="error" x-text="error" style="color:red;"></span>
                                  <canvas x-show="!loading && !error" x-init="draw($el)"></canvas>
                                </td>
                              </tr>
                            </template>
                          </tbody>
                        </table>
                      </div>
                      <!-- Model files -->
                      <div x-show="files.models.length > 0" style="margin-top:16px;">
                        <h3>Models</h3>
                        <ul style="font-size:1.25rem;">
                          <template x-for="mod in files.models">
                            <li x-text="mod.name"></li>
                          </template>
                        </ul>
                      </div>
                    </div>
                  </div>
                </template>
              </div>
            </div>
          </template>
        </div>
      </div>
    </template>
  </div>
</div>

<h2>Notice</h2>

<p>
  This website and its contents are intended strictly to be used for non-commercial, personal, and educational purposes only.
</p>
