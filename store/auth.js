/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */


import Vue from 'vue'
import axios from 'axios'

import { loadFromStorage, saveToStorage } from '~/lib/client-side.js'

const STORAGE_ITEM='auth'
const API_URL=process.env.API_URL
const RPI_URL=process.env.RPI_URL

export const state = () => {
  let defaults = {
    error: false,
    loading: false,
    token: '',
  };
  return defaults
};

const storeState = (state) => {
  saveToStorage(STORAGE_ITEM, JSON.stringify(state))
}

export const actions = {
  nuxtClientInit(context) {
    const saved = loadFromStorage(STORAGE_ITEM)
    if (saved) {
      context.commit('setState', JSON.parse(saved))
    }
  },
  async login({ commit, dispatch }, { login, password }) {
    commit('setLoading', true)
    try {
      const resp = await axios.post(`${API_URL}/login`, {
        handle: login,
        password,
      })
      const token = resp.headers['x-sgl-token']
      const { data: respToken } = await axios.post(`${RPI_URL}/token`, {
        token,
      })
      if (respToken == 'ALREADY_LOGGED_IN') {
        await dispatch('plant/restorePlant', { token }, { root: true })
      }

      commit('setToken', token)
    } catch(e) {
      commit('setError', true)
    }
    commit('setLoading', false)
  },
}

export const mutations = {
  setState(state, newState) {
    Object.assign(state, newState)
  },
  setToken(state, token) {
    state.token = token
    storeState(state)
  },
  setLoading(state, loading) {
    state.loading = loading
  },
  setError(state, error) {
    state.error = error
  },
}

export const getters = {
  loggedIn: state => !!state.token,
  error: state => state.error,
}
