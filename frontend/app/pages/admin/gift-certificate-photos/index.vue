<script setup lang="ts">
definePageMeta({ layout: "admin", middleware: "admin-auth" });

await navigateTo("/admin/gallery-photos", { redirectCode: 301 });
</script>
