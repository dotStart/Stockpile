/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import Vue from 'vue'
import SuiVue from 'semantic-ui-vue'
import moment from 'moment'
import io from 'socket.io-client'

const ProfileIdEvent = 0;
const NameHistoryEvent = 1;
const ProfileEvent = 2;
const BlacklistEvent = 3;

const GOLANG_RFC3399 = 'YYYY-MM-DDTHH:mm:ss.SSSSSSSSSZ ZZ';

const opts = {
  path: '/ui/socket.io'
};
const socket = typeof SERVER_ADDR_OVERRIDE === "undefined" ? io(opts)
    : io('http://' + SERVER_ADDR_OVERRIDE, opts);

// import third party components
Vue.use(SuiVue);

// define a bunch of components for our event feed
Vue.component('profileId-event', {
  template: '#cache-event-profileId',
  props: ['event'],
  computed: {
    firstSeenAt: function () {
      return moment(this.event.Object.FirstSeenAt, GOLANG_RFC3399).calendar()
    },
    validUntil: function () {
      return moment(this.event.Object.ValidUntil, GOLANG_RFC3399).calendar()
    }
  }
});

Vue.component('nameHistory-event', {
  template: '#cache-event-nameHistory',
  props: ['event']
});

Vue.component('profile-event', {
  template: '#cache-event-profile',
  props: ['event']
});

Vue.component('blacklist-event', {
  template: '#cache-event-blacklist',
  props: ['event']
});

Vue.component('event-list', {
  template: '#cache-event-list',
  props: ['events']
});

const app = new Vue({
  el: '#stockpile',
  data: {
    connected: false,

    rateLimitAllocation: 0,
    version: '',
    plugins: [],
    pluginsUnavailable: false,

    events: []
  },
  computed: {
    address: function() {
      const addr = location.host;

      if (addr.indexOf(':') !== -1) {
        return addr
      }

      if (location.origin.indexOf('https') !== -1) {
        return addr + ':443'
      }

      return addr + ':80'
    },
    rateLimitLabel: function () {
      return `Rate Limit: ${this.rateLimitAllocation} / 600`
    },
    rateLimitPercent: function () {
      return this.rateLimitAllocation / 600 * 100
    }
  }
});

socket.on('system', (sys) => {
  console.log('Stockpile v' + sys.version);
  if (sys.pluginsSupported) {
    console.log('Loaded plugins: ' + sys.plugins.map(
        plugin => plugin.Name + ' v' + plugin.Version).toString());
    app.plugins = sys.plugins;
  } else {
    console.log('Plugins are not supported by server');
    app.pluginsUnavailable = true
  }

  app.version = sys.version;
});

socket.on('rate-limit', (allocation) => {
  console.log('Current rate limit allocation: ' + allocation);
  app.rateLimitAllocation = allocation
});

socket.on('cache', (data) => {
  console.log('Event: ' + JSON.stringify(data));

  if (app.events.length > 50) {
    app.events.splice(-1, 1)
  }

  app.events.unshift(data);
});
