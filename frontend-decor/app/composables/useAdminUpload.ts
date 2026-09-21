/**
 * Uploads an image file to the backend (which stores it in Garage) and
 * returns the relative URL to save on the entity (e.g. `coverImage`,
 * `photo`, `gallery[]`). See AdminImageUpload.vue for the picker UI that
 * uses this.
 */

// Translates the backend's (English) error strings from upload_handler.go
// into the Russian text the admin UI shows. Matched by substring since the
// "file too large" message embeds a size that varies with maxUploadSize.
function describeUploadError(message: string | undefined): string {
  if (!message) return "Не удалось загрузить изображение";
  if (message.includes("unsupported image type")) {
    return "Неподдерживаемый формат файла. Разрешены: JPEG, PNG, WEBP, GIF";
  }
  if (message.includes("file too large")) {
    const sizeMatch = /max (\d+) MB/.exec(message);
    const size = sizeMatch ? sizeMatch[1] : "20";
    return `Файл слишком большой (максимум ${size} МБ)`;
  }
  return "Не удалось загрузить изображение";
}

export function useAdminUpload() {
  const api = useApiClient();
  const uploading = ref(false);
  const error = ref("");

  async function upload(file: File): Promise<string> {
    uploading.value = true;
    error.value = "";
    try {
      const body = new FormData();
      body.append("file", file);
      const res = await api<{ url: string }>("/api/v1/admin/uploads", {
        method: "POST",
        body,
      });
      return res.url;
    } catch (err) {
      error.value = describeUploadError(
        (err as { data?: { error?: string } } | undefined)?.data?.error,
      );
      throw new Error("upload failed");
    } finally {
      uploading.value = false;
    }
  }

  return { upload, uploading, error };
}
