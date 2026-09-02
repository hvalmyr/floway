<template>
  <div class="flex min-h-screen flex-col">
    <AmbientTreeBackground />
    <AppHeader />
    <GiftCertificateMarquee />
    <main class="flex-1">
      <slot />
    </main>
    <AppFooter />
    <CookieConsentBanner />
  </div>
</template>
