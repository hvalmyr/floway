<script setup lang="ts">
import { contactInfo } from "~/constants/contact-info";

useSeoMeta({ title: "Реквизиты — flo-way" });

const { text } = await usePageContent();
</script>

<template>
  <UiGlassPage>
    <div class="flex flex-col gap-24">
      <h1 class="font-display text-h1 text-ink">Реквизиты</h1>
      <p class="font-body text-body text-ink">
        {{ text("legal_entity_name", contactInfo.legalEntity || "Реквизиты уточняются") }}
      </p>
      <p v-if="text('legal_inn', contactInfo.inn)" class="font-body text-body text-ink">
        ИНН: {{ text("legal_inn", contactInfo.inn) }}
      </p>
      <p v-if="text('legal_ogrn', contactInfo.ogrn)" class="font-body text-body text-ink">
        ОГРН: {{ text("legal_ogrn", contactInfo.ogrn) }}
      </p>
    </div>
  </UiGlassPage>
</template>
