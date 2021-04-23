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
      <h1>SELECT THE PLANT ON THIS TIMELAPSE</h1>
      <div :class='$style.plants'>
        <div :class='$style.plant' v-for='plant in plants'>
          <div :class='$style.plantinfos'>
            <div :class='$style.plantid'>
              <span :class='$style.plantname'>{{ plant.name }}</span>&nbsp;in&nbsp;<span :class='$style.boxname'>{{ plant.box.name }}</span>
            </div>
            <div :class='$style.plantsettings'>
              <div :class='$style.plantsetting'>
                <img src='~/assets/icon_seeds.svg' />
                <div>
                  {{ plant.settings.strain }}<br/>
                  <span :class='$style.thin'>from</span>&nbsp;<b :class='$style.green'>{{ plant.settings.seedBank }}</b>
                </div>
              </div>
              <div :class='$style.plantsetting'>
                <img src='~/assets/icon_phase.svg' />
                <div>
                  <b>Germinated</b>: 34 days ago<br/>
                  <b>Blooming since</b>: 34 days ago
                </div>
              </div>
            </div>
          </div>
          <div :class='$style.checkboxcontainer'>
            <div :class='$style.checkbox'>
              <Checkbox :checked='selectedPlant == plant' @click='selectPlant(plant)' />
            </div>
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

const API_URL='https://api2.supergreenlab.com'

export default {
  components: {Loading, Checkbox,},
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

    this.$data.plants = plants.map((p, i) => {
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
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  justify-content: center
  height: 100vh

#body
  margin-top: 70pt
  padding: 0 5pt
  width: 100%
  max-width: 600pt

#body > h1
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

.plant
  display: flex
  border-bottom: 1pt dashed #ababab

.plantname
  font-weight: bold
  font-size: 1.2em
  color: #454545

.boxname
  font-weight: bold
  font-size: 1.2em
  color: #3bb30b

.plantinfos
  display: flex
  flex: 1
  flex-direction: column

.plantsettings
  display: flex
  width: 100%
  margin: 15pt 0 10pt 0

.plantsetting
  display: flex
  flex: 1
  color: #454545

.plantsetting > img
  padding-right: 5pt

.green
  color: #3bb30b

.thin
  font-weight: 100

.checkboxcontainer
  display: flex
  justify-content: center
  align-items: center

</style>
