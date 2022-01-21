import Vue from 'vue';
import getDefaultState from "@/store/state";

export default {

  resetState(state) {
    Object.assign(state, getDefaultState())
  },

  setResource(state, {resource, value}) {
    Vue.set(state, resource, value)
  },

  setItem(state, {item, id, resource}) {

    if (!state[resource]) {
      state[resource] = {}
    }
    Vue.set(state[resource], id, item)
  },

  deleteItem(state, {id, resource}) {

    if (!state[resource]) {
      state[resource] = {}
    }

    Vue.delete(state[resource], id)
  },
}
