<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

import type {
  ThankYouPage,
  ThankYouPageFaqItem,
  ThankYouPagePhoto,
  ThankYouPageVariant,
} from "~/types/api";

const variantTabs: { value: ThankYouPageVariant; label: string }[] = [
  { value: "course", label: "Курс" },
  { value: "masterclass", label: "Мастер-класс" },
  { value: "trial_lesson", label: "Пробное занятие" },
];

const activeVariant = ref<ThankYouPageVariant>("course");
const api = useApiClient();

const settings = ref<ThankYouPage | null>(null);
const photos = ref<ThankYouPagePhoto[]>([]);
const faqItems = ref<ThankYouPageFaqItem[]>([]);
const loading = ref(false);
const loadError = ref("");

async function load() {
  loading.value = true;
  loadError.value = "";
  settings.value = null;
  try {
    const data = await api<ThankYouPage>(`/api/v1/thank-you-pages/${activeVariant.value}`);
    settings.value = data;
    photos.value = data.photos.slice().sort((a, b) => a.sortOrder - b.sortOrder);
    faqItems.value = data.faqItems.slice().sort((a, b) => a.sortOrder - b.sortOrder);
  } catch {
    loadError.value = "Не удалось загрузить данные";
  } finally {
    loading.value = false;
  }
}

await load();
watch(activeVariant, load);

// --- settings (title/subtitle/description/show-* flags) -----------------

const savingSettings = ref(false);
const settingsError = ref("");
const settingsSaved = ref(false);

async function saveSettings() {
  if (!settings.value) return;
  savingSettings.value = true;
  settingsError.value = "";
  settingsSaved.value = false;
  try {
    await api(`/api/v1/thank-you-pages/${activeVariant.value}`, {
      method: "PUT",
      body: {
        title: settings.value.title,
        subtitle: settings.value.subtitle,
        description: settings.value.description,
        heroImage: settings.value.heroImage,
        showMessengers: settings.value.showMessengers,
        showSocialLinks: settings.value.showSocialLinks,
        showBlogLink: settings.value.showBlogLink,
        blogLinkText: settings.value.blogLinkText,
        blogLinkUrl: settings.value.blogLinkUrl,
        showCarousel: settings.value.showCarousel,
        showFaq: settings.value.showFaq,
        showCommunity: settings.value.showCommunity,
        communityText: settings.value.communityText,
        communityUrl: settings.value.communityUrl,
      },
    });
    settingsSaved.value = true;
  } catch {
    settingsError.value = "Не удалось сохранить";
  } finally {
    savingSettings.value = false;
  }
}

// --- photos (carousel) ----------------------------------------------------

const photoForm = ref({ image: "" });
const savingPhoto = ref(false);

const { draggingIndex: draggingPhotoIndex, onPointerDown: onPhotoPointerDown } =
  useAdminDragReorder(photos, (item) =>
    api(`/api/v1/thank-you-pages/${activeVariant.value}/photos/${item.id}`, {
      method: "PUT",
      body: { image: item.image, sortOrder: item.sortOrder },
    }),
  );

async function addPhoto() {
  if (!photoForm.value.image) return;
  savingPhoto.value = true;
  try {
    const created = await api<ThankYouPagePhoto>(
      `/api/v1/thank-you-pages/${activeVariant.value}/photos`,
      { method: "POST", body: { image: photoForm.value.image, sortOrder: photos.value.length } },
    );
    photos.value = [...photos.value, created];
    photoForm.value = { image: "" };
  } catch {
    // no dedicated error slot — AdminImageUpload already surfaces upload failures
  } finally {
    savingPhoto.value = false;
  }
}

async function removePhoto(id: number) {
  if (!confirm("Удалить фото из карусели?")) return;
  await api(`/api/v1/thank-you-pages/${activeVariant.value}/photos/${id}`, { method: "DELETE" });
  photos.value = photos.value.filter((p) => p.id !== id);
}

// --- FAQ items --------------------------------------------------------

const faqForm = ref({ question: "", answer: "" });
const savingFaq = ref(false);
const faqError = ref("");

const { draggingIndex: draggingFaqIndex, onPointerDown: onFaqPointerDown } = useAdminDragReorder(
  faqItems,
  (item) =>
    api(`/api/v1/thank-you-pages/${activeVariant.value}/faq-items/${item.id}`, {
      method: "PUT",
      body: { question: item.question, answer: item.answer, sortOrder: item.sortOrder },
    }),
);

async function addFaqItem() {
  faqError.value = "";
  if (!faqForm.value.question || !faqForm.value.answer) return;
  savingFaq.value = true;
  try {
    const created = await api<ThankYouPageFaqItem>(
      `/api/v1/thank-you-pages/${activeVariant.value}/faq-items`,
      {
        method: "POST",
        body: {
          question: faqForm.value.question,
          answer: faqForm.value.answer,
          sortOrder: faqItems.value.length,
        },
      },
    );
    faqItems.value = [...faqItems.value, created];
    faqForm.value = { question: "", answer: "" };
  } catch {
    faqError.value = "Не удалось сохранить — проверьте поля";
  } finally {
    savingFaq.value = false;
  }
}

async function removeFaqItem(id: number) {
  if (!confirm("Удалить вопрос?")) return;
  await api(`/api/v1/thank-you-pages/${activeVariant.value}/faq-items/${id}`, {
    method: "DELETE",
  });
  faqItems.value = faqItems.value.filter((item) => item.id !== id);
}
</script>

<template>
  <div>
    <NuxtLink to="/admin" class="text-sm text-[var(--color-primary)] hover:underline"
      >← К дашборду</NuxtLink
    >
    <h1 class="mt-2 text-2xl font-semibold">Страница благодарности</h1>
    <p class="mt-2 text-sm text-[var(--color-text-muted)]">
      Содержимое страницы, на которую попадает посетитель сразу после отправки заявки — отдельный
      текст и настройки для каждого типа заявки.
    </p>

    <div class="mt-6 flex gap-2 border-b border-gray-200">
      <button
        v-for="tab in variantTabs"
        :key="tab.value"
        type="button"
        class="border-b-2 px-4 py-2 text-sm font-medium"
        :class="
          activeVariant === tab.value
            ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
            : 'border-transparent text-[var(--color-text-muted)] hover:text-ink'
        "
        @click="activeVariant = tab.value"
      >
        {{ tab.label }}
      </button>
    </div>

    <p v-if="loading" class="mt-6 text-[var(--color-text-muted)]">Загрузка…</p>
    <p v-else-if="loadError" class="mt-6 text-red-600">{{ loadError }}</p>

    <template v-else-if="settings">
      <form
        class="mt-6 grid gap-3 rounded border border-gray-200 bg-white p-4"
        @submit.prevent="saveSettings"
      >
        <h2 class="font-semibold">Заголовок, подзаголовок, описание</h2>
        <input
          v-model="settings.title"
          type="text"
          placeholder="Заголовок"
          required
          class="rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="settings.subtitle"
          type="text"
          placeholder="Подзаголовок"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <AdminMarkdownField v-model="settings.description" placeholder="Описание" :rows="4" />
        <AdminImageUpload v-model="settings.heroImage" label="Квадратное фото (Hero)" />

        <hr class="my-2 border-gray-200" />

        <h2 class="font-semibold">Альтернативный контакт (мессенджеры)</h2>
        <p class="text-sm text-[var(--color-text-muted)]">
          Иконки берутся из общих настроек контактов —
          <NuxtLink
            to="/admin/page-content/info"
            class="text-[var(--color-primary)] hover:underline"
            >Контакты и реквизиты</NuxtLink
          >. Здесь только показываете или скрываете их на этой странице.
        </p>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showMessengers" type="checkbox" class="size-5" />
          Показывать мессенджеры
        </label>

        <hr class="my-2 border-gray-200" />

        <h2 class="font-semibold">Соцсети</h2>
        <p class="text-sm text-[var(--color-text-muted)]">
          Ссылки берутся из общего списка —
          <NuxtLink to="/admin/social-links" class="text-[var(--color-primary)] hover:underline"
            >Соцсети</NuxtLink
          >. Здесь только показываете или скрываете их на этой странице.
        </p>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showSocialLinks" type="checkbox" class="size-5" />
          Показывать соцсети
        </label>

        <hr class="my-2 border-gray-200" />

        <h2 class="font-semibold">Кнопка в блог</h2>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showBlogLink" type="checkbox" class="size-5" />
          Показывать кнопку
        </label>
        <input
          v-model="settings.blogLinkText"
          type="text"
          placeholder="Текст кнопки, например «Читать блог»"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="settings.blogLinkUrl"
          type="text"
          placeholder="Ссылка, например /blog"
          class="rounded border border-gray-300 px-3 py-2"
        />

        <hr class="my-2 border-gray-200" />

        <h2 class="font-semibold">Приглашение в сообщество</h2>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showCommunity" type="checkbox" class="size-5" />
          Показывать приглашение
        </label>
        <input
          v-model="settings.communityText"
          type="text"
          placeholder="Текст кнопки, например «Вступить в чат учеников»"
          class="rounded border border-gray-300 px-3 py-2"
        />
        <input
          v-model="settings.communityUrl"
          type="text"
          placeholder="Ссылка на чат или канал"
          class="rounded border border-gray-300 px-3 py-2"
        />

        <hr class="my-2 border-gray-200" />

        <h2 class="font-semibold">Карусель и мини-FAQ</h2>
        <p class="text-sm text-[var(--color-text-muted)]">
          Сами фото и вопросы редактируются в блоках ниже — здесь только показываете или скрываете
          блок целиком.
        </p>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showCarousel" type="checkbox" class="size-5" />
          Показывать карусель фото
        </label>
        <label class="flex items-center gap-2 text-sm">
          <input v-model="settings.showFaq" type="checkbox" class="size-5" />
          Показывать мини-FAQ
        </label>

        <p v-if="settingsError" class="text-sm text-red-600">{{ settingsError }}</p>
        <p v-if="settingsSaved" class="text-sm text-green-700">Сохранено</p>

        <div>
          <button
            type="submit"
            :disabled="savingSettings"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
          >
            {{ savingSettings ? "Сохранение…" : "Сохранить" }}
          </button>
        </div>
      </form>

      <div class="mt-6 rounded border border-gray-200 bg-white p-4">
        <h2 class="font-semibold">Карусель фото (социальное доказательство)</h2>
        <form class="mt-3 flex flex-wrap items-end gap-3" @submit.prevent="addPhoto">
          <AdminImageUpload v-model="photoForm.image" label="Фото" />
          <button
            type="submit"
            :disabled="savingPhoto || !photoForm.image"
            class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
          >
            Добавить
          </button>
        </form>

        <table v-if="photos.length" class="mt-4 w-full border-collapse text-sm">
          <tbody>
            <tr
              v-for="(photo, index) in photos"
              :key="photo.id"
              :data-row-index="index"
              class="border-b border-gray-100 last:border-0"
              :class="draggingPhotoIndex === index ? 'opacity-50' : ''"
            >
              <td class="w-8 px-2 py-2">
                <AdminDragHandle @pointerdown="onPhotoPointerDown(index, $event)" />
              </td>
              <td class="px-2 py-2">
                <img
                  :src="resolveMediaUrl(photo.image)"
                  alt=""
                  class="h-16 w-12 rounded object-cover"
                />
              </td>
              <td class="px-2 py-2 text-right">
                <button class="text-red-600 hover:underline" @click="removePhoto(photo.id)">
                  Удалить
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-else class="mt-4 text-sm text-[var(--color-text-muted)]">Пока пусто</p>
      </div>

      <div class="mt-6 rounded border border-gray-200 bg-white p-4">
        <h2 class="font-semibold">Мини-FAQ</h2>
        <form class="mt-3 grid gap-3 sm:grid-cols-2" @submit.prevent="addFaqItem">
          <input
            v-model="faqForm.question"
            type="text"
            placeholder="Вопрос"
            class="rounded border border-gray-300 px-3 py-2 sm:col-span-2"
          />
          <AdminMarkdownField
            v-model="faqForm.answer"
            placeholder="Ответ"
            :rows="2"
            class="sm:col-span-2"
          />
          <p v-if="faqError" class="text-sm text-red-600 sm:col-span-2">{{ faqError }}</p>
          <div>
            <button
              type="submit"
              :disabled="savingFaq || !faqForm.question || !faqForm.answer"
              class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
            >
              Добавить вопрос
            </button>
          </div>
        </form>

        <div v-if="faqItems.length" class="mt-4 flex flex-col gap-2">
          <div
            v-for="(item, index) in faqItems"
            :key="item.id"
            :data-row-index="index"
            class="flex items-start gap-3 border-b border-gray-100 py-2 last:border-0"
            :class="draggingFaqIndex === index ? 'opacity-50' : ''"
          >
            <AdminDragHandle @pointerdown="onFaqPointerDown(index, $event)" />
            <div class="min-w-0 flex-1">
              <p class="font-medium">{{ item.question }}</p>
              <p class="text-sm text-[var(--color-text-muted)]">{{ item.answer }}</p>
            </div>
            <button class="shrink-0 text-red-600 hover:underline" @click="removeFaqItem(item.id)">
              Удалить
            </button>
          </div>
        </div>
        <p v-else class="mt-4 text-sm text-[var(--color-text-muted)]">Пока пусто</p>
      </div>
    </template>
  </div>
</template>
