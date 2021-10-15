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
    <v-select :options='options' :value='currentDay' label='name' @input='setSelected' />
  </section>
</template>

<script>
export default {
  props: ['value'],
  computed: {
    options() {
      return [
        { name: 'Monday', n: 1 },
        { name: 'Tuesday', n: 2 },
        { name: 'Wednesday', n: 3 },
        { name: 'Thursday', n: 4 },
        { name: 'Friday', n: 5 },
        { name: 'Saturday', n: 6 },
        { name: 'Sunday', n: 0 },
      ]
    },
    currentDay() {
      return this.options.find(o => o.n == this.$props.value)
    },
  },
  methods: {
    setSelected(value) {
      this.$emit('change', value ? value.n : null)
    }
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  flex-direction: column

</style>
