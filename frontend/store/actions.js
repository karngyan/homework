export default {

  fetchAllCustomers({state, commit}, {page = 1, per_page = 25}) {
    return new Promise((resolve, reject) => {
      this.$axios.get('/customers', {
        params: {
          page,
          per_page
        }
      }).then((resp) => {
        const customers = resp.data.customers
        // clear
        commit('setResource', {resource: 'customers', value: {}})
        // set all
        customers.forEach((customer) => {
          commit('setItem', {resource: 'customers', id: customer.id, item: customer})
        })
        // set meta
        commit('setResource', {resource: 'customersMeta', value: resp.data.meta})
        resolve(resp.data)
      }).catch(err => {
        reject(err)
      })
    })
  },

  fetchCustomerById({state, commit}, {id}) {
    return new Promise((resolve, reject) => {
      if (state.customers[id]) {
        resolve(state.customers[id])
        return
      }

      this.$axios.get(`/customers/${id}`)
        .then((resp) => {
          const customer = resp.data.customer
          commit('setItem', {id: customer.id, item: customer, resource: 'customers'})
          resolve(customer)
        }).catch((err) => {
        reject(err)
      })

    })
  },

}
