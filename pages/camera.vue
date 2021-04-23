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
    </div>
    <div :id='$style.button'>
      <button @click='nextHandler'>NEXT</button>
    </div>
  </section>
</template>

<script>
export default {
  methods: {
    nextHandler() {
      this.$router.push("/")
    }
  },
}
</script>

<style module lang=stylus>

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

</style>
