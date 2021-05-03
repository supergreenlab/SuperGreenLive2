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
  <div :id='$style.container'>
    <div :id='$style.plantid'>
      <span :id='$style.name'>{{ plant.name }}</span>&nbsp;in&nbsp;<span :id='$style.boxname'>{{ plant.box.name }}</span>
    </div>
    <div :id='$style.settings'>
      <div :class='$style.setting'>
        <img src='~/assets/icon_seeds.svg' />
        <div>
          {{ plant.settings.strain || 'Not set' }}<br/>
          <span :id='$style.thin'>from</span>&nbsp;<b :id='$style.green'>{{ plant.settings.seedBank || 'Not set' }}</b>
        </div>
      </div>
      <div :class='$style.setting'>
        <img src='~/assets/icon_phase.svg' />
        <div>
          <b>Germinated</b>: {{ germinated }}<br/>
          <b>{{ phase[0] }}</b>: {{ phase[1] }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>

import { DateTime, Interval } from 'luxon'

export default {
  props: ['plant',],
  computed: {
    germinated() {
      const { germinationDate } = this.$props.plant.settings
      if (!germinationDate) {
        return 'Not set'
      }
      const d = DateTime.fromISO(germinationDate)
      const i = Interval.fromDateTimes(d, DateTime.now())
      return `${Math.round(i.toDuration('days').toObject().days)} days ago`
    },
    phase() {
      const { germinationDate, veggingStart, bloomingStart } = this.$props.plant.settings
      if (bloomingStart) {
        const d = DateTime.fromISO(bloomingStart)
        const i = Interval.fromDateTimes(d, DateTime.now())
        return [
          'Blooming since',
          `${Math.round(i.toDuration('days').toObject().days)} days ago`
        ]
      } else if (veggingStart) {
        const d = DateTime.fromISO(veggingStart)
        const i = Interval.fromDateTimes(d, DateTime.now())
        return [
          'Vegging since',
          `${Math.round(i.toDuration('days').toObject().days)} days ago`
        ]
      } else if (germinationDate) {
        const d = DateTime.fromISO(germinationDate)
        const i = Interval.fromDateTimes(d, DateTime.now())
        return [
          'Started',
          `${Math.round(i.toDuration('days').toObject().days)} days ago`
        ]
      }
      return [
        'Phase',
        'Not set'
      ]
    }
  },
}

</script>

<style module lang=stylus>

#container
  display: flex
  flex: 1
  flex-direction: column

#name
  font-weight: bold
  font-size: 1.2em
  color: #454545

#boxname
  font-weight: bold
  font-size: 1.2em
  color: #3bb30b

#settings
  display: flex
  width: 100%
  margin: 15pt 0 10pt 0

.setting
  display: flex
  flex: 1
  color: #454545

.setting > img
  padding-right: 5pt

#green
  color: #3bb30b

#thin
  font-weight: 100

</style>
