<!--
      Copyright (C) 2023  SuperGreenLab <towelie@supergreenlab.com>
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
    <iframe ref='captchaFrame' :src='`${API_URL}/user/captcha`' width='100%'></iframe>
  </section>
</template>

<script>

const API_URL=process.env.API_URL

export default {
  data() {
    return {
      API_URL,
    }
  },
  props: ['onToken',],
  mounted() {
    window.addEventListener('message', this.messageReceived, false)
    this.$refs.captchaFrame.onload = () => {
      this.$refs.captchaFrame.contentWindow.postMessage(`readyCaptcha|6Lc4abcmAAAAAPRQ1EAYfqjm5phbDGSbqefX1EXx|${document.location.origin}`, API_URL)
    }
  },
  destroyed() {
    window.removeEventListener('message', this.messageReceived)
  },

  methods: {
    messageReceived(e) {
      console.log('(page) messageReceived: ', e)
      if (e.data.length > 30) {
        this.onToken(e.data)
      }
    },
  },
}
</script>

<style module lang=stylus>

#container
  display: flex
  width: 100%
  min-width: 450px
  height: 100%
  min-height: 650px

#container > iframe
  border: none

</style>
