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
        <div :class='$style.plant' v-for='plant in plants'>
          <Plant :plant='plant' />
          <div :class='$style.checkboxcontainer'>
            <div :class='$style.checkbox'>
              <Checkbox :checked='selectedPlant == plant' @click='selectPlant(plant)' />
            </div>
          </div>
        </div>
      </div>
      <div :id='$style.button'>
        <button @click='selectedPlant != null ? nextHandler() : null' :class='selectedPlant == null ? $style.disabled : ""'>NEXT</button>
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

const API_URL=process.env.API_URL
const RPI_URL=process.env.RPI_URL

export default {
  components: {Loading, Checkbox, Plant,},
  data() {
    return {
      loading: true,
      plants: [],
      selectedPlant: null,
    }
  },
  async mounted() {
    const { data: { plants } } = await axios.get(`${API_URL}/plants?offset=0&limit=100`, {
      headers: {
        'Authorization': `Bearer ${this.$store.state.auth.token}`,
      },
    })
    const { data: { boxes } } = await axios.get(`${API_URL}/boxes?offset=0&limit=100`, {
      headers: {
        'Authorization': `Bearer ${this.$store.state.auth.token}`,
      },
    })

    this.$data.plants = plants.filter(p => !p.archived && !p.deleted).map((p, i) => {
      p = Object.assign({}, p, {box: boxes.find(b => b.id == p.boxID)})
      p.settings = JSON.parse(p.settings)
      if (typeof p.box.settings == 'string') {
        p.box.settings = JSON.parse(p.box.settings)
      }
      return p
    })
    this.$data.loading = false
  },
  methods: {
    selectPlant(plant) {
      this.$data.selectedPlant = plant
    },
    async nextHandler() {
      this.$data.loading = true
      this.$store.commit('plant/setPlant', this.$data.selectedPlant)
      const token = this.$store.state.auth.token
      const { data: { id } } = await axios.post(`${API_URL}/timelapse`, {
        plantID: this.$data.selectedPlant.id,
        type: 'sglstorage',
        settings: JSON.stringify({}),
      }, { headers: {Authorization: `Bearer ${token}`}})
      await axios.post(`${RPI_URL}/timelapse`, {
        id,
        plantID: this.$data.selectedPlant.id,
        cron: '@every 10m',
      })

      this.$router.push('/camera')
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

.plant
  display: flex
  border-bottom: 1pt dashed #ababab

.checkboxcontainer
  display: flex
  justify-content: center
  align-items: center
  margin: 0 10pt

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
