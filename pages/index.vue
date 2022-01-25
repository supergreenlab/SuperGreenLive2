<!--
      Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
      Author: Constantin Clauzel <constantin.clauzel@gmail.com>

      This program is free software: you can redistribute it and/or modify
      it under the terms of the GNU General Public License as published by
      the Free Software Foundation, either version 3 of the License, or
      (at your option) any later version.

      This program is distributed in the hope that it will be useful,
      but WITHOUT ANY WARRANTY; without even the implied warranty of
      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
      GNU General Public License for more details.

      You should have received a copy of the GNU General Public License
      along with this program.  If not, see <http://www.gnu.org/licenses/>.
 -->

<template>
  <section :id='$style.container'>
    <div :id='$style.body'>
      <div :id='$style.header'>
        <h1>PLANT ON THIS <span :class='$style.green'>TIMELAPSE</span>:</h1>
        <div :id='$style.buttons'>
          <div><a href='javascript:void(0)' :class='$style.button' @click='showTimes'><img src='~assets/timelapse-times.svg' /><span>Timelapse</span></a></div>
          <div><div :class='$style.button'><div :class='$style.checkbox'><Checkbox @click='toggleSkipNight' :checked='skipNight'/></div><span>Skip Night</span></div></div>
          <div><nuxt-link to='/camera' :class='$style.button'><img src='~assets/icon_livecam.svg' /><span>Live cam</span></nuxt-link></div>
          <div><a :href='storage' target='_blank' :class='$style.button'><img src='~assets/icon_download.svg' /><span>Download</span></a></div>
          <div><a href='javascript:void(0)' :class='$style.button' @click='reset'><img src='~assets/icon_reset.svg' /><span>Reset</span></a></div>
        </div>
      </div>
      <div :id='$style.plantInfos'>
        <Plant :plant='plant' />
      </div>
      <div :id='$style.capture'>
        <div :id='$style.loading'>
          <div><Loading label='Capturing pic..' /></div>
        </div>
        <div v-for='src in srcs' v-if='src' :key='src' :style='{"background-image": `url(${src})`}'></div>
      </div>
    </div>
    <Times v-if='showTimelapseSettings' @close='closeTimes' />
    <Premium />
  </section>
</template>

<script>
import axios from 'axios'
import Loading from '~/components/loading.vue'
import Checkbox from '~/components/checkbox.vue'
import Plant from '~/components/plant.vue'
import Times from '~/components/times.vue'
import Premium from '~/components/premium.vue'

const RPI_URL=process.env.RPI_URL

export default {
  components: {Checkbox, Times, Loading, Plant, Premium,},
  data() {
    return {
      n: 0,
      showTimelapseSettings: false,
      skipNight: null,
      srcs: [null, `${RPI_URL}/capture`],
      storage: `${RPI_URL}/storage.zip`,
    }
  },
  mounted() {
    axios.post(`${RPI_URL}/motion/stop`) // in case the page was reloaded and motion never stopped
    this.interval = setInterval(() => {
      this.$data.srcs = [
        this.$data.srcs[1],
        `${RPI_URL}/capture?rand=${new Date().getTime()}`
      ]
    }, 120000)
    axios.get(`${RPI_URL}/timelapse`).then(({ data: { skipNight } }) => {
      skipNight = skipNight == null ? true : skipNight == 'true'
      this.$data.skipNight = skipNight
    })
  },
  destroyed() {
    clearInterval(this.interval)
  },
  methods: {
    reset() {
      const c = confirm('Start new timelapse? This is not reversible.')
      if (!c) return;
      this.$store.commit('plant/setPlant', null)
      this.$router.push("/plant")
    },
    toggleSkipNight() {
      this.$data.skipNight = !this.$data.skipNight
      axios.post(`${RPI_URL}/timelapse`, {skipNight: `${this.$data.skipNight}`})
    },
    closeTimes() {
      this.$data.showTimelapseSettings = false
    },
    showTimes() {
      this.$data.showTimelapseSettings = true
    },
  },
  computed: {
    plant() {
      return this.$store.state.plant.plant
    },
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  justify-content: center
  height: 100vh

#body
  display: flex
  flex-direction: column
  margin-top: 70pt
  padding: 0 20pt
  width: 100%
  max-width: 600pt
  @media only screen and (max-width: 1000pt)
    margin-top: 60pt

#header
  display: flex
  justify-content: space-between
  align-items: center
  @media only screen and (max-width: 900pt)
    flex-direction: column

#header > h1
  margin: 20pt 0
  color: #454545
  @media only screen and (max-width: 1000pt)
    font-size: 1.2em
    margin: 10pt 0

#buttons
  display: flex
  @media only screen and (max-width: 900pt)
    margin-bottom: 15pt

#buttons > div
  height: 40px
  text-align: center

.button
  display: flex
  flex-direction: column
  justify-content: space-between
  font-weight: 600
  color: #3bb30b
  text-decoration: none
  margin: 0 4pt
  font-size: 0.7em

.button:hover
  text-decoration: underline

.checkbox, .button > img
  height: 35px
  margin-bottom: 4pt

.green
  color: #3bb30b

#capture
  position: relative
  height: 100%
  margin: 30pt 0
  @media only screen and (max-width: 1000pt)
    margin: 15pt 0

#capture > div
  position: absolute
  top: 0
  right: 0
  bottom: 0
  left: 0
  background-position: center
  background-size: contain
  background-repeat: no-repeat

#loading
  display: flex
  align-items: center
  justify-content: center

#loading > div
  position: relative
  width: 100pt
  height: 100pt

#plantInfos
  @media only screen and (max-width: 1000pt)
    font-size: 0.8em

</style>
