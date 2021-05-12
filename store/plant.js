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

const STORAGE_ITEM='plant'
const API_URL=process.env.API_URL
const RPI_URL=process.env.RPI_URL

export const state = () => {
  let defaults = {
    plant: null
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
  async restorePlant(context, { token }) {
    const { data: timelapse } = await axios.get(`${RPI_URL}/timelapse`)
    if (timelapse.plantID) {
      const { data: plant } = await axios.get(`${API_URL}/plant/${timelapse.plantID}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })
      const { data: box } = await axios.get(`${API_URL}/box/${plant.boxID}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })
      plant.box = box

      context.commit('setPlant', plant)
    }
  },
}

export const mutations = {
  setState(state, newState) {
    Object.assign(state, newState)
  },
  setPlant(state, plant) {
    state.plant = plant
    storeState(state)
  },
}

export const getters = {
}
