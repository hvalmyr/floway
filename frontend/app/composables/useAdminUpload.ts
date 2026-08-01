/**
 * Uploads an image file to the backend (which stores it in Garage) and
 * returns the relative URL to save on the entity (e.g. `coverImage`,
 * `photo`, `gallery[]`). See AdminImageUpload.vue for the picker UI that
 * uses this.
 */
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
    } catch {
      error.value = "Не удалось загрузить изображение";
      throw new Error("upload failed");
    } finally {
      uploading.value = false;
    }
  }

  return { upload, uploading, error };
}
