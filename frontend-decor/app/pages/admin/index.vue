<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

// Placeholder dashboard, proves the login/session flow works end to end.
// Real admin CRUD screens (object-types, landing-pages, portfolio, leads,
// content) land in Phase 3.
</script>

<template>
  <div>
    <h1 class="text-3xl font-semibold">Админ-панель flo-way</h1>
    <p class="mt-4 text-[var(--color-text-muted)]">
      Разделы редактирования появятся в следующей фазе.
    </p>
  </div>
</template>
