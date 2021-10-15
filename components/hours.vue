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
    <div @click='minus' :id='$style.minus' :class='$style.button'><img src='~/assets/minus.svg' /></div>
    <div :id='$style.hour'>{{ String(value).padStart(2, '0') }}h</div>
    <div @click='plus' :id='$style.plus' :class='$style.button'><img src='~/assets/plus.svg' /></div>
  </section>
</template>

<script>
export default {
  props: ['value'],
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
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  align-items: center
  justify-content: center
  color: #454545

#hour
  display: flex
  align-items: center
  justify-content: center
  font-size: 3em
  font-weight: 600

.button
  cursor: pointer

</style>
