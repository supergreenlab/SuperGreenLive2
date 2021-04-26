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
        <h1>PLANTS ON THIS <span :class='$style.green'>TIMELAPSE</span>:</h1>
        <nuxt-link to='/plant' :id='$style.change'>change</nuxt-link></div>
      <Plant :plant='plant' />
      <div :id='$style.capture'><div v-for='src in srcs' v-if='src' :key='src' :style='{"background-image": `url(${src})`}'></div></div>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      n: 0,
      srcs: [null, 'http://192.168.1.26:8080/capture'],
    }
  },
  mounted() {
    this.interval = setInterval(() => {
      this.$data.srcs = [
        this.$data.srcs[1],
        `http://192.168.1.26:8080/capture?rand=${new Date().getTime()}`
      ]
    }, 20000)
  },
  destroyed() {
    clearInterval(this.interval)
  },
  methods: {
    nextHandler() {
      this.$router.push("/")
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
  padding: 0 5pt
  width: 100%
  max-width: 600pt

#header
  display: flex
  justify-content: space-between
  align-items: center

#header > h1
  margin: 20pt 0
  color: #454545

#change
  font-weight: 600
  color: #3bb30b
  text-decoration: none

.green
  color: #3bb30b

#capture
  position: relative
  height: 100%
  margin: 30pt 0

#capture > div
  position: absolute
  top: 0
  right: 0
  bottom: 0
  left: 0
  background-position: center
  background-repeat: center
  background-size: cover

</style>
