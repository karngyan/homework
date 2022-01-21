<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="customer">
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            {{ customer.attributes.email }}
          </h2>
        </div>
      </div>

      <div class="text-md text-gray-600 pt-2">Last updated at: {{ customer.last_updated | formatDate }}</div>

      <div class="pt-6">
        <h3 class="text-xl leading-6 font-medium text-gray-900">
          Attributes
        </h3>
      </div>
      <div class="mt-5">
        <dl class="divide-y divide-gray-100">
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              id
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.id }}</span>
            </dd>
          </div>

          <!-- keeping these 2 fields on top as they are most common and fixed -->
          <div class="py-1 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              email
            </dt>
            <dd class="flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <input v-model="attributes['email']" name="email" type="email" class="block max-w-lg w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border-gray-300 rounded-md">
            </dd>
          </div>
          <div class="py-1 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              created_at
            </dt>
            <dd class="flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <input v-model="attributes['created_at']" name="created_at" type="text" class="block max-w-lg w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border-gray-300 rounded-md">
            </dd>
          </div>

          <!--  other attributes -->
          <div v-if="key !== 'email' && key !== 'created_at'" v-for="(key, ind) in Object.keys(attributes)" :key="`${key}-${ind}`"
               class="py-1 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              {{ key }}
            </dt>
            <dd class="flex items-center text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <input v-model="attributes[key]" :name="key" type="text" class="block max-w-lg w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border-gray-300 rounded-md">
              <div @click="removeAttribute(key)" class="ml-6 cursor-pointer text-red-600 hover:text-red-700 underline underline-offset-2">Remove</div>
            </dd>
          </div>

          <!-- new attribute -->
          <div class="py-10 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              <input v-model="newAttribute.name" placeholder="name" name="new-name" type="text" class="block w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border-gray-300 rounded-md">
            </dt>
            <dd class="flex items-center text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <input v-model="newAttribute.value" placeholder="value" name="new-value" type="text" class="block max-w-lg w-full shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm border-gray-300 rounded-md">
              <div @click="addNewAttribute" class="ml-6 cursor-pointer text-gray-600 hover:text-gray-700 underline underline-offset-2">Add</div>
            </dd>
            <div class="-mt-3 ml-1 text-sm text-red-500 font-normal">{{ newAttribute.error }}</div>
          </div>
        </dl>
      </div>

      <div class="flex flex-row space-x-4 mr-2 justify-end items-center">
        <button @click="discardChanges" class="cursor-pointer text-md underline text-gray-500 hover:text-gray-700 underline-offset-2">
          <!-- TODO: Add a confirm modal -->
          Discard changes
        </button>
        <button @click="updateAttributes" type="button" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          Save changes
        </button>
        <div class="h-5 w-5">
          <svg v-show="showSavedTick" xmlns="http://www.w3.org/2000/svg" class="text-emerald-600 h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
        </div>
      </div>


    </div>
    <div v-else>
      <div class="min-h-screen flex items-center justify-center">
        {{ message }}
      </div>
    </div>
  </div>
</template>

<script>
import {format} from "date-fns";

export default {
  async fetch() {
    try {
      await this.$store.dispatch('fetchCustomerById', {id: this.id})
      this.attributes = { ...this.customer.attributes }
    } catch (err) {
      console.error(err)
      this.message = err.toString()
    }
  },
  data() {
    return {
      message: 'fetching customer data ...',
      attributes: {},
      newAttribute: {
        name: '',
        value: '',
        error: ''
      },
      showSavedTick: false,
    }
  },
  computed: {
    id() {
      return parseInt(this.$route.params.id)
    },
    customer() {
      return this.$store.state.customers[this.id]
    }
  },
  methods: {
    removeAttribute(key) {
      // tricky so that vue knows about the changes
      const xAttributes = { ...this.attributes }
      delete xAttributes[key]
      this.attributes = xAttributes
    },
    addNewAttribute() {
      this.newAttribute.name = this.newAttribute.name.trim()
      this.newAttribute.value = this.newAttribute.value.trim()

      if (this.newAttribute.name === '' || this.newAttribute.value === '') {
        this.newAttribute.error = 'empty values are not acceptable'
        return
      }

      if (this.newAttribute.name in this.attributes) {
        this.newAttribute.error = 'attribute name already exists'
        return
      }

      // new attribute name and non-empty name/value
      this.attributes[this.newAttribute.name] = this.newAttribute.value
      this.newAttribute = {
        name: '',
        value: '',
        error: ''
      }
    },
    discardChanges() {
      this.attributes = {...this.customer.attributes}
    },
    updateAttributes() {
      this.$store.dispatch('patchCustomerAttributes', {id: this.id, attributes: this.attributes})
        .then(() => {
          this.showSavedTick = true
          setTimeout(() => {
            this.showSavedTick = false
          }, 3000)
        }).catch((err) => {
          console.error(err)
      })
    }
  },
  filters: {
    formatDate(secondsSinceEpoch) {
      return format(new Date(secondsSinceEpoch * 1000), 'MMM do yyyy, KK:mm aaa')
    }
  }
}
</script>

<style scoped>

</style>
