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
  <div :id='$style.loading'>
    <div :id="$style.spinner" :style='{width, height, "margin-bottom": label ? "20pt" : ""}'>
      <div :id="$style.rect1"></div>
      <div :id="$style.rect2"></div>
      <div :id="$style.rect3"></div>
      <div :id="$style.rect4"></div>
      <div :id="$style.rect5"></div>
    </div>
    {{ label }}
  </div>
</template>

<script>
export default {
  props: {
    label: String,
    width: {
      default: '100%',
    },
    height: {
      default: '100%',
    },
  },
}
</script>

<style module lang=stylus>

#loading
  display: flex
  flex-direction: column
  align-items: center
  justify-content: center
  position: absolute
  top: 0
  left: 0
  width: 100%
  height: 100%

#spinner
  text-align: center
  font-size: 10px
  opacity: 0.5

#spinner > div
  background-color: #3BB30B
  height: 100%
  width: 10%
  display: inline-block
  
  -webkit-animation: sk-stretchdelay 1.2s infinite ease-in-out
  animation: sk-stretchdelay 1.2s infinite ease-in-out

#spinner #rect2
  -webkit-animation-delay: -1.1s
  animation-delay: -1.1s

#spinner #rect3
  -webkit-animation-delay: -1.0s
  animation-delay: -1.0s

#spinner #rect4
  -webkit-animation-delay: -0.9s
  animation-delay: -0.9s

#spinner #rect5
  -webkit-animation-delay: -0.8s
  animation-delay: -0.8s

@-webkit-keyframes sk-stretchdelay {
  0%, 40%, 100% { -webkit-transform: scaleY(0.4) }  
  20% { -webkit-transform: scaleY(1.0) }
}

@keyframes sk-stretchdelay {
  0%, 40%, 100% { 
    transform: scaleY(0.4)
    -webkit-transform: scaleY(0.4)
  }  20% { 
    transform: scaleY(1.0)
    -webkit-transform: scaleY(1.0)
  }
}
</style>
