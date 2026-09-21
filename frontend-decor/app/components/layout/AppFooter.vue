<script setup lang="ts">
import { socialIcons } from "~/constants/social-icons";
import { contactInfo } from "~/constants/contact-info";

const { text } = await usePageContent();
const { socialLinks } = await useSocialLinks();
// "Связь со школой": те же мессенджеры, что и в форме заявки (contactMethod),
// плюс Telegram уже есть в "соцсети" ниже.
const { contactChannels } = await useContactChannels();
const route = useRoute();
const year = new Date().getFullYear();

function isActive(to: string) {
  return to !== "/" && route.path.startsWith(to);
}

// Placeholder — see AppHeader.vue's own note; real page set lands Phase 4.
const schoolLinks = [
  { to: "/portfolio", label: "Портфолио" },
  { to: "/how-it-works", label: "Как мы работаем" },
  { to: "/contacts", label: "Контакты" },
];

const documentLinks = [
  { to: "/privacy", label: "политика конфиденциальности" },
  { to: "/cookie-policy", label: "политика использования cookie" },
  { to: "/terms", label: "пользовательское соглашение" },
  { to: "/pd-consent", label: "согласие на обработку персональных данных" },
];
</script>

<template>
  <footer class="mt-auto bg-primary text-white">
    <div class="container flex flex-col gap-32 py-64">
      <!-- Same size as the desktop header wordmark (AppHeader.vue's centered
      logo) at every breakpoint — not shrunk down on mobile/tablet. -->
      <div class="flex items-center gap-16">
        <NuxtLink to="/" class="font-display text-h2 text-white">Фловей</NuxtLink>
        <IconFirBranch class="h-32 w-auto text-white" aria-hidden="true" />
      </div>

      <div class="grid grid-cols-1 gap-32 md:grid-cols-3">
        <div class="flex flex-col gap-16">
          <h3 class="font-display text-h4 text-white">связаться с нами</h3>
          <div class="flex gap-16">
            <a
              v-for="channel in contactChannels"
              :key="channel.label"
              :href="channel.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="channel.label"
              class="grid size-[44px] place-items-center rounded-full bg-white text-primary hover:opacity-80"
            >
              <component :is="channel.icon" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>

          <h3 class="mt-16 font-display text-h4 text-white">соцсети</h3>
          <div class="flex gap-16">
            <a
              v-for="social in socialLinks"
              :key="social.label"
              :href="social.href"
              target="_blank"
              rel="noopener noreferrer"
              :aria-label="social.label"
              class="grid size-[44px] place-items-center rounded-full bg-white text-primary hover:opacity-80"
            >
              <component :is="socialIcons[social.label]" class="size-[20px]" aria-hidden="true" />
            </a>
          </div>
          <p
            v-for="social in socialLinks.filter((s) => s.disclaimer)"
            :key="`${social.label}-disclaimer`"
            class="font-body text-body text-white/70"
          >
            {{ social.disclaimer }}
          </p>
        </div>

        <nav class="flex flex-col gap-16" aria-label="Оформление">
          <h3 class="font-display text-h4 text-white">оформление</h3>
          <NuxtLink
            v-for="link in schoolLinks"
            :key="link.to"
            :to="link.to"
            class="font-body text-body text-white hover:underline"
            :class="isActive(link.to) ? 'underline' : ''"
          >
            {{ link.label }}
          </NuxtLink>
        </nav>

        <div class="flex flex-col gap-16">
          <h3 class="font-display text-h4 text-white">документы</h3>
          <NuxtLink
            v-for="link in documentLinks"
            :key="link.to"
            :to="link.to"
            class="font-body text-body text-white hover:underline"
            :class="isActive(link.to) ? 'underline' : ''"
          >
            {{ link.label }}
          </NuxtLink>
          <p class="font-body text-body text-white">
            {{ text("legal_entity_name", contactInfo.legalEntity) }}
          </p>
          <p class="font-body text-body text-white">
            ИНН: {{ text("legal_inn", contactInfo.inn) }}
          </p>
          <p class="font-body text-body text-white">
            ОГРН: {{ text("legal_ogrn", contactInfo.ogrn) }}
          </p>
        </div>
      </div>
    </div>

    <p class="container py-24 font-body text-body text-white/80">
      &copy; {{ year }} flo-way — новогоднее оформление
    </p>
  </footer>
</template>
