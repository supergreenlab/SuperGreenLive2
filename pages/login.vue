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
        <div v-if='!showCaptcha' :id='$style.title'>S<span :id='$style.green'>G</span>L LOGIN</div>
        <input v-if='!showCaptcha' type='text' placeholder='Login' v-model='login' @change=''/>
        <input v-if='!showCaptcha' type='password' placeholder='Password' v-model='password' />
        <Captcha v-if='showCaptcha' :onToken='onToken' />
        <div v-if='!showCaptcha' :id='$style.app'>No account yet? create one on the <a target='_blank' href='http://www.supergreenlab.com/app'>sgl app</a></div>
        <span :id='$style.error' v-if='error'>Wrong login/password</span>
        <div :id='$style.button'>
          <button @click='loginHandler'>LOGIN</button>
        </div>
      </div>
    </form>
  </section>
</template>

<script>

import Captcha from '~/components/captcha.vue'

export default {
  components: {Captcha,},
  data() {
    return {
      showCaptcha: false,
      login: '',
      password: '',
      token: '',
    }
  },
  watch: {
    loggedIn(val) {
      if (val == true) {
        this.$router.replace('/plant')
      }
    },
  },
  methods: {
    loginHandler(e) {
      if (e) {
        e.preventDefault()
        e.stopPropagation()
      }
      const { login, password, token } = this.$data
      if (!this.$data.showCaptcha && !login && !password) {
        return
      }
      if (this.$data.showCaptcha && !token) {
        return
      }
      if (!this.$data.showCaptcha) {
        this.$data.showCaptcha = true
        return
      }
      this.$store.dispatch('auth/login', { login, password, captcha: token })
      return false
    },
    onToken(token) {
      this.$data.token = token
      this.loginHandler()
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
  min-width: 400px
  justify-content: center
  align-items: center

#body > input
  margin: 3pt 0
  padding: 3pt 6pt
  max-width: 200px

#green
  color: #3bb30b

#title
  color: #454545
  font-weight: bold

#button
  display: flex
  justify-content: flex-end
  align-items: flex-end
  padding: 5pt 0 0 0

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

</style>
