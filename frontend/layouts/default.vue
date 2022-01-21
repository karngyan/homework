<template>
<div>
  <Nav />
  <Nuxt />
  <Footer />
</div>
</template>

<script>
export default {
  name: "default"
}
</script>

<style scoped>
html {
  @apply font-sans bg-white;
}

/* Page-transition styles must be global */
/* page is used by default */
.page-leave-active,
.page-enter-active {
  transition: opacity 0.3s;
}
.page-enter,
.page-leave-to {
  opacity: 0;
}

::-moz-selection { @apply bg-violet-600 text-white; }
::selection { @apply bg-violet-600 text-white; }
</style>
