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
      <h1>PLACE YOUR <span :class='$style.green'>CAMERA</span>:</h1>
      <div :id='$style.videocontainer'>
        <div :id='$style.motion'>
          <img v-if='motionStarted' :src='src' @error='imgError'/>
        </div>
      </div>
      <div :id='$style.button'>
        <button @click='nextHandler'>NEXT</button>
      </div>
    </div>
  </section>
</template>

<script>
import axios from 'axios'

const RPI_URL='http://192.168.1.26:8080'
const LIVEVIEW_URL='http://192.168.1.26:8081'

export default {
  data() {
    return {
      src: LIVEVIEW_URL,
      motionStarted: false,
    }
  },
  async mounted() {
    await axios.post(`${RPI_URL}/motion/start`)
    setTimeout(() => {
      this.$data.motionStarted = true
    }, 1000)
  },
  methods: {
    async nextHandler() {
      await axios.post(`${RPI_URL}/motion/stop`)
      await axios.post(`${RPI_URL}/timelapse`, {
        cron: '@every 10m',
      })
      this.$router.push("/")
    },
    async imgError() {
      setTimeout(() => {
        this.$data.src = `${LIVEVIEW_URL}?rand=${new Date().getTime()}`
      }, 1000)
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

.green
  color: #3bb30b
  
#videocontainer
  flex: 1
  height: 100%
  display: flex
  flex-direction: column
  align-items: center
  justify-content: center

#motion
  flex: 1
  display: flex
  flex-direction: column
  align-items: center
  justify-content: center

#motion > img
  max-width: 100%
  max-height: 50vh

#button
  display: flex
  justify-content: flex-end
  align-items: flex-end
  padding: 15pt 0 0 0
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

#quality
  color: #676767
  font-size: 0.8em

</style>
