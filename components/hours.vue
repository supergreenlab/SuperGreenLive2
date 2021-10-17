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
    <div v-if='!inputMode' :id='$style.input'>
      <div @click='minus' :id='$style.minus' :class='$style.button'><img src='~/assets/minus.svg' /></div>
      <div :id='$style.hour' @click='clickInputMode'>{{ String(value).padStart(2, '0') }}h</div>
      <div @click='plus' :id='$style.plus' :class='$style.button'><img src='~/assets/plus.svg' /></div>
    </div>
    <div v-else :id='$style.input' :class='$style.inputMode'>
      <input @keypress='listenEnter' ref='input' type='number' :value='value' />
      <button @click='save'>save</button>
    </div>
  </section>
</template>

<script>
export default {
  props: ['value'],
  data() {
    return {
      inputMode: false
    }
  },
  methods: {
    minus() {
      let value = this.$props.value - 1
      if (value < 0) {
        value = 23
      }
      this.$emit('change', value)
    },
    plus() {
      this.$emit('change', (this.$props.value + 1) % 24)
    },
    clickInputMode() {
      this.$data.inputMode = true
      setTimeout(() => {
        this.$refs['input'].select()
      }, 50)
    },
    save() {
      let val = parseInt(this.$refs['input'].value) % 24
      if (val < 0) {
        val = 0
      }
      if (isNaN(val)) {
        return
      }
      this.$emit('change', val)
      this.$data.inputMode = false
    },
    listenEnter(e) {
      if (e.key == "Enter") {
        this.save()
      }
    },
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  color: #454545

#input
  display: flex
  align-items: center
  justify-content: center
  height: 40pt
  overflow: hidden

.inputMode
  border: 1pt solid #ababab
  border-radius: 5pt

#input > input
  border: none
  background-color: transparent
  text-align: center
  font-size: 3em
  font-weight: 600
  width: 75pt
  color: #454545

#input > button
  background-color: #3bb30b
  height: 100%
  border: none
  color: white
  padding: 0 10pt
  cursor: pointer
  border-radius: 0

#hour
  display: flex
  align-items: center
  justify-content: center
  font-size: 3em
  font-weight: 600

.button
  cursor: pointer
  margin: 0 10pt

</style>
