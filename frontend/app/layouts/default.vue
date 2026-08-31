<template>
  <div class="flex min-h-screen flex-col">
    <AmbientTreeBackground />
    <AppHeader />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>
