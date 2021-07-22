<!--
      Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
      Author: Constantin Clauzel <constantin.clauzel@gmail.com>

      This program is free software: you can redistribute it and/or modify
      it under the terms of the GNU General Public License as published by
      the Free Software Foundation, either version 3 of the License, or
      (at your option) any later version.

      This program is distributed in the hope that it will be useful,
      but WITHOUT ANY WARRANTY without even the implied warranty of
      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
      GNU General Public License for more details.

      You should have received a copy of the GNU General Public License
      along with this program.  If not, see <http://www.gnu.org/licenses/>.
 -->

<template>
  <section :id="$style.container">
    <div :id='$style.body'>
      <h1>SELECT <span :class='$style.green'>THE PLANT</span> ON THIS TIMELAPSE</h1>
      <div :id='$style.plants'>
        <div v-if='plants.length == 0' :id='$style.noplants'>
          <img src='~assets/icon_noplant.svg' />
          <h2>No plant yet</h2>
          <small>Head to the <a href='https://www.supergreenlab.com/app' target='_blank'>app</a> to create one</small>
        </div>
        <div v-else :class='$style.dashed' v-for='plant in plants'>
          <Plant :plant='plant' />
          <div :id='$style.timelapses' v-if='plant.timelapses.length'>
            <h3>Current timelapses</h3>
            <div :class='$style.timelapse' v-for='timelapse in plant.timelapses'>
              <div :class='$style.timelapseinfos'>
                <b>{{ timelapse.name }}</b> - {{ timelapse.nFrames || 0 }} frames<br />
                <span :class='$style.date'>Started: {{ new Date(timelapse.cat).toLocaleString() }}</span>
              </div>
              <div :class='$style.buttonscontainer'>
                <a @click='start(plant, timelapse)'><img src='~/assets/icon_continue.svg' height='20pt' />&nbsp;Continue timelapse</a>
              </div>
            </div>
          </div>
          <div :class='$style.newtimelapse'>
            <a @click='start(plant)'><img src='~/assets/icon_add.svg' height='20pt' />&nbsp;Create new timelapse</a>
          </div>
        </div>
      </div>
    </div>
    <div v-if='loading' :id='$style.loading'>
      <div :id='$style.loadingbox'>
        <Loading />
      </div>
    </div>
  </section>
</template>

<script>
import axios from 'axios'

import Loading from '~/components/loading.vue'
import Checkbox from '~/components/checkbox.vue'
import Plant from '~/components/plant.vue'

const RPI_URL=process.env.RPI_URL

export default {
  components: {Loading, Checkbox, Plant,},
  data() {
    return {
      loading: true,
      plants: [],
    }
  },
  async mounted() {
    const { data: { plants } } = await axios.get(`${RPI_URL}/api/plants?offset=0&limit=100`)
    const { data: { boxes } } = await axios.get(`${RPI_URL}/api/boxes?offset=0&limit=100`)
    const { data: { timelapses } } = await axios.get(`${RPI_URL}/api/timelapses?offset=0&limit=100`)

    this.$data.plants = plants.filter(p => !p.archived && !p.deleted).map((p, i) => {
      p = Object.assign({}, p, {
        box: boxes.find(b => b.id == p.boxID),
        timelapses: timelapses.filter(t => !t.deleted && t.plantID == p.id),
      })
      p.settings = JSON.parse(p.settings)
      if (typeof p.box.settings == 'string') {
        p.box.settings = JSON.parse(p.box.settings)
      }
      return p
    })
    this.$data.loading = false
  },
  methods: {
    async start(plant, timelapse) {
      this.$data.loading = true
      let timelapseID
      if (timelapse) {
        timelapseID = timelapse.id
      } else {
        const name = prompt('Please name this timelapse:', 'Timelapse')
        const { data: { id } } = await axios.post(`${RPI_URL}/api/timelapse`, {
          name: name,
          plantID: plant.id,
          type: 'sglstorage',
          settings: JSON.stringify({}),
        })
        timelapseID = id
      }
      await axios.post(`${RPI_URL}/timelapse`, {
        id: timelapseID,
        plantID: plant.id,
        cron: '@every 10m',
      })
      this.showCamera = true
      this.$store.commit('plant/setPlant', plant)
    },
  },
  watch: {
    plant(val) {
      if (val) {
        console.log(this.showCamera)
        this.$router.replace(this.showCamera ? '/camera' : '/')
      }
    }
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
  padding: 0 5pt
  width: 100%
  max-width: 600pt

#body > h1
  margin: 20pt 0
  color: #454545
  @media only screen and (max-width: 900pt)
    font-size: 1.6em

#loading
  position: absolute
  display: flex
  align-items: center
  justify-content: center
  top: 0; right: 0; bottom: 0; left: 0;
  background-color: rgba(255, 255, 255, 0.5)

#loadingbox
  position: relative
  display: flex
  width: 100pt
  height: 50pt

#plants
  flex: 1 
  overflow: auto

#noplants
  display: flex
  flex-direction: column
  height: 100%
  align-items: center
  justify-content: center
  text-transform: uppercase
  color: #454545

#noplants a
  color: #454545

.dashed
  border-bottom: 1pt dashed #ababab

#timelapses
  margin: 0 20pt

#timelapses h3
  margin: 5pt 0
  color: #454545

.timelapse
  display: flex
  margin: 5pt 0
  border-bottom: 1pt dashed #ababab

.timelapseinfos
  flex: 1
  color: #454545

.timelapse b
  font-weight: bold
  color: #3bb30b

.timelapse .date
  font-size: 0.8em
  color: #898989

.newtimelapse a
  display: flex
  justify-content: center
  align-items: center
  text-align: center
  color: #454545
  cursor: pointer
  margin: 10pt 0

.newtimelapse:hover
  text-decoration: underline

.buttonscontainer a
  display: flex
  justify-content: center
  align-items: center
  margin: 0 10pt
  text-align: center
  cursor: pointer
  color: #454545

.buttonscontainer:hover
  text-decoration: underline

#button
  display: flex
  justify-content: flex-end
  align-items: flex-end
  padding: 15pt 0 15pt 0
  margin: 20pt 0

#button > button
  border: none
  color: white
  border-radius: 2.5px
  background-color: #3bb30b
  padding: 2pt 20pt
  cursor: pointer

#button > button:hover
  background-color: #4bc31b

#button > button:active
  background-color: #2ba300

.disabled
  background-color: #ababab !important

</style>
