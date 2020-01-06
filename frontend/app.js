window.app = new Vue({

  created() {
    // Use defaults with custom icon paths
    this.icon = L.icon({
      ...L.Icon.Default.prototype.options,
      iconUrl: '/asset/leaflet/marker-icon.png',
      iconRetinaUrl: '/asset/leaflet/marker-icon-2x.png',
      shadowUrl: '/asset/leaflet/marker-shadow.png',
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

  methods: {
    initMap() {
      this.map = L.map('map')
        .setView([0, 0], 13)

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

      this.socket = new WebSocket(`${window.location.href.split('#')[0].replace(/^http/, 'ws')}/ws`)
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

      return axios.put(window.location.href.split('#')[0], data)
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
