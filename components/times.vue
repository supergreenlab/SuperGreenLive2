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
      <img src='~assets/close.svg' :id='$style.close' @click='close' />
      <h1>Daily timelapse</h1>
      <p>Trigger the generation of my <b>daily</b> timelapse at:</p>
      <div :class='$style.params'>
        <Hours @change='changedDailyHour' :value='dailyHour' />
      </div>
      <h1>Weekly timelapse</h1>
      <p>Trigger the generation of my <b>weekly</b> timelapse at:</p>
      <div :class='$style.params'>
        <div :class='$style.param'>
          <Hours @change='changedWeeklyHour' :value='weeklyHour' />
        </div>
        <div :class='$style.param'>
          <Days @change='changedWeeklyDay' :value='weeklyDay' />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Hours from '~/components/hours.vue'
import Days from '~/components/days.vue'

export default {
  component: {Hours, Days,},
  data() {
    return {
      dailyHour: 0,
      weeklyDay: 0,
      weeklyHour: 0,
    }
  },
  methods: {
    close() {
      this.$emit('close')
    },
    changedDailyHour(value) {
      this.$data.dailyHour = value
    },
    changedWeeklyDay(value) {
      this.$data.weeklyDay = value
    },
    changedWeeklyHour(value) {
      this.$data.weeklyHour = value
    },
  },
}
</script>

<style module lang=stylus>

#container
  position: absolute
  display: flex
  align-items: center
  justify-content: center
  top: 0
  left: 0
  width: 100vw
  height: 100vh
  background-color: rgba(255, 255, 255, 0.5)

#close
  position: absolute
  top: 15pt
  right: 15pt
  cursor: pointer

#body
  position: relative
  display: flex
  flex-direction: column
  border-radius: 5pt
  border: 1pt solid #ababab
  background-color: white
  padding: 20pt 20pt
  width: 100%
  max-width: 600pt
  color: #323232

#body > h1
  text-transform: uppercase

#body > p > b
  font-weight: 600
  color: #3bb30b

.params
  display: flex
  align-items: center
  justify-content: center
  padding: 20pt 0
  @media only screen and (max-width: 1000pt)
    flex-direction: column

.param
  width: 200pt
  margin: 10pt

</style>
