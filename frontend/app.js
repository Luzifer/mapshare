/* global VueI18n */

const get_locale = (fallback = 'en') => {
  const urlParams = new URLSearchParams(window.location.search)

  for (const lc of [
    urlParams.get('hl'),
    navigator.languages,
    navigator.language,
    navigator.browserLanguage,
    navigator.userLanguage,
    fallback,
  ]) {
    if (!lc) {
      continue
    }

    switch (typeof lc) {
    case 'object':
      if (lc.length > 0) {
        return lc[0].split('-')[0]
      }
      break
    case 'string':
      return lc.split('-')[0]
    }
  }

  return fallback
}

const i18n = new VueI18n({
  locale: get_locale(),

  messages: {
    de: {
      btnModalOK: 'OK',
      optKeepSending: 'Position kontinuierlich senden',
      optKeepSendingSub: '(wenn aktiviert, wird die Position gesendet, solange dieses Fenster offen ist)',
      optRetainLocation: 'Position auf dem Server speichern',
      optRetainLocationSub: '(neue Beobachter sehen die Position sofort)',
      btnShareMyLocation: 'Meine Position senden!',
      shareSettings: 'Einstellungen',
      waitingForLocation: 'Warte auf Position...',
    },
    en: {
      btnModalOK: 'OK',
      optKeepSending: 'Keep sending location',
      optKeepSendingSub: '(when enabled location is updated as long as this window is open)',
      optRetainLocation: 'Retain location on server',
      optRetainLocationSub: '(new viewers instantly see your location)',
      btnShareMyLocation: 'Share my location!',
      shareSettings: 'Share-Settings',
      waitingForLocation: 'Waiting for location...',
    },
  },
})

window.app = new Vue({

  computed: {
    zoom() {
      const params = new URLSearchParams(window.location.search)
      return params.has('zoom') ? parseInt(params.get('zoom')) : 13
    },
  },

  created() {
    // Use defaults with custom icon paths
    this.icon = L.icon({
      ...L.Icon.Default.prototype.options,
      iconUrl: '/asset/images/marker-icon.png',
      iconRetinaUrl: '/asset/images/marker-icon-2x.png',
      shadowUrl: '/asset/images/marker-shadow.png',
    })

    /*
     * This is only to detect another user updated the location
     * therefore this is NOT cryptographically safe!
     */
    this.browserID = localStorage.getItem('browserID')
    if (!this.browserID) {
      this.browserID = Math.random().toString(16)
        .substr(2)
      localStorage.setItem('browserID', this.browserID)
    }
  },

  data: {
    browserID: null,
    icon: null,
    loc: null,
    map: null,
    marker: null,
    shareActive: false,
    shareSettings: {
      continuous: true,
      retained: false,
    },
    shareSettingsOpen: false,
    socket: null,
  },

  el: '#app',

  i18n,

  methods: {
    initMap() {
      this.map = L.map('map')
        .setView([0, 0], this.zoom)

      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.map)
    },

    shareLocation() {
      this.shareActive = true

      const opts = {
        enableHighAccuracy: true,
        timeout: 5000,
        maximumAge: 0,
      }

      if (!this.shareSettings.continuous) {
        navigator.geolocation.getCurrentPosition(this.updateLocation, err => console.error(err), opts)
        return
      }

      navigator.geolocation.watchPosition(this.updateLocation, err => console.error(err), opts)
    },

    subscribe() {
      if (this.socket) {
        // Dispose old socket
        this.socket.close()
        this.socket = null
      }

      let socketAddr = window.location.href.replace(/^http/, 'ws')
      socketAddr = socketAddr.split('#')[0]
      socketAddr = socketAddr.split('?')[0]
      socketAddr = `${socketAddr}/ws`

      this.socket = new WebSocket(socketAddr)
      this.socket.onclose = () => window.setTimeout(this.subscribe, 1000) // Restart socket
      this.socket.onmessage = evt => {
        const loc = JSON.parse(evt.data)
        loc.time = new Date(loc.time)
        this.loc = loc
      }
    },

    updateLocation(pos) {
      const data = {
        lat: pos.coords.latitude,
        lon: pos.coords.longitude,
        retained: this.shareSettings.retained,
        sender_id: this.browserID,
      }

      return fetch(window.location.href.split('#')[0], {
        body: JSON.stringify(data),
        credentials: 'same-origin',
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'PUT',
      })
        .catch(err => console.error(err))
    },

    updateMap() {
      const center = [this.loc.lat, this.loc.lon]

      if (!this.marker) {
        this.marker = L.marker(center, { icon: this.icon })
          .addTo(this.map)
      }

      this.map.panTo(center)
      this.marker.setLatLng(center)
    },
  },

  mounted() {
    this.initMap()
    this.subscribe()
  },

  watch: {
    loc() {
      this.updateMap()
    },
  },
})
