export default defineNuxtRouteMiddleware(async () => {
  const { adminUser, authChecked, fetchMe } = useAdminAuth();

  if (!authChecked.value) {
    await fetchMe();
  }

  if (!adminUser.value) {
    return navigateTo("/admin/login");
  }
});
