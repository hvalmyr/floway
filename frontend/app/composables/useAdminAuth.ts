interface AdminUser {
  id: number;
  login: string;
}

const adminUser = () => useState<AdminUser | null>("admin-user", () => null);
const authChecked = () => useState<boolean>("admin-auth-checked", () => false);

export function useAdminAuth() {
  const api = useApi();
  const user = adminUser();
  const checked = authChecked();

  async function fetchMe() {
    try {
      user.value = await api<AdminUser>("/api/v1/admin/me");
    } catch {
      user.value = null;
    } finally {
      checked.value = true;
    }
  }

  async function login(loginValue: string, password: string) {
    await api("/api/v1/admin/login", {
      method: "POST",
      body: { login: loginValue, password },
    });
    await fetchMe();
  }

  async function logout() {
    try {
      await api("/api/v1/admin/logout", { method: "POST" });
    } finally {
      user.value = null;
      checked.value = true;
    }
  }

  return { adminUser: user, authChecked: checked, fetchMe, login, logout };
}
