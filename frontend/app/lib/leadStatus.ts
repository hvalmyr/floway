import type { LeadStatus } from "~/types/api";

/** Ordered list of every lead status — the single source of truth so the
 * status filter pills, bulk-change dropdown, and per-row dropdown never
 * drift out of sync with each other or with the backend's LeadStatus enum
 * (backend/internal/model/model.go). */
export const LEAD_STATUSES: LeadStatus[] = [
  "new",
  "in_progress",
  "waiting_client",
  "booked",
  "postponed",
  "closed_won",
  "closed_lost",
];

export const LEAD_STATUS_LABELS: Record<LeadStatus, string> = {
  new: "Новая",
  in_progress: "В работе",
  waiting_client: "Ждём ответа клиента",
  booked: "Записан/оплачено",
  postponed: "Отложено",
  closed_won: "Закрыта: успешно",
  closed_lost: "Закрыта: отказ",
};
