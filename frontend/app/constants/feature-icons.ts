import type { Component } from "vue";
import IconCalendarCheck from "~/components/ui/IconCalendarCheck.vue";
import IconChecklist from "~/components/ui/IconChecklist.vue";
import IconEightyTwenty from "~/components/ui/IconEightyTwenty.vue";
import IconFlexStart from "~/components/ui/IconFlexStart.vue";
import IconGift from "~/components/ui/IconGift.vue";
import IconLevels from "~/components/ui/IconLevels.vue";
import IconMapPin from "~/components/ui/IconMapPin.vue";
import IconMcLocation from "~/components/ui/IconMcLocation.vue";
import IconPeopleGroup from "~/components/ui/IconPeopleGroup.vue";
import IconPeopleTrio from "~/components/ui/IconPeopleTrio.vue";
import IconTulips from "~/components/ui/IconTulips.vue";
import IconTwoFormats from "~/components/ui/IconTwoFormats.vue";

/**
 * Site-wide "built-in" icon set — every place that accepts an icon value
 * (features, the gift-certificate marquee, ...) resolves a value that isn't
 * an uploaded icon ("icon:<id>", see useIconLibrary/AppIcon.vue) against
 * this map instead. Add a new built-in icon here, not per-page; a one-off
 * icon that doesn't need to live in code at all can just be uploaded
 * through the admin icon library instead.
 */
export const FEATURE_ICONS: Record<string, { component: Component; label: string }> = {
  "two-formats": { component: IconTwoFormats, label: "Два формата" },
  "people-group": { component: IconPeopleGroup, label: "Группа людей" },
  "map-pin": { component: IconMapPin, label: "Метка на карте" },
  "flex-start": { component: IconFlexStart, label: "Гибкий старт" },
  "eighty-twenty": { component: IconEightyTwenty, label: "80/20" },
  levels: { component: IconLevels, label: "Уровни" },
  "calendar-check": { component: IconCalendarCheck, label: "Календарь" },
  "people-trio": { component: IconPeopleTrio, label: "Трое людей" },
  "mc-location": { component: IconMcLocation, label: "Локация" },
  checklist: { component: IconChecklist, label: "Чек-лист" },
  gift: { component: IconGift, label: "Подарок" },
  tulips: { component: IconTulips, label: "Тюльпаны" },
};

export function featureIconComponent(key: string): Component | undefined {
  return FEATURE_ICONS[key]?.component;
}
