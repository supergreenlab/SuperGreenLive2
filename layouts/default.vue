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
  <section :id='$style.container'>
    <section :id='$style.header'>
      <img src='~/assets/logo_sgl_white.svg' />
    </section>
    <section :id='$style.body'>
      <Nuxt />
    </section>
  </section>
</template>

<style module lang=stylus>
html
  font-family:
    'Source Sans Pro',
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    'Helvetica Neue',
    Arial,
    sans-serif
  font-size: 16px
  word-spacing: 1px
  -ms-text-size-adjust: 100%
  -webkit-text-size-adjust: 100%
  -moz-osx-font-smoothing: grayscale
  -webkit-font-smoothing: antialiased
  box-sizing: border-box
  min-height: 100vh

*,
*::before,
*::after
  box-sizing: border-box
  margin: 0

#container
  margin: 0 auto
  min-height: 100vh
  display: flex
  flex-direction: column

#header
  position: absolute
  display: flex
  padding: 10pt
  width: 100%
  top: 0
  left: 0

#header > img
  height: 30pt

#body
  flex: 1

</style>
