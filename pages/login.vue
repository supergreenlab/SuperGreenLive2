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
  <section :id="$style.container">
    <form @submit='loginHandler'>
      <div :id='$style.body'>
        <div :id='$style.title'>S<span :id='$style.green'>G</span>L LOGIN</div>
        <input type='text' placeholder='Login' v-model='login' @change=''/>
        <input type='password' placeholder='Password' v-model='password' />
        <iframe ref='captchaFrame' :src='`${API_URL}/user/captcha`' width='100%' height='400px'></iframe>
        <div :id='$style.app'>No account yet? create one on the <a target='_blank' href='http://www.supergreenlab.com/app'>sgl app</a></div>
        <span :id='$style.error' v-if='error'>Wrong login/password</span>
        <div :id='$style.button'>
          <button @click='loginHandler'>LOGIN</button>
        </div>
      </div>
    </form>
  </section>
</template>

<script>

const API_URL=process.env.API_URL

export default {
  data() {
    return {
      API_URL,
      login: '',
      password: '',
    }
  },
  watch: {
    loggedIn(val) {
      if (val == true) {
        this.$router.replace('/plant')
      }
    },
  },
  mounted() {
    window.addEventListener('message', this.messageReceived, false)
    this.$refs.captchaFrame.onload = () => {
      console.log('pouet')
      this.$refs.captchaFrame.contentWindow.postMessage(`readyCaptcha|6Lc4abcmAAAAAPRQ1EAYfqjm5phbDGSbqefX1EXx|${document.location.origin}`, API_URL)
    }
  },
  destroyed() {
    window.removeEventListener('message', this.messageReceived)
  },
  methods: {
    loginHandler(e) {
      e.preventDefault()
      e.stopPropagation()
      const { login, password } = this.$data
      this.$store.dispatch('auth/login', { login, password })
      return false
    },
    messageReceived(e) {
      console.log('(page) messageReceived: ', e)
    },
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.loggedIn
    },
    error() {
      return this.$store.state.auth.error
    },
  },
}
</script>

<style module lang=stylus>
#container
  display: flex
  height: 100vh
  justify-content: center
  align-items: center

#body
  display: flex
  flex-direction: column

#body > input
  margin: 3pt 0
  padding: 3pt 6pt

#green
  color: #3bb30b

#title
  color: #454545
  font-weight: bold

#button
  display: flex
  justify-content: flex-end
  align-items: flex-end
  padding: 15pt 0 0 0

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

#error
  color: red

#app
  color: #454545

#body > iframe
  border: none

</style>
