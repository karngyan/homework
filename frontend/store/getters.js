export default {

  numberOfPages(state) {
    const total = state.customersMeta?.total ?? 0
    const per_page = state.customersMeta?.per_page ?? 25
    return Math.ceil(total / per_page)
  }
}
