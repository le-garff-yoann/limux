<template>
  <div id="app">
    <b-navbar fixed="top" type="dark" variant="primary">
      <b-navbar-brand class="font-weight-light">Limux</b-navbar-brand>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav class="ml-auto">
          <b-spinner
            v-show="ws.isConnected"
            type="grow"
            variant="success"

            v-b-popover.hover="'Connected.'"
            data-toggle="tooltip"
          ></b-spinner>
          <b-spinner
            v-show="! ws.isConnected"
            variant="danger"

            v-b-popover.hover="'Trying to connect....'"
          ></b-spinner>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-list-group class="body">
      <b-list-group-item v-for="item in processors" v-bind:key="item.key">
        <b-row align-v="center">
          <b-col>
            <highlight-code lang="yaml" class="shadow p-3 mb-5 bg-white rounded">{{ item.processor | toYaml }}</highlight-code>
          </b-col>
          <b-col>
            <trend
              v-bind:data="item.values"
              :gradient="['lightgray', 'gray', 'black']"
              smooth
              auto-draw

              class="shadow p-3 mb-5 bg-white rounded"
            ></trend>
          </b-col>
        </b-row>
      </b-list-group-item>
    </b-list-group>
  </div>
</template>

<script>
import Vue from 'vue'
import VueNativeSock from 'vue-native-websocket'
Vue.use(VueNativeSock, 'ws://localhost:8080', {
  connectManually: true
})

import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
Vue.use(BootstrapVue)

import Trend from 'vuetrend'
Vue.use(Trend)

import VueHighlightJS from 'vue-highlight.js'
import yaml from 'highlight.js/lib/languages/yaml'
import 'highlight.js/styles/monokai.css'
Vue.use(VueHighlightJS, {
  languages: { yaml }
})

import Jsum from 'jsum'
import JsYaml from 'js-yaml'

export default {
  name: 'app',
  filters: {
    toYaml(v) {
      return JsYaml.safeDump(v)
    }
  },
  data() {
    return {
      processors: [],
      ws: { isConnected: false }
    }
  },
  created() {
    this.$options.sockets.onopen = () => this.ws.isConnected = true
    this.$options.sockets.onclose = () => {
      this.processors = []

      this.ws.isConnected = false
    }

    this.$options.sockets.onmessage = e => this.pushEvent(JSON.parse(e.data))

    setInterval(() => {
      this.ws.isConnected || this.$connect(`ws${window.location.protocol === 'https:' ? 's' : ''}://${window.location.host}${process.env.VUE_APP_ROOT_API}/ws/events`)
    }, 1000)
  },
  methods: {
    pushEvent(o) {
      const k = Jsum.digest(o.processor, 'SHA256', 'hex')

      if (this.processors.findIndex(p => p.key == k) == -1) {
        this.processors.push({
          key: k,
          processor: o.processor,
          values: [0, 0]
        })
      }

      const pr = this.processors.find(p => p.key == k)
      const oldValue = pr.values[pr.values.length - 1]

      if (o.type == 1) {
        pr.values.push(oldValue + 1)
      } else if (o.type == 2) {
        pr.values.push(oldValue - 1)
      }
    }
  }
}
</script>

<style>
html, #app {
  background-color: whitesmoke;
}

#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.body {
  margin-top: 70px;
  margin-right: 10px;
  margin-left: 10px;
}

.trend {
  width: 100%;
  position: absolute;
}
</style>
