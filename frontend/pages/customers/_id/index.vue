<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div v-if="customer">
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            {{ customer.attributes.email }}
          </h2>
        </div>
        <div class="mt-4 flex md:mt-0">
          <nuxt-link :to="`/customers/${id}/edit`"
                     class="ml-3 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Edit attributes
          </nuxt-link>
        </div>
      </div>

      <div class="text-md text-gray-600 pt-2">Last updated at: {{ customer.last_updated | formatDate }}</div>

      <div class="pt-6">
        <h3 class="text-xl leading-6 font-medium text-gray-900">
          Attributes
        </h3>
      </div>
      <div class="mt-5">
        <dl class="">
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              id
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.id }}</span>
            </dd>
          </div>
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              email
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.attributes.email }}</span>
            </dd>
          </div>
          <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              created_at
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.attributes.created_at }}</span>
            </dd>
          </div>
          <div v-if="key !== 'email' && key !== 'created_at'" v-for="(key, ind) in Object.keys(customer.attributes)" :key="`${key}-${ind}`"
               class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              {{ key }}
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.attributes[key] }}</span>
            </dd>
          </div>
        </dl>
      </div>

      <div class="pt-6">
        <h3 class="text-xl leading-6 font-medium text-gray-900">
          Events
        </h3>
      </div>
      <div class="mt-5">
        <dl class="">
          <div class="py-2 border-b bg-gray-100 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              Event name
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">Count</span>
            </dd>
          </div>
          <div v-for="(key, ind) in Object.keys(customer.events)" :key="`${key}-${ind}`"
               class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
            <dt class="text-sm font-medium text-gray-500">
              {{ key }}
            </dt>
            <dd class="mt-1 flex text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span class="flex-grow">{{ customer.events[key] }}</span>
            </dd>
          </div>
        </dl>
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
    } catch (err) {
      console.error(err)
      this.message = err.toString()
    }
  },
  data() {
    return {
      message: 'fetching customer data ...'
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
  filters: {
    formatDate(secondsSinceEpoch) {
      return format(new Date(secondsSinceEpoch * 1000), 'MMM do yyyy, KK:mm aaa')
    }
  }
}
</script>

<style scoped>

</style>
