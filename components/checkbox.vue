<!--
      Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
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
  <section :id='$style.container' v-on:click='click()'>
    <div :id='$style.border' :class='checked ? $style.selected : ""'>
      <div :id='$style.checkboxcenter'></div>
    </div>
    {{ label }} 
  </section>
</template>

<script>
export default {
  props: ['label', 'checked',],
  methods: {
    click() {
      this.$emit('click')
    },
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  align-items: center
  justify-content: center
  white-space: nowrap
  font-weight: bold
  color: #454545
  cursor: pointer

#border
  position: relative
  display: flex
  align-items: center
  justify-content: center
  width: 20pt
  height: 20pt
  margin-right: 5pt
  border: 1px solid #C4C4C4
  border-radius: 2pt
  transition: box-shadow 0.2s

#checkboxcenter
  width: 100%
  height: 100%
  border: 2pt solid white
  border-radius: 2pt
  transition: background-color 0.4s

.selected
  box-shadow: 0 0 0 2pt #3AB20B

.selected #checkboxcenter
  background-color: #3AB20B

</style>
