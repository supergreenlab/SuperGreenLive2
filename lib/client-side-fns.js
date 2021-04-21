/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

module.exports = {
  onNuxtReady: window.onNuxtReady,
  loadFromStorage: window.localStorage.getItem.bind(window.localStorage),
  saveToStorage: window.localStorage.setItem.bind(window.localStorage),
  addEventListener: window.addEventListener,
  removeEventListener: window.removeEventListener,
  innerHeight: () => window.innerHeight,
  innerWidth: () => window.innerWidth,
  screenX: () => window.screenX,
  availHeight: () => window.screen.availHeight,
  availWidth: () => window.screen.availWidth,
  open: window.open,
  setHref: (href) => window.location.href = href,
  postMessage: () => window.opener.postMessage.bind(window.opener),
  userAgent: navigator.userAgent,
  baseUrl: `${window.location.protocol}//${window.location.host}`
}
