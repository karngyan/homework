<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">

    <div class="pb-5 border-b border-gray-200">
      <h3 class="text-lg leading-6 font-medium text-gray-900">
        Customers
      </h3>
    </div>

    <div class="pt-12 flex flex-col">
      <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
          <div class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ID
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Email
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Last Updated
                </th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(customer, ind) in customers" :key="customer.id" @click="routeToCustomerPage(customer.id)"
                  class="cursor-pointer hover:bg-gray-100" :class="ind%2 === 0 ? 'bg-white' : 'bg-gray-50'">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {{ customer.id }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ customer.attributes.email }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ customer.last_updated | formatDate }}
                </td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <nav class="border-t border-gray-200 px-4 mt-12 flex items-center justify-between sm:px-0">
      <div class="-mt-px w-0 flex-1 flex">
        <button @click.prevent="previousPage"
                class="border-t-2 border-transparent pt-4 pr-1 inline-flex items-center text-sm font-medium text-gray-500 hover:text-gray-700 hover:border-gray-300">
          <svg class="mr-3 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
               fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd"
                  d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
                  clip-rule="evenodd"/>
          </svg>
          Previous
        </button>
      </div>
      <div class="hidden md:-mt-px md:flex">
      <span v-if="page - 2 >= 1"
            class="border-transparent text-gray-500 border-t-2 pt-4 px-4 inline-flex items-center text-sm font-medium">
        ...
      </span>
        <button v-if="(page - 1) >= 1" @click="fetchPage(page - 1)"
                class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 border-t-2 pt-4 px-4 inline-flex items-center text-sm font-medium"
                aria-current="page">
          {{ page - 1 }}
        </button>
        <div
          class="border-indigo-500 text-indigo-600 border-t-2 pt-4 px-4 inline-flex items-center text-sm font-medium">
          {{ page }}
        </div>
        <button v-if="(page + 1) <= numberOfPages" @click="fetchPage(page + 1)"
                class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300  border-t-2 pt-4 px-4 inline-flex items-center text-sm font-medium"
                aria-current="page">
          {{ page + 1 }}
        </button>
        <span v-if="page + 2 <= numberOfPages"
              class="border-transparent text-gray-500 border-t-2 pt-4 px-4 inline-flex items-center text-sm font-medium">
        ...
      </span>
      </div>
      <div class="-mt-px w-0 flex-1 flex justify-end">
        <button @click.prevent="nextPage"
                class="border-t-2 border-transparent pt-4 pl-1 inline-flex items-center text-sm font-medium text-gray-500 hover:text-gray-700 hover:border-gray-300">
          Next
          <svg class="ml-3 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
               fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd"
                  d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
                  clip-rule="evenodd"/>
          </svg>
        </button>
      </div>
    </nav>

  </div>
</template>

<script>
import {format} from 'date-fns'


export default {
  async fetch() {
    try {
      await this.fetchCustomers(this.page)
    } catch (err) {
      console.error(err)
    }
  },
  methods: {
    fetchCustomers(page, per_page = 25) {
      return this.$store.dispatch('fetchAllCustomers', {page, per_page})
    },
    previousPage() {
      if (this.page === 1) return
      this.page--
      this.fetchCustomers(this.page)
    },
    nextPage() {
      console.debug(this.numberOfPages, this.page)
      if (this.page >= this.numberOfPages) return
      this.page++
      this.fetchCustomers(this.page)
    },
    fetchPage(page) {
      if (this.page === page) return
      this.page = page
      this.fetchCustomers(this.page)
    },
    routeToCustomerPage(id) {
      this.$router.push(`/customers/${id}`)
    }
  },
  computed: {
    customers() {
      return Object.values(this.$store.state.customers)
    },
    numberOfPages() {
      return this.$store.getters.numberOfPages
    },
    meta() {
      return this.$store.state.customersMeta
    }
  },
  data() {
    return {
      page: 1,
      per_page: 25,
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
