<script setup lang="ts">
import { contactInfo } from "~/constants/contact-info";

// TODO: остальные обязательные реквизиты (юридический адрес, ФИО директора
// и т.д.) заказчик не предоставлял — здесь только то, что видно на макете
// футера/контактов.
useSeoMeta({ title: "Реквизиты — ФлоВей" });
</script>

<template>
  <div class="container flex flex-col gap-24 py-64 sm:py-96">
    <h1 class="font-display text-h1 text-ink">Реквизиты</h1>
    <p class="font-body text-body text-ink">{{ contactInfo.legalEntity }}</p>
    <p class="font-body text-body text-ink">ИНН: {{ contactInfo.inn }}</p>
    <p class="font-body text-body text-ink">ОГРН: {{ contactInfo.ogrn }}</p>
  </div>
</template>
