const CYRILLIC_TO_LATIN: Record<string, string> = {
  а: "a",
  б: "b",
  в: "v",
  г: "g",
  д: "d",
  е: "e",
  ё: "e",
  ж: "zh",
  з: "z",
  и: "i",
  й: "y",
  к: "k",
  л: "l",
  м: "m",
  н: "n",
  о: "o",
  п: "p",
  р: "r",
  с: "s",
  т: "t",
  у: "u",
  ф: "f",
  х: "h",
  ц: "c",
  ч: "ch",
  ш: "sh",
  щ: "sch",
  ъ: "",
  ы: "y",
  ь: "",
  э: "e",
  ю: "yu",
  я: "ya",
};

/**
 * Transliterates Cyrillic to Latin, lowercases, and replaces any run of
 * non-alphanumeric characters with a single hyphen (leading/trailing
 * hyphens trimmed) — "Курс флористики" -> "kurs-floristiki".
 */
export function slugify(input: string): string {
  const transliterated = input
    .toLowerCase()
    .split("")
    .map((char) => CYRILLIC_TO_LATIN[char] ?? char)
    .join("");
  return transliterated.replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");
}
