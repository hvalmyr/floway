<script setup lang="ts">
/**
 * Clickable status dot for list-view "visible/hidden" columns — replaces
 * plain "Показан"/"Скрыт" text. Purely presentational: the parent owns the
 * actual toggle (API call), this just renders the dot and emits a click.
 *
 * @example
 * <AdminVisibilityDot :visible="course.visible" @click="onToggleVisible(course)" />
 */
defineProps<{ visible: boolean; disabled?: boolean }>();
defineEmits<{ click: [] }>();
</script>

<template>
  <button
    type="button"
    class="size-5 rounded-full transition-colors disabled:cursor-not-allowed disabled:opacity-50"
    :class="visible ? 'bg-green-500 hover:bg-green-600' : 'bg-gray-300 hover:bg-gray-400'"
    :title="visible ? 'Показан — нажмите, чтобы скрыть' : 'Скрыт — нажмите, чтобы показать'"
    :disabled="disabled"
    @click="$emit('click')"
  />
</template>
